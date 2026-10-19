package seedlink

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/slgo"
	"github.com/bclswl0827/slgo/handlers"
)

const seedLinkWriteTimeout = 5 * time.Second

// Hardware publishes one event per second; keep enough to bridge batched
// archiver writes and preserve the recent tail while history is queried.
const seedLinkHistoryBufferSize = 120

func (s *SeedLinkServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("SeedLink service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	runCtx, runCancel := s.ctx, s.cancelFn
	seedlinkMessageBus := message.NewBus[explorer.Event](ID)
	recentBuffer := ringbuf.New[explorer.Event](seedLinkHistoryBufferSize)
	clients := newSeedLinkClientRegistry()
	server := slgo.New(
		&provider{
			ctx:           runCtx,
			hardwareDev:   s.hardwareDev,
			timeSource:    s.timeSource,
			actionHandler: s.actionHandler,
			recentBuffer:  recentBuffer,
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

	compression := slgo.CompressionNone
	if s.useCompress {
		switch s.compressScheme {
		case COMPRESS_SCHEME_STEIM_2:
			compression = slgo.CompressionSteim2
		default:
			compression = slgo.CompressionSteim1
		}
	}
	listenHost, listenPort := s.listenHost, s.listenPort

	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			runCancel()
			seedlinkMessageBus.Close()
			_ = s.hardwareDev.Unsubscribe(ID)
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
			defer service.RecoverCallback(ID, runCancel)
			// Cache before publishing so a history query can always see events
			// emitted before the client's live subscription was installed.
			recentBuffer.Push(event)
			seedlinkMessageBus.Publish(event)
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		logger.GetLogger(ID).Infof("service seedlink is listening on %s:%d", listenHost, listenPort)
		if err := server.Start(runCtx, listenHost, listenPort, compression); err != nil {
			logger.GetLogger(ID).Errorf("failed to start seedlink server: %v", err)
		}
	}()

	return nil
}

func (s *SeedLinkServiceImpl) GetListenPort() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.listenPort
}

type provider struct {
	ctx           context.Context
	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler
	recentBuffer  *ringbuf.Buffer[explorer.Event]
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
		if _, ok := c.clients.Get(clientId); !ok {
			return
		}
		for _, data := range event.ChannelData {
			if _, exists := channelSet[data.ChannelCode]; exists {
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

func (h *hooks) OnData(_ *handlers.SeedLinkClient, _ []byte) {}
func (h *hooks) OnConnection(client *handlers.SeedLinkClient) {
	client.Conn = &seedLinkConn{Conn: client.Conn}
	h.clients.Add(client)
	logger.GetLogger(ID).Infof("%s - client connected to SeedLink service", client.RemoteAddr().String())
}
func (h *hooks) OnClose(client *handlers.SeedLinkClient) {
	h.clients.Remove(client)
	logger.GetLogger(ID).Infof("%s - client disconnected from SeedLink service", client.RemoteAddr().String())
}
func (h *hooks) OnCommand(client *handlers.SeedLinkClient, command []string) {
	logger.GetLogger(ID).Infof("%s - client sent command to SeedLink service: %s", client.RemoteAddr().String(), command)
}
