package dnsquery

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestNewServerTypesAndDefaultPorts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		endpoint string
		check    func(IServer) bool
	}{
		{"https://dns.example.test/dns-query", func(server IServer) bool {
			value, ok := server.(*DoH)
			return ok && value.server == "https://dns.example.test/dns-query"
		}},
		{"tls://dns.example.test", func(server IServer) bool {
			value, ok := server.(*DoT)
			return ok && value.server == "dns.example.test:853"
		}},
		{"tls://dns.example.test:8853", func(server IServer) bool {
			value, ok := server.(*DoT)
			return ok && value.server == "dns.example.test:8853"
		}},
		{"udp://8.8.8.8", func(server IServer) bool { value, ok := server.(*UDP); return ok && value.server == "8.8.8.8:53" }},
		{"udp://[2001:4860:4860::8888]", func(server IServer) bool {
			value, ok := server.(*UDP)
			return ok && value.server == "[2001:4860:4860::8888]:53"
		}},
		{"sdns://AQcAAAAAAAAAEzE0OS4xMTIuMTEyLjExMjoyNDQz", func(server IServer) bool { _, ok := server.(*DNSCrypt); return ok }},
	}
	for _, tt := range tests {
		server, err := New(tt.endpoint)
		if err != nil || !tt.check(server) {
			t.Errorf("New(%q) = %#v, %v", tt.endpoint, server, err)
		}
	}
}

func TestNewServerRejectsUnsupportedEndpoints(t *testing.T) {
	t.Parallel()
	for _, endpoint := range []string{"", "://bad", "tcp://dns.example.test", "dns.example.test"} {
		if server, err := New(endpoint); err == nil || server != nil {
			t.Errorf("New(%q) = %#v, %v, want error", endpoint, server, err)
		}
	}
}

func TestResolverLifecycle(t *testing.T) {
	t.Parallel()
	query := new(dns.Msg).SetQuestion("example.test.", dns.TypeA)
	for name, server := range map[string]IServer{
		"UDP": &UDP{}, "DoT": &DoT{}, "DoH": &DoH{}, "DNSCrypt": &DNSCrypt{},
	} {
		if _, err := server.Query(query, time.Millisecond); err == nil {
			t.Errorf("%s Query() succeeded before Open()", name)
		}
		if err := server.Close(); err != nil {
			t.Errorf("%s Close() error = %v", name, err)
		}
	}
	for _, server := range []IServer{&UDP{}, &DoT{}, &DoH{}} {
		if err := server.Open(); err != nil {
			t.Fatal(err)
		}
		if err := server.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err := server.Query(query, time.Millisecond); err == nil {
			t.Errorf("%T Query() succeeded after Close()", server)
		}
	}
}

func TestResolverCatalogAndPickRandom(t *testing.T) {
	t.Parallel()
	resolvers := NewResolvers()
	if len(resolvers) == 0 {
		t.Fatal("resolver catalog is empty")
	}
	for i := 0; i < 100; i++ {
		picked := resolvers.PickRandom()
		found := false
		for _, resolver := range resolvers {
			if picked == resolver {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("PickRandom() returned unknown resolver: %+v", picked)
		}
	}
	if got := (Resolvers(nil)).PickRandom(); got != (Resolver{}) {
		t.Fatalf("empty PickRandom() = %+v", got)
	}
}
