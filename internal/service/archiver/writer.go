package archiver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/pkg/logger"
)

// recordQueue decouples the hardware message bus from the database. The bus
// callback only appends encoded records here and never waits for SQLite, so a
// slow insert, a long purge or a long-running query can no longer back-pressure
// the hardware stream (which used to surface as "message delivery timed out").
type recordQueue struct {
	mu      sync.Mutex
	pending []model.SeisRecord
	dropped int
	notify  chan struct{}

	cleanupCountDown int
}

func newRecordQueue() *recordQueue {
	return &recordQueue{
		pending:          make([]model.SeisRecord, 0, RECORDS_INSERT_INTERVAL),
		notify:           make(chan struct{}, 1),
		cleanupCountDown: RECORDS_CLEANUP_INTERVAL,
	}
}

// push appends a record without blocking. When the database has been stalled
// long enough to fill the queue, the oldest records are discarded.
func (q *recordQueue) push(record model.SeisRecord) {
	q.mu.Lock()
	q.pending = append(q.pending, record)
	q.trimLocked()
	ready := len(q.pending) >= RECORDS_INSERT_INTERVAL
	q.mu.Unlock()
	if ready {
		select {
		case q.notify <- struct{}{}:
		default:
		}
	}
}

// tickCleanup counts one archived second and reports whether a purge is due.
func (q *recordQueue) tickCleanup() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.cleanupCountDown--
	if q.cleanupCountDown > 0 {
		return false
	}
	q.cleanupCountDown = RECORDS_CLEANUP_INTERVAL
	return true
}

func (q *recordQueue) take(limit int) []model.SeisRecord {
	q.mu.Lock()
	defer q.mu.Unlock()
	n := min(limit, len(q.pending))
	if n == 0 {
		return nil
	}
	batch := slices.Clone(q.pending[:n])
	q.pending = append(q.pending[:0], q.pending[n:]...)
	return batch
}

// requeue puts a failed batch back in front of newer records.
func (q *recordQueue) requeue(batch []model.SeisRecord) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.pending = append(slices.Clone(batch), q.pending...)
	q.trimLocked()
}

func (q *recordQueue) trimLocked() {
	if excess := len(q.pending) - RECORDS_PENDING_LIMIT; excess > 0 {
		q.pending = append(q.pending[:0], q.pending[excess:]...)
		q.dropped += excess
	}
}

func (q *recordQueue) takeDropped() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	dropped := q.dropped
	q.dropped = 0
	return dropped
}

func (q *recordQueue) size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.pending)
}

// recordWriter is the only goroutine that inserts waveform records. Every
// insert runs with its own timeout, so a stuck database connection delays the
// archive but never the hardware stream.
type recordWriter struct {
	queue      *recordQueue
	insert     func(context.Context, ...model.SeisRecord) error
	dbStats    func() (sql.DBStats, bool)
	failures   int
	retryAfter time.Time
}

// describeFailure records enough context to tell a connection that never became
// available from a statement that ran too long.
func (w *recordWriter) describeFailure(batch []model.SeisRecord, elapsed time.Duration) string {
	description := fmt.Sprintf("batch=%d records (%d..%d), elapsed=%s, backlog=%d",
		len(batch), batch[0].RecordTime, batch[len(batch)-1].RecordTime,
		elapsed.Round(time.Millisecond), w.queue.size())
	if w.dbStats == nil {
		return description
	}
	if stats, ok := w.dbStats(); ok {
		description += fmt.Sprintf(", pool in_use=%d idle=%d max_open=%d wait_count=%d wait_duration=%s",
			stats.InUse, stats.Idle, stats.MaxOpenConnections, stats.WaitCount, stats.WaitDuration.Round(time.Millisecond))
	}
	return description
}

func (w *recordWriter) run(ctx context.Context) {
	ticker := time.NewTicker(RECORDS_FLUSH_INTERVAL)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			w.finalFlush()
			return
		case <-w.queue.notify:
			if time.Now().Before(w.retryAfter) {
				continue
			}
		case <-ticker.C:
		}
		w.flush(ctx)
	}
}

func (w *recordWriter) flush(ctx context.Context) {
	if dropped := w.queue.takeDropped(); dropped > 0 {
		logger.GetLogger(ID).Errorf("discarded %d seismic waveform records because the database did not keep up", dropped)
	}
	for ctx.Err() == nil {
		batch := w.queue.take(RECORDS_INSERT_BATCH_LIMIT)
		if len(batch) == 0 {
			return
		}
		insertCtx, cancel := context.WithTimeout(ctx, RECORDS_INSERT_TIMEOUT)
		started := time.Now()
		err := w.insert(insertCtx, batch...)
		elapsed := time.Since(started)
		cancel()
		if err == nil {
			w.failures = 0
			w.retryAfter = time.Time{}
			logger.GetLogger(ID).Infof("%d seismic waveform records have been inserted to database", len(batch))
			continue
		}
		if ctx.Err() != nil {
			// Stopping: keep the batch for the final flush.
			w.queue.requeue(batch)
			return
		}
		// A busy or slow database is temporary: keep the batch (bounded by the
		// queue limit) instead of turning a stall into a permanent gap. Only
		// errors that retries cannot fix count toward discarding the batch.
		if !isTransientInsertError(err) {
			w.failures++
		}
		if w.failures >= RECORDS_INSERT_MAX_ATTEMPTS {
			logger.GetLogger(ID).Errorf("discarded %d seismic waveform records after %d failed insert attempts: %v [%s]", len(batch), w.failures, err, w.describeFailure(batch, elapsed))
			w.failures = 0
			continue
		}
		w.queue.requeue(batch)
		w.retryAfter = time.Now().Add(RECORDS_INSERT_RETRY_DELAY)
		logger.GetLogger(ID).Warnf("failed to create seismic waveform records, will retry: %v [%s]", err, w.describeFailure(batch, elapsed))
		return
	}
}

// finalFlush writes what is left after the run was canceled, bounded by its
// own deadline so that stopping the service cannot hang on the database.
func (w *recordWriter) finalFlush() {
	ctx, cancel := context.WithTimeout(context.Background(), RECORDS_FINAL_FLUSH_TIMEOUT)
	defer cancel()
	for ctx.Err() == nil {
		batch := w.queue.take(RECORDS_INSERT_BATCH_LIMIT)
		if len(batch) == 0 {
			return
		}
		if err := w.insert(ctx, batch...); err != nil {
			logger.GetLogger(ID).Errorf("failed to write %d pending seismic waveform records while stopping: %v", len(batch)+w.queue.size(), err)
			return
		}
	}
}

func isTransientInsertError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	message := strings.ToLower(err.Error())
	for _, marker := range []string{"database is locked", "sqlite_busy", "sqlite_interrupt", "interrupted", "deadline exceeded"} {
		if strings.Contains(message, marker) {
			return true
		}
	}
	return false
}
