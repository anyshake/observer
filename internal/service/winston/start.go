package winston

import (
	"context"
	"errors"
	"fmt"
	"math"
	"runtime/debug"
	"sort"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/bclswl0827/winsgo"
)

type waveformStore interface {
	SeisRecordsGetQueryWindow() time.Duration
	SeisRecordsQuery(time.Time, time.Time) ([]model.SeisRecord, error)
}

type waveformRecord struct {
	timestamp  time.Time
	sampleRate int
	samples    []int32
}

func cloneChannelData(data []explorer.ChannelData) []explorer.ChannelData {
	cloned := make([]explorer.ChannelData, len(data))
	for i := range data {
		cloned[i] = data[i]
		cloned[i].Data = append([]int32(nil), data[i].Data...)
	}
	return cloned
}

func (s *WinstonServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("Winston service is already running")
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}
	if s.ringBuffer == nil {
		return errors.New("Winston service is not initialized")
	}

	server := winsgo.New(
		&provider{
			hardwareDev:  s.hardwareDev,
			timeSource:   s.timeSource.Now,
			stationCode:  s.stationCode,
			networkCode:  s.networkCode,
			locationCode: s.locationCode,
		},
		&consumer{
			stationCode:  s.stationCode,
			networkCode:  s.networkCode,
			locationCode: s.locationCode,
			ringBuffer:   s.ringBuffer,
			store:        s.actionHandler,
		},
		&hooks{},
	)
	server.DisableHelicorder = true

	err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
		s.ringBuffer.Push(winstonRingBuffer{
			timestamp:   event.Timestamp,
			sampleRate:  event.SampleRate,
			channelData: cloneChannelData(event.ChannelData),
		})
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to hardware message bus: %w", err)
	}

	s.wg.Add(1)
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	go func() {
		defer s.wg.Done()
		defer func() {
			_ = s.hardwareDev.Unsubscribe(ID)
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		logger.GetLogger(ID).Infof("service Winston is listening on %s:%d", s.listenHost, s.listenPort)
		if err := server.Start(s.ctx, s.listenHost, s.listenPort); err != nil {
			logger.GetLogger(ID).Errorf("failed to start Winston server: %v", err)
		}
	}()

	return nil
}

type provider struct {
	hardwareDev  hardware.IHardware
	timeSource   func() time.Time
	stationCode  string
	networkCode  string
	locationCode string
}

func (p *provider) Channels(ctx context.Context) ([]winsgo.Channel, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	deviceConfig := p.hardwareDev.GetConfig()
	channelCodes := deviceConfig.GetChannelCodes()
	latitude, longitude, elevation, _ := p.hardwareDev.GetCoordinates(false)
	now := p.timeSource().UTC()
	instrument := winsgo.Instrument{
		ID:          1,
		Name:        deviceConfig.GetModel(),
		Description: "Powered by AnyShake Project",
		Longitude:   longitude,
		Latitude:    latitude,
		Height:      elevation,
		TimeZone:    time.Local.String(),
	}

	channels := make([]winsgo.Channel, len(channelCodes))
	for i, channelCode := range channelCodes {
		channels[i] = winsgo.Channel{
			ID: i + 1,
			SCNL: winsgo.SCNL{
				Station:  p.stationCode,
				Channel:  channelCode,
				Network:  p.networkCode,
				Location: p.locationCode,
			},
			StartTime:  time.Unix(0, 0).UTC(),
			EndTime:    now,
			Instrument: instrument,
			Alias:      channelCode,
			Unit:       "count",
			LinearA:    1,
			Groups:     []string{ID},
		}
	}
	return channels, nil
}

type consumer struct {
	stationCode  string
	networkCode  string
	locationCode string
	ringBuffer   *ringbuf.Buffer[winstonRingBuffer]
	store        waveformStore
}

func (c *consumer) Consume(ctx context.Context, request winsgo.WaveformRequest, deliver winsgo.WaveformHandler) error {
	if request.Channel.Station != c.stationCode ||
		request.Channel.Network != c.networkCode ||
		request.Channel.Location != c.locationCode ||
		request.Channel.Channel == "" ||
		!request.EndTime.After(request.StartTime) {
		return winsgo.ErrNoData
	}
	if c.store == nil {
		return errors.New("Winston waveform store is not configured")
	}
	queryWindow := c.store.SeisRecordsGetQueryWindow()
	if request.EndTime.Sub(request.StartTime) > queryWindow {
		return fmt.Errorf("Winston query duration exceeds %s limit", queryWindow)
	}

	ringRecords, err := c.queryRing(ctx, request)
	if err != nil {
		return err
	}
	records := ringRecords
	if !recordsCover(ringRecords, request.StartTime, request.EndTime) {
		daoRecords, err := c.queryDAO(ctx, request)
		if err != nil {
			return err
		}
		records = mergeRecords(daoRecords, ringRecords)
	}

	waveform, ok, err := buildWaveform(request.StartTime, request.EndTime, records)
	if err != nil {
		return err
	}
	if !ok {
		return winsgo.ErrNoData
	}
	return deliver(waveform)
}

func (c *consumer) queryRing(ctx context.Context, request winsgo.WaveformRequest) ([]waveformRecord, error) {
	if c.ringBuffer == nil {
		return nil, nil
	}

	values := c.ringBuffer.Values()
	records := make([]waveformRecord, 0, len(values))
	for _, value := range values {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		for _, channel := range value.channelData {
			if channel.ChannelCode != request.Channel.Channel || value.sampleRate <= 0 || len(channel.Data) == 0 {
				continue
			}
			recordEnd := value.timestamp.Add(samplesDuration(len(channel.Data), value.sampleRate))
			if value.timestamp.Before(request.EndTime) && recordEnd.After(request.StartTime) {
				records = append(records, waveformRecord{
					timestamp:  value.timestamp,
					sampleRate: value.sampleRate,
					samples:    append([]int32(nil), channel.Data...),
				})
			}
		}
	}
	sortRecords(records)
	return records, nil
}

func (c *consumer) queryDAO(ctx context.Context, request winsgo.WaveformRequest) ([]waveformRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	daoRecords, err := c.store.SeisRecordsQuery(request.StartTime, request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("query Winston waveform history: %w", err)
	}

	var records []waveformRecord
	for i := range daoRecords {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		timestamp, sampleRate, channelData, err := daoRecords[i].Decode()
		if err != nil {
			return nil, fmt.Errorf("decode Winston waveform history: %w", err)
		}
		for _, channel := range channelData {
			if channel.ChannelCode != request.Channel.Channel || sampleRate <= 0 || len(channel.Data) == 0 {
				continue
			}
			recordEnd := timestamp.Add(samplesDuration(len(channel.Data), sampleRate))
			if timestamp.Before(request.EndTime) && recordEnd.After(request.StartTime) {
				records = append(records, waveformRecord{
					timestamp:  timestamp,
					sampleRate: sampleRate,
					samples:    append([]int32(nil), channel.Data...),
				})
			}
		}
	}
	sortRecords(records)
	return records, nil
}

func recordsCover(records []waveformRecord, start, end time.Time) bool {
	if len(records) == 0 || !end.After(start) {
		return false
	}

	coveredUntil := start
	for _, record := range records {
		if record.sampleRate <= 0 || len(record.samples) == 0 {
			continue
		}
		recordEnd := record.timestamp.Add(samplesDuration(len(record.samples), record.sampleRate))
		if !recordEnd.After(coveredUntil) {
			continue
		}
		if record.timestamp.After(coveredUntil) {
			return false
		}
		coveredUntil = recordEnd
		if !coveredUntil.Before(end) {
			return true
		}
	}
	return false
}

func mergeRecords(daoRecords, ringRecords []waveformRecord) []waveformRecord {
	merged := make(map[int64]waveformRecord, len(daoRecords)+len(ringRecords))
	for _, record := range daoRecords {
		merged[record.timestamp.UnixNano()] = record
	}
	// Ring data is newer and avoids the archiver's write delay, so it wins on
	// duplicate timestamps.
	for _, record := range ringRecords {
		merged[record.timestamp.UnixNano()] = record
	}

	records := make([]waveformRecord, 0, len(merged))
	for _, record := range merged {
		records = append(records, record)
	}
	sortRecords(records)
	return records
}

func sortRecords(records []waveformRecord) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].timestamp.Before(records[j].timestamp)
	})
}

func samplesDuration(count, sampleRate int) time.Duration {
	return time.Duration(float64(count) / float64(sampleRate) * float64(time.Second))
}

func sampleOffset(index, sampleRate int) time.Duration {
	return time.Duration(float64(index) / float64(sampleRate) * float64(time.Second))
}

func buildWaveform(start, end time.Time, records []waveformRecord) (winsgo.Waveform, bool, error) {
	if len(records) == 0 || !end.After(start) {
		return winsgo.Waveform{}, false, nil
	}

	sortRecords(records)
	sampleRate := records[0].sampleRate
	if sampleRate <= 0 {
		return winsgo.Waveform{}, false, errors.New("invalid Winston sample rate")
	}

	type clippedRecord struct {
		record waveformRecord
		first  int
		last   int
	}
	clipped := make([]clippedRecord, 0, len(records))
	var waveformStart, waveformEnd time.Time
	for _, record := range records {
		if record.sampleRate != sampleRate {
			return winsgo.Waveform{}, false, errors.New("sample rate changed within Winston query range")
		}
		if len(record.samples) == 0 {
			continue
		}

		first := 0
		if start.After(record.timestamp) {
			first = int(math.Ceil(start.Sub(record.timestamp).Seconds()*float64(sampleRate) - 1e-9))
		}
		last := len(record.samples)
		if end.Before(record.timestamp.Add(samplesDuration(len(record.samples), sampleRate))) {
			last = int(math.Ceil(end.Sub(record.timestamp).Seconds()*float64(sampleRate) - 1e-9))
		}
		first = max(0, first)
		last = min(len(record.samples), last)
		if first >= last {
			continue
		}

		recordStart := record.timestamp.Add(sampleOffset(first, sampleRate))
		recordEnd := record.timestamp.Add(sampleOffset(last, sampleRate))
		if waveformStart.IsZero() || recordStart.Before(waveformStart) {
			waveformStart = recordStart
		}
		if recordEnd.After(waveformEnd) {
			waveformEnd = recordEnd
		}
		clipped = append(clipped, clippedRecord{record: record, first: first, last: last})
	}
	if len(clipped) == 0 {
		return winsgo.Waveform{}, false, nil
	}

	count := int(math.Ceil(waveformEnd.Sub(waveformStart).Seconds()*float64(sampleRate) - 1e-9))
	samples := make([]int32, count)
	for i := range samples {
		samples[i] = math.MinInt32
	}
	for _, item := range clipped {
		recordStart := item.record.timestamp.Add(sampleOffset(item.first, sampleRate))
		offset := int(math.Round(recordStart.Sub(waveformStart).Seconds() * float64(sampleRate)))
		for source := item.first; source < item.last; source++ {
			destination := offset + source - item.first
			if destination >= 0 && destination < len(samples) {
				samples[destination] = item.record.samples[source]
			}
		}
	}

	return winsgo.Waveform{
		StartTime:          waveformStart.UTC(),
		SampleRate:         float64(sampleRate),
		RegistrationOffset: math.NaN(),
		Samples:            samples,
		DataType:           "s4",
	}, true, nil
}

type hooks struct{}

func (*hooks) OnData(_ *winsgo.Client, _ []byte) {}
func (*hooks) OnConnection(client *winsgo.Client) {
	logger.GetLogger(ID).Infof("%s - client connected to Winston service", client.RemoteAddr())
}
func (*hooks) OnClose(client *winsgo.Client) {
	logger.GetLogger(ID).Infof("%s - client disconnected from Winston service", client.RemoteAddr())
}
func (*hooks) OnCommand(client *winsgo.Client, command winsgo.Command) {
	logger.GetLogger(ID).Infof("%s - client sent command to Winston service: %s", client.RemoteAddr(), command.Raw)
}

func (s *WinstonServiceImpl) GetListenPort() int {
	return s.listenPort
}
