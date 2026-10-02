package seedlink

import (
	"context"
	"runtime/debug"
	"sync"
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

const seedLinkWriteTimeout = 5 * time.Second

func (s *SeedLinkServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	seedlinkMessageBus := message.NewBus[explorer.Event](ID)
	clients := newSeedLinkClientRegistry()
	server := slgo.New(
		&provider{
			hardwareDev:   s.hardwareDev,
			timeSource:    s.timeSource,
			actionHandler: s.actionHandler,
			startTime:     s.timeSource.Now(),
			stationCode:   s.stationCode,
			networkCode:   s.networkCode,
			locationCode:  s.locationCode,
		},
		&consumer{
			messageBus: seedlinkMessageBus,
			clients:    clients,
		},
		&hooks{clients: clients},
	)

	go func() {
		s.status.SetStartedAt(s.timeSource.Now())
		s.status.SetIsRunning(true)
		defer func() {
			seedlinkMessageBus.Close()
			_ = s.hardwareDev.Unsubscribe(ID)
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
			s.wg.Done()
		}()

		err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
			seedlinkMessageBus.Publish(event)
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		logger.GetLogger(ID).Infof("service seedlink is listening on %s:%d", s.listenHost, s.listenPort)
		if err := server.Start(s.ctx, s.listenHost, s.listenPort, s.useCompress); err != nil {
			logger.GetLogger(ID).Errorf("failed to start seedlink server: %v", err)
		}
	}()

	s.wg.Add(1)
	return nil
}

func (s *SeedLinkServiceImpl) GetListenPort() int {
	return s.listenPort
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
func (p *provider) GetCurrentTime() time.Time { return p.timeSource.Now() }
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
		endTime = p.timeSource.Now()
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
	messageBus *message.Bus[explorer.Event]
	clients    *seedLinkClientRegistry
}

func (c *consumer) Subscribe(clientId string, channels []handlers.SeedLinkChannel, eventHandler func(handlers.SeedLinkDataPacket)) error {
	channelSet := make(map[string]struct{}, len(channels))
	for _, ch := range channels {
		channelSet[ch.ChannelName] = struct{}{}
	}

	handler := func(event explorer.Event) {
		client, ok := c.clients.Get(clientId)
		if !ok {
			return
		}
		for _, data := range event.ChannelData {
			if _, exists := channelSet[data.ChannelCode]; exists {
				if err := client.SetWriteDeadline(time.Now().Add(seedLinkWriteTimeout)); err != nil {
					_ = client.Close()
					return
				}
				eventHandler(handlers.SeedLinkDataPacket{
					Timestamp:  event.Timestamp.UnixMilli(),
					SampleRate: event.SampleRate,
					Channel:    data.ChannelCode,
					DataArr:    data.Data,
				})
			}
		}
	}

	return c.messageBus.Subscribe(clientId, message.SubscriptionOptions{
		BufferSize: 16,
		Overflow:   message.OverflowDisconnect,
		OnError: func(err error) {
			logger.GetLogger(ID).Warnf("%s - SeedLink subscription closed: %v", clientId, err)
			if client, ok := c.clients.Get(clientId); ok {
				_ = client.Close()
			}
		},
	}, handler)
}
func (c *consumer) Unsubscribe(clientId string) error {
	return c.messageBus.Unsubscribe(clientId)
}

type seedLinkClientRegistry struct {
	mu      sync.RWMutex
	clients map[string]*handlers.SeedLinkClient
}

func newSeedLinkClientRegistry() *seedLinkClientRegistry {
	return &seedLinkClientRegistry{clients: make(map[string]*handlers.SeedLinkClient)}
}

func (r *seedLinkClientRegistry) Add(client *handlers.SeedLinkClient) {
	r.mu.Lock()
	r.clients[client.RemoteAddr().String()] = client
	r.mu.Unlock()
}

func (r *seedLinkClientRegistry) Remove(client *handlers.SeedLinkClient) {
	clientID := client.RemoteAddr().String()
	r.mu.Lock()
	if current, ok := r.clients[clientID]; ok && current == client {
		delete(r.clients, clientID)
	}
	r.mu.Unlock()
}

func (r *seedLinkClientRegistry) Get(clientID string) (*handlers.SeedLinkClient, bool) {
	r.mu.RLock()
	client, ok := r.clients[clientID]
	r.mu.RUnlock()
	return client, ok
}

type hooks struct {
	clients *seedLinkClientRegistry
}

func (h *hooks) OnData(client *handlers.SeedLinkClient, _ []byte) {
	_ = client.SetWriteDeadline(time.Now().Add(seedLinkWriteTimeout))
}
func (h *hooks) OnConnection(client *handlers.SeedLinkClient) {
	h.clients.Add(client)
	logger.GetLogger(ID).Infof("%s - client connected to SeedLink service", client.RemoteAddr().String())
}
func (h *hooks) OnClose(client *handlers.SeedLinkClient) {
	h.clients.Remove(client)
	logger.GetLogger(ID).Infof("%s - client disconnected from SeedLink service", client.RemoteAddr().String())
}
func (h *hooks) OnCommand(client *handlers.SeedLinkClient, command []string) {
	_ = client.SetWriteDeadline(time.Now().Add(seedLinkWriteTimeout))
	logger.GetLogger(ID).Infof("%s - client sent command to SeedLink service: %s", client.RemoteAddr().String(), command)
}
