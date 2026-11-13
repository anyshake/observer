package archiver

import (
	"context"
	"errors"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
)

type fakeInserter struct {
	mu      sync.Mutex
	block   chan struct{}
	fail    int
	failErr error
	calls   int
	written []int64
}

func (f *fakeInserter) insert(ctx context.Context, records ...model.SeisRecord) error {
	f.mu.Lock()
	f.calls++
	failing := f.fail > 0
	if failing {
		f.fail--
	}
	block := f.block
	f.mu.Unlock()
	if block != nil {
		select {
		case <-block:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if failing {
		if f.failErr != nil {
			return f.failErr
		}
		return errors.New("database is locked")
	}
	f.mu.Lock()
	for _, record := range records {
		f.written = append(f.written, record.RecordTime)
	}
	f.mu.Unlock()
	return nil
}

func (f *fakeInserter) snapshot() []int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]int64(nil), f.written...)
}

func record(at int64) model.SeisRecord { return model.SeisRecord{RecordTime: at} }

func TestQueuePushNeverBlocksOnStalledDatabase(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inserter := &fakeInserter{block: make(chan struct{})}
		queue := newRecordQueue()
		writer := &recordWriter{queue: queue, insert: inserter.insert}
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { writer.run(ctx); close(done) }()

		// One hour of records while every insert hangs: push must return
		// immediately and the backlog must stay bounded.
		for i := range 3600 {
			queue.push(record(int64(i)))
			time.Sleep(time.Second)
		}
		if size := queue.size(); size > RECORDS_PENDING_LIMIT {
			t.Fatalf("backlog = %d, want <= %d", size, RECORDS_PENDING_LIMIT)
		}

		close(inserter.block)
		time.Sleep(RECORDS_FLUSH_INTERVAL + time.Second)
		written := inserter.snapshot()
		if len(written) == 0 || written[len(written)-1] != 3599 {
			t.Fatalf("latest record was not written after recovery: %v", len(written))
		}
		cancel()
		<-done
	})
}

func TestWriterRetriesFailedBatchInOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inserter := &fakeInserter{fail: 2}
		queue := newRecordQueue()
		writer := &recordWriter{queue: queue, insert: inserter.insert}
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { writer.run(ctx); close(done) }()

		for i := range RECORDS_INSERT_INTERVAL * 3 {
			queue.push(record(int64(i)))
			time.Sleep(time.Second)
		}
		time.Sleep(time.Minute)
		written := inserter.snapshot()
		if len(written) != RECORDS_INSERT_INTERVAL*3 {
			t.Fatalf("written = %d records, want %d", len(written), RECORDS_INSERT_INTERVAL*3)
		}
		for i, at := range written {
			if at != int64(i) {
				t.Fatalf("record %d has time %d; order was not preserved", i, at)
			}
		}
		cancel()
		<-done
	})
}

func TestWriterDiscardsPermanentlyFailingBatch(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inserter := &fakeInserter{fail: RECORDS_INSERT_MAX_ATTEMPTS, failErr: errors.New("NOT NULL constraint failed")}
		queue := newRecordQueue()
		writer := &recordWriter{queue: queue, insert: inserter.insert}
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { writer.run(ctx); close(done) }()

		for i := range RECORDS_INSERT_INTERVAL {
			queue.push(record(int64(i)))
		}
		time.Sleep(2 * time.Minute)
		queue.push(record(100))
		time.Sleep(RECORDS_FLUSH_INTERVAL + time.Second)
		if written := inserter.snapshot(); len(written) != 1 || written[0] != 100 {
			t.Fatalf("written = %v, want only the record after the discarded batch", written)
		}
		cancel()
		<-done
	})
}

func TestWriterFlushesPendingRecordsOnStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inserter := &fakeInserter{}
		queue := newRecordQueue()
		writer := &recordWriter{queue: queue, insert: inserter.insert}
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { writer.run(ctx); close(done) }()

		queue.push(record(1))
		queue.push(record(2))
		cancel()
		<-done
		if written := inserter.snapshot(); len(written) != 2 {
			t.Fatalf("written = %v, want pending records flushed on stop", written)
		}
	})
}

func TestWriterStopIsBoundedWhenDatabaseHangs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inserter := &fakeInserter{block: make(chan struct{})}
		queue := newRecordQueue()
		writer := &recordWriter{queue: queue, insert: inserter.insert}
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { writer.run(ctx); close(done) }()

		for i := range RECORDS_INSERT_INTERVAL {
			queue.push(record(int64(i)))
		}
		time.Sleep(time.Second)
		started := time.Now()
		cancel()
		<-done
		if elapsed := time.Since(started); elapsed > RECORDS_FINAL_FLUSH_TIMEOUT {
			t.Fatalf("stop took %s with a hung database", elapsed)
		}
	})
}

func TestWriterKeepsRetryingTransientFailures(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inserter := &fakeInserter{fail: RECORDS_INSERT_MAX_ATTEMPTS * 3, failErr: context.DeadlineExceeded}
		queue := newRecordQueue()
		writer := &recordWriter{queue: queue, insert: inserter.insert}
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { writer.run(ctx); close(done) }()

		for i := range RECORDS_INSERT_INTERVAL {
			queue.push(record(int64(i)))
		}
		time.Sleep(5 * time.Minute)
		if written := inserter.snapshot(); len(written) != RECORDS_INSERT_INTERVAL {
			t.Fatalf("written = %v, want the batch kept through transient failures", written)
		}
		cancel()
		<-done
	})
}
