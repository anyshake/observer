package seedlink

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/slgo"
	"github.com/bclswl0827/slgo/handlers"
)

func (s *SeedLinkServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	seedlinkMessageBus := message.NewBus[explorer.EventHandler](ID, 65535)
	server := slgo.New(
		&provider{
			hardwareDev:   s.hardwareDev,
			timeSource:    s.timeSource,
			actionHandler: s.actionHandler,
			startTime:     s.timeSource.Get(),
			stationCode:   s.stationCode,
			networkCode:   s.networkCode,
			locationCode:  s.locationCode,
		},
		&consumer{
			messageBus: seedlinkMessageBus,
		},
		&hooks{},
	)

	go func() {
		err := s.hardwareDev.Subscribe(ID, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			seedlinkMessageBus.Publish(t, di, dv, cd)
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)

		logger.GetLogger(ID).Infof("service seedlink is listening on %s:%d", s.listenHost, s.listenPort)
		if err := server.Start(s.ctx, s.listenHost, s.listenPort, s.useCompress); err != nil {
			logger.GetLogger(ID).Errorf("failed to start seedlink server: %v", err)
			s.status.SetStoppedAt(s.timeSource.Get())
			s.status.SetIsRunning(false)
			_ = s.hardwareDev.Unsubscribe(ID)
		}

		s.wg.Done()
	}()

	s.wg.Add(1)
	return nil
}

type provider struct {
	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler
	startTime     time.Time
	stationCode   string
	networkCode   string
	locationCode  string
}

func (p *provider) GetSoftware() string       { return "anyshake_observer" }
func (p *provider) GetOrganization() string   { return "anyshake.org" }
func (p *provider) GetCurrentTime() time.Time { return p.timeSource.Get() }
func (p *provider) GetStartTime() time.Time   { return p.startTime }
func (p *provider) GetCapabilities() []handlers.SeedLinkCapability {
	return []handlers.SeedLinkCapability{
		{Name: "info:all"}, {Name: "info:gaps"}, {Name: "info:streams"},
		{Name: "dialup"}, {Name: "info:id"}, {Name: "multistation"},
		{Name: "window-extraction"}, {Name: "info:connections"},
		{Name: "info:capabilities"}, {Name: "info:stations"},
	}
}
func (p *provider) GetStations() []handlers.SeedLinkStation {
	return []handlers.SeedLinkStation{
		{
			BeginSequence: "000000",
			EndSequence:   "FFFFFF",
			Station:       p.stationCode,
			Network:       p.networkCode,
			Description:   "AnyShake Observer SeedLink Service",
		},
	}
}
func (p *provider) GetStreams() []handlers.SeedLinkStream {
	hardwareCfg := p.hardwareDev.GetConfig()
	channelCodes := hardwareCfg.GetChannelCodes()

	streams := make([]handlers.SeedLinkStream, len(channelCodes))
	for idx, channelCode := range channelCodes {
		streams[idx] = handlers.SeedLinkStream{
			BeginTime: p.GetStartTime().Format("2006-01-02 15:04:05"),
			EndTime:   p.GetCurrentTime().Format("2006-01-02 15:04:05"),
			SeedName:  channelCode,
			Location:  p.locationCode,
			Station:   p.stationCode,
			Type:      "D",
		}
	}

	return streams
}
func (p *provider) QueryHistory(startTime, endTime time.Time, channels []handlers.SeedLinkChannel) ([]handlers.SeedLinkDataPacket, error) {
	if endTime.IsZero() {
		endTime = p.timeSource.Get()
	}
	recordsRawData, err := p.actionHandler.SeisRecordsQuery(startTime, endTime)
	if err != nil {
		return nil, err
	}

	channelSet := make(map[string]struct{}, len(channels))
	for _, ch := range channels {
		channelSet[ch.ChannelName] = struct{}{}
	}

	var dataPackets []handlers.SeedLinkDataPacket
	for _, record := range recordsRawData {
		tm, sampleRate, channelData, err := record.Decode()
		if err != nil {
			return nil, err
		}

		for _, data := range channelData {
			if _, exists := channelSet[data.ChannelCode]; exists {
				dataPackets = append(dataPackets, handlers.SeedLinkDataPacket{
					Timestamp:  tm.UnixMilli(),
					SampleRate: sampleRate,
					Channel:    data.ChannelCode,
					DataArr:    data.Data,
				})
			}
		}
	}

	return dataPackets, nil
}

type consumer struct {
	messageBus message.Bus[explorer.EventHandler]
}

func (c *consumer) Subscribe(clientId string, channels []handlers.SeedLinkChannel, eventHandler func(handlers.SeedLinkDataPacket)) error {
	channelSet := make(map[string]struct{}, len(channels))
	for _, ch := range channels {
		channelSet[ch.ChannelName] = struct{}{}
	}

	handler := func(tm time.Time, dc *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
		sampleRate := dc.GetSampleRate()
		for _, data := range cd {
			if _, exists := channelSet[data.ChannelCode]; exists {
				eventHandler(handlers.SeedLinkDataPacket{
					Timestamp:  tm.UnixMilli(),
					SampleRate: sampleRate,
					Channel:    data.ChannelCode,
					DataArr:    data.Data,
				})
			}
		}
	}

	return c.messageBus.Subscribe(clientId, handler)
}
func (c *consumer) Unsubscribe(clientId string) error {
	return c.messageBus.Unsubscribe(clientId)
}

type hooks struct{}

func (h *hooks) OnData(client *handlers.SeedLinkClient, data []byte) {}
func (h *hooks) OnConnection(client *handlers.SeedLinkClient) {
	logger.GetLogger(ID).Infof("%s - client connected to SeedLink service", client.RemoteAddr().String())
}
func (h *hooks) OnClose(client *handlers.SeedLinkClient) {
	logger.GetLogger(ID).Infof("%s - client disconnected from SeedLink service", client.RemoteAddr().String())
}
func (h *hooks) OnCommand(client *handlers.SeedLinkClient, command []string) {
	logger.GetLogger(ID).Infof("%s - client sent command to SeedLink service: %s", client.RemoteAddr().String(), command)
}
