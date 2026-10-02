package archiver

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
)

func (s *ArchiverServiceImpl) handleInterrupt() {
	s.wg.Done()
}

func (s *ArchiverServiceImpl) startExpiredSeisRecordsPurge(endTime time.Time) {
	if !s.cleanupRunning.CompareAndSwap(false, true) {
		logger.GetLogger(ID).Infoln("skipped expired seismic waveform records purge because previous archiver cleanup task is still running")
		return
	}

	go s.purgeExpiredSeisRecords(endTime)
}

func (s *ArchiverServiceImpl) purgeExpiredSeisRecords(endTime time.Time) {
	defer s.cleanupRunning.Store(false)
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger(ID).Errorf("expired seismic waveform records purge unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()

	if err := s.actionHandler.SeisRecordsPurge(time.Unix(0, 0), endTime); err != nil {
		if errors.Is(err, action.ErrCleanupRunning) {
			logger.GetLogger(ID).Infoln("skipped expired seismic waveform records purge because another cleanup task is running")
			return
		}

		logger.GetLogger(ID).Errorf("failed to purge expired seismic waveform records: %v", err)
		return
	}
	logger.GetLogger(ID).Infoln("expired seismic waveform records have been purged from database")
}

func (s *ArchiverServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	s.cleanupCountDown = RECORDS_CLEANUP_INTERVAL
	s.insertCountDown = RECORDS_INSERT_INTERVAL
	s.recordBuffer = make([]model.SeisRecord, RECORDS_INSERT_INTERVAL)

	go func() {
		s.status.SetStartedAt(s.timeSource.Now())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt()
				_ = s.Stop()
			}
		}()

		err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
			s.mu.Lock()
			defer s.mu.Unlock()

			record := model.SeisRecord{}
			if err := record.Encode(event.Timestamp, event.SampleRate, event.ChannelData); err != nil {
				logger.GetLogger(ID).Errorf("failed to encode seismic waveform record: %v", err)
				return
			}
			s.recordBuffer[len(s.recordBuffer)-s.insertCountDown] = record

			s.insertCountDown--
			if s.rotation > 0 {
				s.cleanupCountDown--
			}

			if s.insertCountDown == 0 {
				s.insertCountDown = RECORDS_INSERT_INTERVAL
				if err := s.actionHandler.SeisRecordsCreate(s.recordBuffer...); err != nil {
					logger.GetLogger(ID).Errorf("failed to create seismic waveform records: %v", err)
					return
				}
				logger.GetLogger(ID).Infof("%d seismic waveform records have been inserted to database", len(s.recordBuffer))
			}
			if s.cleanupCountDown == 0 {
				s.cleanupCountDown = RECORDS_CLEANUP_INTERVAL
				endTime := event.Timestamp.Add(time.Duration(-s.rotation) * time.Hour * 24)
				s.startExpiredSeisRecordsPurge(endTime)
			}
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		<-s.ctx.Done()
		s.handleInterrupt()
	}()

	s.wg.Add(1)
	return nil
}
