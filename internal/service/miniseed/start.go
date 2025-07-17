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
	"github.com/anyshake/observer/pkg/logger"
	"github.com/bclswl0827/mseedio"
	"github.com/samber/lo"
)

func (s *MiniSeedServiceImpl) handleInterrupt() {
	s.wg.Done()
}

func (s *MiniSeedServiceImpl) getAppendInterval(cfg *explorer.DeviceConfig) int {
	if cfg == nil {
		cfgVal := s.hardwareDev.GetConfig()
		cfg = &cfgVal
	}
	// Set write interval to 1 if GNSS time is not available
	// This is a simple solution to sample rate and timestamp jittering
	// However, it will increase the disk I/O and file size
	if !cfg.GetGnssAvailability() {
		return 1
	}
	return MINISEED_APPEND_INTERVAL
}

func (s *MiniSeedServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	logger.GetLogger(ID).Infof("generated MiniSEED files will be saved to %s", s.filePath)
	s.dataSequence.filePath = fmt.Sprintf("%s/.miniseed_sequence", s.filePath)
	s.dataSequence.sequenceData = make(map[string]uint32)
	s.cleanupCountDown = MINISEED_CLEANUP_INTERVAL

	s.appendCountDown = s.getAppendInterval(nil)
	s.recordBuffer = make([][]buffer, s.appendCountDown)

	go func() {
		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt()
				_ = s.Stop()
			}
		}()

		err := s.hardwareDev.Subscribe(ID, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			s.mu.Lock()
			defer s.mu.Unlock()

			currentInterval := s.getAppendInterval(di)
			if len(s.recordBuffer) != currentInterval {
				s.recordBuffer = make([][]buffer, s.appendCountDown)
				s.appendCountDown = currentInterval
			}

			bufferIndex := len(s.recordBuffer) - s.appendCountDown
			if bufferIndex >= 0 && bufferIndex < len(s.recordBuffer) {
				s.recordBuffer[bufferIndex] = make([]buffer, len(cd))
			}
			lo.ForEach(cd, func(ch explorer.ChannelData, idx int) {
				s.recordBuffer[bufferIndex][idx] = buffer{
					Timestamp:   t.UnixMilli(),
					SampleRate:  di.GetSampleRate(),
					ChannelData: ch,
				}
			})

			s.appendCountDown--
			if s.lifeCycle > 0 {
				s.cleanupCountDown--
			}

			if s.appendCountDown == 0 {
				s.appendCountDown = currentInterval
				channels, err := s.saveMiniSeedRecords(di)
				if err != nil {
					logger.GetLogger(ID).Errorf("failed to append records to MiniSEED file: %v", err)
					return
				}
				logger.GetLogger(ID).Infof("%d records in %d channels have been added to MiniSEED file", len(s.recordBuffer), channels)
			}
			if s.cleanupCountDown == 0 {
				s.cleanupCountDown = MINISEED_CLEANUP_INTERVAL
				endTime := t.Add(time.Duration(-s.lifeCycle) * time.Hour * 24)
				if err := s.cleanupMiniSeedRecords(endTime); err != nil {
					logger.GetLogger(ID).Errorf("failed to purge expired MiniSEED files: %v", err)
					return
				}
				logger.GetLogger(ID).Infoln("expired MiniSEED files have been purged from storage")
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

func (s *MiniSeedServiceImpl) saveMiniSeedRecords(cfg *explorer.DeviceConfig) (int, error) {
	if len(s.recordBuffer) == 0 || len(s.recordBuffer[0]) == 0 {
		return 0, errors.New("no data to save")
	}

	dataSeqMap, err := s.dataSequence.Read()
	if err != nil {
		logger.GetLogger(ID).Warnf("failed to read data sequence, starting from 0: %v", err)
	}

	if cfg == nil {
		cfgVal := s.hardwareDev.GetConfig()
		cfg = &cfgVal
	}
	allowedJitterMs := lo.Ternary[float64](cfg.GetGnssAvailability(), explorer.ALLOWED_JITTER_MS_GNSS, explorer.ALLOWED_JITTER_MS_NTP)

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

		encodeing := mseedio.INT32
		if s.useCompress {
			encodeing = mseedio.STEIM2
		}
		var miniseed mseedio.MiniSeedData
		if err := miniseed.Init(encodeing, mseedio.MSBFIRST); err != nil {
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

func (s *MiniSeedServiceImpl) cleanupMiniSeedRecords(until time.Time) error {
	err := filepath.Walk(s.filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dirTime := info.ModTime()
			if dirTime.Before(until) {
				err := os.RemoveAll(path)
				if err != nil {
					return fmt.Errorf("failed to remove directory: %w", err)
				}
			}
		}

		return nil
	})

	return err
}
