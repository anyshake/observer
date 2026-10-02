package tailscale

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"testing"

	"github.com/anyshake/observer/pkg/timesource"
	"tailscale.com/ipn"
)

type memoryStateStore struct {
	mu     sync.Mutex
	values map[ipn.StateKey][]byte
}

func newMemoryStateStore() *memoryStateStore {
	return &memoryStateStore{values: make(map[ipn.StateKey][]byte)}
}

func (s *memoryStateStore) ReadState(id ipn.StateKey) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.values[id]
	if !ok {
		return nil, ipn.ErrStateNotExist
	}
	return bytes.Clone(value), nil
}

func (s *memoryStateStore) WriteState(id ipn.StateKey, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value == nil {
		delete(s.values, id)
		return nil
	}
	s.values[id] = bytes.Clone(value)
	return nil
}

type testLogger struct {
	mu       sync.Mutex
	messages []string
}

func (l *testLogger) Infof(format string, args ...any) {
	l.append(format, args...)
}

func (l *testLogger) Errorf(format string, args ...any) {
	l.append(format, args...)
}

func (l *testLogger) append(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf(format, args...))
}

func (l *testLogger) contains(value string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Contains(strings.Join(l.messages, "\n"), value)
}

type fakeNode struct {
	upStarted chan struct{}
	closed    chan struct{}
	startOnce sync.Once
	closeOnce sync.Once
	listener  *blockingListener
}

type blockingListener struct {
	closed    chan struct{}
	closeOnce sync.Once
}

func newBlockingListener() *blockingListener {
	return &blockingListener{closed: make(chan struct{})}
}

func (l *blockingListener) Accept() (net.Conn, error) {
	<-l.closed
	return nil, net.ErrClosed
}

func (l *blockingListener) Close() error {
	l.closeOnce.Do(func() { close(l.closed) })
	return nil
}

func (l *blockingListener) Addr() net.Addr {
	return fakeAddr("100.64.0.1:8073")
}

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeNode() *fakeNode {
	return &fakeNode{
		upStarted: make(chan struct{}),
		closed:    make(chan struct{}),
		listener:  newBlockingListener(),
	}
}

func (n *fakeNode) Start() error {
	return nil
}

func (n *fakeNode) Up(ctx context.Context) error {
	n.startOnce.Do(func() { close(n.upStarted) })
	return nil
}

func (n *fakeNode) Listen(_, _ string) (net.Listener, error) {
	return n.listener, nil
}

func (n *fakeNode) ListenPacket(_, _ string) (net.PacketConn, error) {
	return nil, errors.New("unexpected UDP listener")
}

func (n *fakeNode) TailscaleIPs() (netip.Addr, netip.Addr) {
	return netip.MustParseAddr("100.64.0.1"), netip.Addr{}
}

func (n *fakeNode) Close() error {
	_ = n.listener.Close()
	n.closeOnce.Do(func() { close(n.closed) })
	return nil
}

func newTestService(store ipn.StateStore, log serviceLogger) *TailscaleServiceImpl {
	service := New("0.0.0.0:8073", nil, timesource.New(nil))
	service.stateStore = store
	service.log = log
	service.hostname = "observer-test"
	service.controlURL = "https://headscale.example.com"
	return service
}

func TestStartRequiresAuthKeyForInitialRegistration(t *testing.T) {
	service := newTestService(newMemoryStateStore(), &testLogger{})
	service.newNode = func(nodeConfig) embeddedNode {
		t.Fatal("node must not be created without an auth key or persisted identity")
		return nil
	}

	err := service.Start()
	if err == nil || !strings.Contains(err.Error(), "auth key is required") {
		t.Fatalf("unexpected start error: %v", err)
	}
}

func TestStartUsesAuthKeyOnlyForInitialRegistration(t *testing.T) {
	const authKey = "tskey-auth-secret"
	store := newMemoryStateStore()
	log := &testLogger{}
	service := newTestService(store, log)
	service.authKey = authKey
	node := newFakeNode()
	service.newNode = func(config nodeConfig) embeddedNode {
		if config.AuthKey == "" {
			t.Fatal("auth key was not provided for initial registration")
		}
		config.UserLogf("auth key: %s", config.AuthKey)
		return node
	}

	if err := service.Start(); err != nil {
		t.Fatalf("failed to start service: %v", err)
	}
	<-node.upStarted
	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service: %v", err)
	}
	<-node.closed

	if log.contains(authKey) {
		t.Fatal("auth key was written to the service logger")
	}
	if !log.contains("[redacted]") {
		t.Fatal("redacted auth key marker was not logged")
	}
}

func TestStartRestoresPersistedIdentityWithoutAuthKey(t *testing.T) {
	store := newMemoryStateStore()
	if err := store.WriteState(ipn.CurrentProfileStateKey, []byte("profile-key")); err != nil {
		t.Fatalf("failed to seed state store: %v", err)
	}

	service := newTestService(store, &testLogger{})
	service.authKey = "unused-auth-key"
	node := newFakeNode()
	service.newNode = func(config nodeConfig) embeddedNode {
		if config.AuthKey != "" {
			t.Fatal("auth key must not be used when a persisted identity exists")
		}
		if config.Store != store {
			t.Fatal("database state store was not passed to the embedded node")
		}
		return node
	}

	if err := service.Start(); err != nil {
		t.Fatalf("failed to start service with persisted identity: %v", err)
	}
	<-node.upStarted
	if !service.GetStatus().GetIsRunning() {
		t.Fatal("service did not enter running state")
	}

	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service: %v", err)
	}
	<-node.closed
	if service.GetStatus().GetIsRunning() {
		t.Fatal("service remained running after context cancellation")
	}
}
