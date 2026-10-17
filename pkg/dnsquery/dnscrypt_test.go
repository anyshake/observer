package dnsquery

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/ameshkov/dnscrypt/v2"
	"github.com/jedisct1/go-dnsstamps"
	"github.com/miekg/dns"
)

type cryptHandler struct{}

func (cryptHandler) ServeDNS(writer dnscrypt.ResponseWriter, request *dns.Msg) error {
	return writer.WriteMsg(exampleReply(request))
}

func TestDNSCryptRoundTripAndPinnedKey(t *testing.T) {
	for _, construction := range []dnscrypt.CryptoConstruction{dnscrypt.XSalsa20Poly1305, dnscrypt.XChacha20Poly1305} {
		t.Run(construction.String(), func(t *testing.T) {
			config, err := dnscrypt.GenerateResolverConfig("example.test", nil)
			if err != nil {
				t.Fatal(err)
			}
			config.EsVersion = construction
			certificate, err := config.CreateCert()
			if err != nil {
				t.Fatal(err)
			}
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			server := &dnscrypt.Server{ProviderName: config.ProviderName, ResolverCert: certificate, Handler: cryptHandler{}}
			done := make(chan struct{})
			go func() {
				defer close(done)
				_ = server.ServeTCP(listener)
			}()
			t.Cleanup(func() {
				_ = listener.Close()
				_ = server.Shutdown(context.Background())
				<-done
			})
			stamp, err := config.CreateStamp(listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			backend, err := New(stamp.String())
			if err != nil {
				t.Fatal(err)
			}
			defer backend.Close()
			if err := backend.Open(); err != nil {
				t.Fatal(err)
			}
			queryExample(t, backend)
			// Expired cached certificates must be refreshed before encrypting.
			backend.(*DNSCrypt).resolver.ResolverCert.NotAfter = 0
			queryExample(t, backend)

			pinned, err := dnsstamps.NewServerStampFromString(stamp.String())
			if err != nil {
				t.Fatal(err)
			}
			pinned.ServerPk[0] ^= 0xff
			spoofed := &DNSCrypt{server: pinned.String()}
			if err := spoofed.Open(); err == nil {
				t.Fatal("certificate signed by an unpinned key was accepted")
			}
			if _, err := spoofed.Query(exampleQuestion(), time.Second); err == nil {
				t.Fatal("failed handshake left a usable DNSCrypt client")
			}
		})
	}
}

func TestDNSCryptHandshakeCancellationClosesConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	config, err := dnscrypt.GenerateResolverConfig("example.test", nil)
	if err != nil {
		t.Fatal(err)
	}
	stamp, err := config.CreateStamp(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	started, disconnected := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(disconnected)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		client := &dns.Conn{Conn: conn}
		if _, err := client.ReadMsg(); err != nil {
			return
		}
		close(started)
		_, _ = client.ReadMsg() // No certificate reply; cancellation must close I/O.
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	backend := &DNSCrypt{server: stamp.String()}
	result := make(chan error, 1)
	go func() { result <- backend.OpenContext(ctx, time.Minute) }()
	select {
	case <-started:
	case <-disconnected:
		t.Fatal("certificate query did not reach local server")
	}
	cancel()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatalf("handshake cancellation = %v", err)
	}
	<-disconnected
}
