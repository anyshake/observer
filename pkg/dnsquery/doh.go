package dnsquery

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

type DoH struct {
	server string
	client *http.Client
}

func (d *DoH) Open() error {
	return d.OpenContext(context.Background(), 5*time.Second)
}

func (d *DoH) OpenContext(ctx context.Context, _ time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d.client == nil {
		d.client = &http.Client{
			Transport: &http.Transport{
				ForceAttemptHTTP2: true,
				TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS12},
				IdleConnTimeout:   90 * time.Second,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if req.URL.Scheme != "https" {
					return errors.New("DoH redirect must use HTTPS")
				}
				if len(via) >= 10 {
					return errors.New("too many DoH redirects")
				}
				return nil
			},
		}
	}
	return nil
}

func (d *DoH) Query(msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	return d.QueryContext(context.Background(), msg, timeout)
}

func (d *DoH) QueryContext(ctx context.Context, msg *dns.Msg, timeout time.Duration) (*dns.Msg, error) {
	if d.client == nil {
		return nil, errors.New("DoH client is not initialized")
	}

	dsnReq, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DoH query message: %w", err)
	}

	ctx, cancel := queryContext(ctx, timeout)
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected DoH status: %s", resp.Status)
	}
	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/dns-message" {
		return nil, errors.New("invalid DoH response content type")
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, dns.MaxMsgSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if len(bodyBytes) > dns.MaxMsgSize {
		return nil, errors.New("DoH response exceeds maximum DNS message size")
	}

	response := dns.Msg{}
	err = response.Unpack(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack DoH response: %w", err)
	}

	return &response, nil
}

func (d *DoH) Close() error {
	if d.client != nil {
		d.client.CloseIdleConnections()
		d.client = nil
	}
	return nil
}
