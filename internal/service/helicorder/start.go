package helicorder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/bclswl0827/heligo"
)

type queryCacheData struct {
	sampleRate  int
	timestamp   int64
	channelData []explorer.ChannelData
}

type provider struct {
	actionHandler *action.Handler
	queryCache    cache.KvCache[[]queryCacheData]

	stationCode  string
	networkCode  string
	locationCode string

	channelCode      string
	channelCodeIndex int
}

func (d *provider) GetPlotName() string { return "AnyShake Observer" }
func (d *provider) GetStation() string  { return d.stationCode }
func (d *provider) GetNetwork() string  { return d.networkCode }
func (d *provider) GetChannel() string  { return d.channelCode }
func (d *provider) GetLocation() string { return d.locationCode }
func (d *provider) GetPlotData(startTime, endTime time.Time) ([]heligo.PlotData, error) {
	startTimestamp := startTime.Add(-time.Second) // Also used as key of cache
	endTimestamp := endTime.Add(time.Second)
	cacheData, found := d.queryCache.Get(startTimestamp)

	var seisRecords []queryCacheData
	if found {
		seisRecords = cacheData
	} else {
		records, err := d.actionHandler.SeisRecordsQuery(startTimestamp, endTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to query seismic waveform records: %w", err)
		}
		seisRecords = make([]queryCacheData, len(records))
		for idx, record := range records {
			_, sampleRate, channelData, err := record.Decode()
			if err != nil {
				return nil, fmt.Errorf("failed to decode seismic waveform record on timestamp %d: %w", record.Timestamp, err)
			}
			seisRecords[idx].sampleRate = sampleRate
			seisRecords[idx].channelData = channelData
			seisRecords[idx].timestamp = record.Timestamp
		}
		d.queryCache.Set(startTimestamp, seisRecords)
	}

	var plotData []heligo.PlotData
	for _, record := range seisRecords {
		data := make([]heligo.PlotData, record.sampleRate)
		for i := 0; i < record.sampleRate; i++ {
			timeOffset := int64(i * 1000 / record.sampleRate)
			data[i].Time = time.UnixMilli(record.timestamp + timeOffset)
			if d.channelCodeIndex > len(record.channelData) {
				continue
			}
			data[i].Value = float64(record.channelData[d.channelCodeIndex].Data[i])
		}
		plotData = append(plotData, data...)
	}

	return plotData, nil
}
func (d *provider) setChannelCode(channelCode string, channelCodeIndex int) {
	d.channelCode = channelCode
	d.channelCodeIndex = channelCodeIndex
}

func (s *HelicorderServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	logger.GetLogger(ID).Infof("generated helicorder images will be saved to %s", s.filePath)

	go func() {
		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)

		timer := time.NewTimer(time.Minute)

		for {
			currentTime := s.timeSource.Get().Add(-time.Minute) // Subtract one minute to avoid date rollover
			timer.Reset(s.getDurationToNextTime(s.timeSource.Get()))

			select {
			case <-s.ctx.Done():
				timer.Stop()
				s.wg.Done()
				return
			case <-timer.C:
				hardwareConfig := s.hardwareDev.GetConfig()
				s.channelCodes = hardwareConfig.GetChannelCodes()

				for channelIdx, channelCode := range s.channelCodes {
					helicorderCtx, err := heligo.New(&s.dataProvider, 24*time.Hour, time.Duration(s.timeSpan)*time.Minute)
					if err != nil {
						logger.GetLogger(ID).Errorf("failed to create helicorder context: %v", err)
						continue
					}

					// Update current channel code
					s.dataProvider.setChannelCode(channelCode, channelIdx)
					logger.GetLogger(ID).Infof("start plotting helicorder for channel %s", channelCode)

					err = helicorderCtx.Plot(currentTime, s.spanSamples, s.scaleFactor, s.lineWidth)
					if err != nil {
						logger.GetLogger(ID).Errorf("failed to plot helicorder for %s: %v", channelCode, err)
						continue
					}

					dateDir := filepath.Join(s.filePath, currentTime.UTC().Format("2006-01-02"))
					filePath := filepath.Join(dateDir, s.getHelicorderFileName(currentTime, channelCode))

					if err := os.MkdirAll(dateDir, 0755); err != nil {
						logger.GetLogger(ID).Errorf("failed to create directory %s: %v", dateDir, err)
						continue
					}

					err = helicorderCtx.Save(s.imageSize, filePath)
					if err != nil {
						logger.GetLogger(ID).Errorf("failed to save helicorder for %s: %v", channelCode, err)
						continue
					}

					logger.GetLogger(ID).Infof("helicorder for %s has been saved to %s", channelCode, filePath)
				}

				s.dataProvider.queryCache.Clear()

				if s.lifeCycle > 0 {
					endTime := currentTime.Add(time.Duration(-s.lifeCycle) * time.Hour * 24)
					if err := s.cleanupHelicorderFiles(endTime); err != nil {
						logger.GetLogger(ID).Errorf("failed to purge expired helicorder files: %v", err)
					}
				}
			}
		}
	}()

	s.wg.Add(1)
	return nil
}

func (s *HelicorderServiceImpl) getDurationToNextTime(currentTime time.Time) time.Duration {
	timsSpanMinute := int(time.Hour.Minutes())
	currentMinute := currentTime.Minute()
	// Minutes to next time span
	nextQuarter := (currentMinute/timsSpanMinute + 1) * timsSpanMinute % 60
	nextTime := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		nextQuarter,
		0, // Reset seconds
		0,
		currentTime.Location(),
	)
	if nextQuarter <= currentMinute {
		nextTime = nextTime.Add(time.Hour)
	}
	return nextTime.Sub(currentTime)
}

func (m *HelicorderServiceImpl) getHelicorderFileName(tm time.Time, channelCode string) string {
	return fmt.Sprintf("%s.%s.%s.%s.%s.%s.%s",
		m.dataProvider.networkCode, m.dataProvider.stationCode, m.dataProvider.locationCode, channelCode,
		tm.UTC().Format("2006"),
		tm.UTC().Format("002"),
		m.imageFormat,
	)
}

func (s *HelicorderServiceImpl) cleanupHelicorderFiles(until time.Time) error {
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
