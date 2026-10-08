package dnsquery

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestUDPQuery(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	go func() {
		buf := make([]byte, 1500)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return
		}
		var request dns.Msg
		if err := request.Unpack(buf[:n]); err != nil {
			return
		}
		packed, _ := exampleReply(&request).Pack()
		_, _ = conn.WriteTo(packed, addr)
	}()

	endpoint := (&url.URL{Scheme: "udp", Host: conn.LocalAddr().String()}).String()
	server, err := New(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if err := server.Open(); err != nil {
		t.Fatal(err)
	}
	response := queryExample(t, server)
	if answer, ok := response.Answer[0].(*dns.A); !ok || !answer.A.Equal(net.ParseIP("1.2.3.4")) {
		t.Fatalf("answer = %#v", response.Answer)
	}
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := server.Query(exampleQuestion(), time.Second); err == nil {
		t.Fatal("closed UDP client accepted a query")
	}
}

func TestDoHQuery(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bad" {
			_, _ = w.Write([]byte("not-dns"))
			return
		}
		body, _ := io.ReadAll(r.Body)
		var request dns.Msg
		if err := request.Unpack(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		packed, _ := exampleReply(&request).Pack()
		_, _ = io.Copy(w, bytes.NewReader(packed))
	}))
	defer httpServer.Close()

	server := &DoH{server: httpServer.URL}
	if _, err := server.Query(exampleQuestion(), time.Second); err == nil {
		t.Fatal("unopened DoH client accepted a query")
	}
	if err := server.Open(); err != nil {
		t.Fatal(err)
	}
	server.client = httpServer.Client()
	queryExample(t, server)

	bad := &DoH{server: httpServer.URL + "/bad"}
	if err := bad.Open(); err != nil {
		t.Fatal(err)
	}
	bad.client = httpServer.Client()
	if _, err := bad.Query(exampleQuestion(), time.Second); err == nil {
		t.Fatal("invalid DoH body accepted")
	}
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDoTAndDNSCryptErrors(t *testing.T) {
	dot, err := New("tls://127.0.0.1:1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dot.Query(exampleQuestion(), time.Second); err == nil {
		t.Fatal("unopened DoT client accepted a query")
	}
	if err := dot.Open(); err != nil {
		t.Fatal(err)
	}
	if _, err := dot.Query(exampleQuestion(), 200*time.Millisecond); err == nil {
		t.Fatal("closed DoT port accepted a query")
	}
	if err := dot.Close(); err != nil {
		t.Fatal(err)
	}

	crypt, err := New("sdns://not-a-stamp")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := crypt.Query(exampleQuestion(), time.Second); err == nil {
		t.Fatal("unopened DNSCrypt client accepted a query")
	}
	if err := crypt.Open(); err == nil {
		t.Fatal("invalid DNSCrypt stamp accepted")
	}
	if err := crypt.Close(); err != nil {
		t.Fatal(err)
	}
}

func exampleQuestion() *dns.Msg {
	message := new(dns.Msg)
	message.SetQuestion("example.com.", dns.TypeA)
	return message
}

func queryExample(t *testing.T, server IServer) *dns.Msg {
	t.Helper()
	ask := server.Query
	response, err := ask(exampleQuestion(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Answer) != 1 {
		t.Fatalf("answers = %#v", response.Answer)
	}
	return response
}

func exampleReply(request *dns.Msg) *dns.Msg {
	response := new(dns.Msg)
	response.SetReply(request)
	response.Answer = []dns.RR{&dns.A{
		Hdr: dns.RR_Header{Name: request.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		A:   net.ParseIP("1.2.3.4"),
	}}
	return response
}
