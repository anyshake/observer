package dnsquery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/net/http2"
)

type DoH struct {
	server string
	client *http.Client
}

func (d *DoH) Open() error {
	d.client = &http.Client{Transport: &http2.Transport{}}
	return nil
}

func (d *DoH) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if d.client == nil {
		return nil, errors.New("DoH client is not initialized")
	}

	dsnReq, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DoH query message: %w", err)
	}

	// Create a context with timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", d.server, bytes.NewBuffer(dsnReq))
	if err != nil {
		return nil, fmt.Errorf("failed to create DoH query request: %w", err)
	}
	req.Header.Add("Accept", "application/dns-message")
	req.Header.Add("Content-Type", "application/dns-message")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send DoH query: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := dns.Msg{}
	err = response.Unpack(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack DoH response: %w", err)
	}

	return &response, nil
}

func (o *DoH) Close() error {
	o.client = nil
	return nil
}
