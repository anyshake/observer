package archiver

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
)

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
		err := s.hardwareDev.Subscribe(ID, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			s.mu.Lock()
			defer s.mu.Unlock()

			record := model.SeisRecord{}
			if err := record.Encode(t, di.GetSampleRate(), cd); err != nil {
				logger.GetLogger(ID).Errorf("failed to encode siesmic waveform record: %v", err)
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
					logger.GetLogger(ID).Errorf("failed to create siesmic waveform records: %v", err)
					return
				}
				logger.GetLogger(ID).Infof("%d siesmic waveform records have been inserted to database", len(s.recordBuffer))
			}
			if s.cleanupCountDown == 0 {
				s.cleanupCountDown = RECORDS_CLEANUP_INTERVAL
				endTime := t.Add(time.Duration(-s.rotation) * time.Hour * 24)
				if err := s.actionHandler.SeisRecordsPurge(time.Unix(0, 0), endTime); err != nil {
					logger.GetLogger(ID).Errorf("failed to purge expired siesmic waveform records: %v", err)
					return
				}
				logger.GetLogger(ID).Infoln("expired siesmic waveform records have been purged from database")
			}
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)

		<-s.ctx.Done()
		s.wg.Done()
	}()

	s.wg.Add(1)
	return nil
}
