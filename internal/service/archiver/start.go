package archiver

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/logger"
)

func (s *ArchiverServiceImpl) startExpiredSeisRecordsPurge(ctx context.Context, endTime time.Time) {
	if ctx.Err() != nil {
		return
	}
	if !s.cleanupRunning.CompareAndSwap(false, true) {
		logger.GetLogger(ID).Infoln("skipped expired seismic waveform records purge because previous archiver cleanup task is still running")
		return
	}

	s.wg.Add(1)
	go s.purgeExpiredSeisRecords(ctx, endTime)
}

func (s *ArchiverServiceImpl) purgeExpiredSeisRecords(ctx context.Context, endTime time.Time) {
	defer s.wg.Done()
	defer s.cleanupRunning.Store(false)
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger(ID).Errorf("expired seismic waveform records purge unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()

	if err := s.actionHandler.SeisRecordsPurge(ctx, time.Unix(0, 0), endTime); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, action.ErrCleanupRunning) {
			logger.GetLogger(ID).Infoln("skipped expired seismic waveform records purge because another cleanup task is running")
			return
		}

		logger.GetLogger(ID).Errorf("failed to purge expired seismic waveform records: %v", err)
		return
	}
	logger.GetLogger(ID).Infoln("expired seismic waveform records have been purged from database")
}

func (s *ArchiverServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("archiver service is already running")
	}

	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	rotation := s.rotation
	queue := newRecordQueue()
	writer := &recordWriter{queue: queue, insert: s.actionHandler.SeisRecordsCreate}

	runCtx, runCancel := s.ctx, s.cancelFn
	// The callback runs on the hardware message bus. It must never wait for the
	// database: it only encodes the event and hands it to the writer goroutine.
	err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
		defer service.RecoverCallback(ID, runCancel)
		if runCtx.Err() != nil {
			return
		}
		record := model.SeisRecord{}
		if err := record.Encode(event.Timestamp, event.SampleRate, event.ChannelData); err != nil {
			logger.GetLogger(ID).Errorf("failed to encode seismic waveform record: %v", err)
			return
		}
		queue.push(record)
		if rotation > 0 && queue.tickCleanup() {
			endTime := event.Timestamp.Add(time.Duration(-rotation) * time.Hour * 24)
			s.startExpiredSeisRecordsPurge(runCtx, endTime)
		}
	})
	if err != nil {
		runCancel()
		return fmt.Errorf("failed to subscribe to hardware message bus: %w", err)
	}

	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(2)
	go func() {
		defer s.wg.Done()
		defer service.RecoverWorker(ID, runCancel)
		writer.run(runCtx)
	}()
	go func() {
		defer s.wg.Done()
		defer func() {
			runCancel()
			_ = s.hardwareDev.Unsubscribe(ID)
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()
		<-runCtx.Done()
	}()

	return nil
}
