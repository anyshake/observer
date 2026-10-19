package miniseed

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/fileutil"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/bclswl0827/mseedio"
	"github.com/samber/lo"
)

func (s *MiniSeedServiceImpl) getAppendInterval(gnssAvailable bool) int {
	// Set write interval to 1 if GNSS time is not available
	// This is a simple solution to sample rate and timestamp jittering
	// However, it will increase the disk I/O and file size
	if !gnssAvailable {
		return 1
	}
	return MINISEED_APPEND_INTERVAL
}

func (s *MiniSeedServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("MiniSEED service is already running")
	}

	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	logger.GetLogger(ID).Infof("generated MiniSEED files will be saved to %s", s.filePath)
	s.dataSequence.filePath = fmt.Sprintf("%s/.miniseed_sequence", s.filePath)
	s.dataSequence.sequenceData = make(map[string]uint32)
	s.cleanupCountDown = MINISEED_CLEANUP_INTERVAL

	hardwareConfig := s.hardwareDev.GetConfig()
	s.appendCountDown = s.getAppendInterval(hardwareConfig.GetGnssAvailability())
	s.recordBuffer = make([][]buffer, s.appendCountDown)

	runCtx, runCancel := s.ctx, s.cancelFn
	err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
		defer service.RecoverCallback(ID, runCancel)
		s.mu.Lock()
		defer s.mu.Unlock()

		if runCtx.Err() != nil {
			return
		}
		currentInterval := s.getAppendInterval(event.GNSSAvailable)
		if len(s.recordBuffer) != currentInterval {
			s.recordBuffer = make([][]buffer, currentInterval)
			s.appendCountDown = currentInterval
		}

		bufferIndex := len(s.recordBuffer) - s.appendCountDown
		if bufferIndex >= 0 && bufferIndex < len(s.recordBuffer) {
			s.recordBuffer[bufferIndex] = make([]buffer, len(event.ChannelData))
		}
		lo.ForEach(event.ChannelData, func(ch explorer.ChannelData, idx int) {
			s.recordBuffer[bufferIndex][idx] = buffer{
				Timestamp:   event.Timestamp.UnixMilli(),
				SampleRate:  event.SampleRate,
				ChannelData: ch,
			}
		})

		s.appendCountDown--
		if s.lifeCycle > 0 {
			s.cleanupCountDown--
		}

		if s.appendCountDown == 0 {
			s.appendCountDown = currentInterval
			channels, err := s.saveMiniSeedRecords(event.GNSSAvailable)
			if err != nil {
				logger.GetLogger(ID).Errorf("failed to append records to MiniSEED file: %v", err)
				return
			}
			logger.GetLogger(ID).Infof("%d records in %d channels have been added to MiniSEED file", len(s.recordBuffer), channels)
		}
		if s.cleanupCountDown == 0 {
			s.cleanupCountDown = MINISEED_CLEANUP_INTERVAL
			endTime := event.Timestamp.Add(time.Duration(-s.lifeCycle) * time.Hour * 24)
			if err := s.cleanupMiniSeedRecords(runCtx, endTime); err != nil {
				if !errors.Is(err, context.Canceled) {
					logger.GetLogger(ID).Errorf("failed to purge expired MiniSEED files: %v", err)
				}
				return
			}
			logger.GetLogger(ID).Infoln("expired MiniSEED files have been purged from storage")
		}
	})
	if err != nil {
		runCancel()
		return fmt.Errorf("failed to subscribe to hardware message bus: %w", err)
	}

	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
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

type sequence struct {
	filePath     string
	sequenceData map[string]uint32
}

func (s *sequence) Save(data map[string]uint32) error {
	s.sequenceData = data

	if s.filePath == "" {
		return errors.New("file path is empty")
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func (s *sequence) Read() (map[string]uint32, error) {
	if s.filePath == "" {
		return nil, errors.New("file path is empty")
	}

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return s.sequenceData, nil
		}
		return nil, err
	}
	defer file.Close()

	data := make(map[string]uint32)
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return s.sequenceData, nil
	}

	return data, nil

}

func (m *MiniSeedServiceImpl) getMiniSeedFileName(tm time.Time, channelCode string) string {
	return fmt.Sprintf("%s.%s.%s.%s.D.%s.%s.mseed",
		m.networkCode, m.stationCode, m.locationCode, channelCode,
		tm.Format("2006"),
		tm.Format("002"),
	)
}

func (s *MiniSeedServiceImpl) saveMiniSeedRecords(gnssAvailable bool) (int, error) {
	if len(s.recordBuffer) == 0 || len(s.recordBuffer[0]) == 0 {
		return 0, errors.New("no data to save")
	}

	dataSeqMap, err := s.dataSequence.Read()
	if err != nil {
		logger.GetLogger(ID).Warnf("failed to read data sequence, starting from 0: %v", err)
	}

	allowedJitterMs := lo.Ternary[float64](gnssAvailable, explorer.ALLOWED_JITTER_MS_GNSS, explorer.ALLOWED_JITTER_MS_NTP)

	startTimestamp := s.recordBuffer[0][0].Timestamp
	startSampleRate := s.recordBuffer[0][0].SampleRate
	startChannels := len(s.recordBuffer[0])
	for i := 1; i < len(s.recordBuffer); i++ {
		expectedTimestamp := startTimestamp + int64(i)*1000
		if math.Abs(float64(s.recordBuffer[i][0].Timestamp-expectedTimestamp)) >= allowedJitterMs {
			return 0, fmt.Errorf(
				"timestamp out of jitter range %d ms, expected %d, got %d",
				int(allowedJitterMs),
				expectedTimestamp,
				s.recordBuffer[i][0].Timestamp,
			)
		}

		if s.recordBuffer[i][0].SampleRate != startSampleRate {
			return 0, fmt.Errorf(
				"sample rate mismatch, expected %d, got %d",
				startSampleRate,
				s.recordBuffer[i][0].SampleRate,
			)
		}

		if len(s.recordBuffer[i]) != startChannels {
			return 0, fmt.Errorf(
				"number of channels mismatch, expected %d, got %d",
				startChannels,
				len(s.recordBuffer[i]),
			)
		}
	}

	for channelIdx := 0; channelIdx < startChannels; channelIdx++ {
		var channelBuffer []int32
		for recordIdx := 0; recordIdx < len(s.recordBuffer); recordIdx++ {
			channelBuffer = append(channelBuffer, s.recordBuffer[recordIdx][channelIdx].Data...)
		}

		encoding := mseedio.INT32
		if s.useCompress {
			switch s.compressScheme {
			case COMPRESS_SCHEME_STEIM_2:
				encoding = mseedio.STEIM2
			default:
				encoding = mseedio.STEIM1
			}
		}
		var miniseed mseedio.MiniSeedData
		if err := miniseed.Init(encoding, mseedio.MSBFIRST); err != nil {
			return 0, err
		}

		channelCode := s.recordBuffer[0][channelIdx].ChannelCode
		seq, ok := dataSeqMap[channelCode]
		if !ok {
			seq = 0
		}

		dateDir := filepath.Join(s.filePath, time.UnixMilli(startTimestamp).UTC().Format("2006-01-02"))
		filePath := filepath.Join(dateDir, s.getMiniSeedFileName(time.UnixMilli(startTimestamp).UTC(), channelCode))

		if err := os.MkdirAll(dateDir, 0755); err != nil {
			return 0, fmt.Errorf("failed to create directory %s: %w", dateDir, err)
		}

		err = miniseed.Append(channelBuffer, &mseedio.AppendOptions{
			SequenceNumber: fmt.Sprintf("%06d", seq),
			SampleRate:     float64(startSampleRate),
			StartTime:      time.UnixMilli(startTimestamp).UTC(),
			ChannelCode:    channelCode,
			StationCode:    s.stationCode,
			NetworkCode:    s.networkCode,
			LocationCode:   s.locationCode,
		})
		if err != nil {
			return 0, err
		}

		dataBytes, err := miniseed.Encode(mseedio.APPEND, mseedio.MSBFIRST)
		if err != nil {
			return 0, err
		}

		if err = miniseed.Write(filePath, mseedio.APPEND, dataBytes); err != nil {
			return 0, fmt.Errorf("failed to write MiniSEED file: %w", err)
		}

		dataSeqMap[channelCode]++
		dataSeqMap[channelCode] %= 999999 + 1
	}

	if err = s.dataSequence.Save(dataSeqMap); err != nil {
		return 0, fmt.Errorf("failed to save data sequence: %w", err)
	}

	return startChannels, nil
}

func (s *MiniSeedServiceImpl) cleanupMiniSeedRecords(ctx context.Context, until time.Time) error {
	return s.actionHandler.CleanupExclusive(ctx, func() error {
		return s.removeMiniSeedRecordsByDateRange(ctx, time.Unix(0, 0), until)
	})
}

func (s *MiniSeedServiceImpl) PurgeMiniSeedFiles(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.actionHandler.CleanupExclusive(ctx, func() error {
		if err := fileutil.RemoveAll(ctx, s.filePath); err != nil {
			return fmt.Errorf("failed to remove MiniSEED directory %s: %w", s.filePath, err)
		}
		return nil
	})
}

func (s *MiniSeedServiceImpl) PurgeMiniSeedFilesByDate(ctx context.Context, startDate, endDate time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	startTime, endTime, err := normalizeDateRange(startDate, endDate)
	if err != nil {
		return err
	}

	return s.actionHandler.CleanupExclusive(ctx, func() error {
		return s.removeMiniSeedRecordsByDateRange(ctx, startTime, endTime)
	})
}

func (s *MiniSeedServiceImpl) removeMiniSeedRecordsByDateRange(ctx context.Context, startDate, endDate time.Time) error {
	startTime, endTime, err := normalizeDateRange(startDate, endDate)
	if err != nil {
		return err
	}

	for currentDate := startTime; !currentDate.After(endTime); currentDate = currentDate.AddDate(0, 0, 1) {
		if err := ctx.Err(); err != nil {
			return err
		}
		dirPath := filepath.Join(s.filePath, currentDate.UTC().Format("2006-01-02"))
		if err := fileutil.RemoveAll(ctx, dirPath); err != nil {
			return fmt.Errorf("failed to remove MiniSEED directory %s: %w", dirPath, err)
		}
	}

	return nil
}

func normalizeDateRange(startDate, endDate time.Time) (time.Time, time.Time, error) {
	startTime := time.Date(startDate.UTC().Year(), startDate.UTC().Month(), startDate.UTC().Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(endDate.UTC().Year(), endDate.UTC().Month(), endDate.UTC().Day(), 0, 0, 0, 0, time.UTC)
	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, errors.New("start date is after end date")
	}

	return startTime, endTime, nil
}
