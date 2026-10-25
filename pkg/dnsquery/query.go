package dnsquery

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// ErrNoAnswer reports a valid negative response or a missing requested record.
// It does not mark the resolver unhealthy or impose failure backoff.
var ErrNoAnswer = errors.New("DNS response has no usable answer for the requested record type")

func (c *Client) Query(msg *dns.Msg) (*dns.Msg, error) {
	return c.QueryContext(context.Background(), msg)
}

// QueryContext returns the first usable answer and joins canceled probes before
// returning. Only resolver quality is cached; every query uses fresh DNS data.
func (c *Client) QueryContext(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	if c.ctx.Err() != nil {
		return nil, ErrClosed
	}
	if msg == nil || len(msg.Question) != 1 {
		return nil, errors.New("DNS query must contain exactly one question")
	}
	ctx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(c.ctx, cancel)
	defer stop()
	defer cancel()
	select {
	case c.gate <- struct{}{}:
	case <-ctx.Done():
		return nil, c.contextError(ctx)
	}
	defer func() { <-c.gate }()
	if err := c.contextError(ctx); err != nil {
		return nil, err
	}

	discover := c.lastDiscovery.IsZero() || time.Since(c.lastDiscovery) >= poolRefreshInterval
	c.infof("DNS query requested: name=%s type=%d pool=%d discovery=%t", msg.Question[0].Name, msg.Question[0].Qtype, len(c.pool), discover)
	seen := make(map[string]bool, len(c.pool))
	var lastErr error
	for round := 0; round <= c.retries; round++ {
		if err := c.contextError(ctx); err != nil {
			return nil, err
		}
		servers := c.selectServers(seen, time.Now(), discover)
		if len(servers) == 0 {
			break
		}
		limit := concurrentQueries
		if !discover && round == 0 {
			limit = 1
		}
		servers = servers[:min(len(servers), limit)]
		for _, server := range servers {
			seen[server] = true
			// Keep discovery moving even if the caller cancels this batch.
			c.servers[server].lastProbe = time.Now()
		}
		c.infof("DNS probe batch %d: endpoints=%d", round+1, len(servers))
		results := c.probe(ctx, servers, msg)
		if err := c.contextError(ctx); err != nil {
			// Caller cancellation and Close are not resolver failures.
			return nil, err
		}
		var winner *probeResult
		for i := range results {
			result := &results[i]
			state := c.servers[result.server]
			if result.canceled {
				continue
			}
			if result.err != nil {
				lastErr = fmt.Errorf("DNS resolver %s: %w", result.server, result.err)
				if errors.Is(result.err, ErrNoAnswer) {
					state.valid, state.rtt = true, result.rtt
					state.backoff, state.nextTry = 0, time.Time{}
					c.infof("DNS resolver returned no answer: server=%s", result.server)
					continue
				}
				state.valid = false
				state.backoff = min(maxFailureBackoff, max(minFailureBackoff, state.backoff*2))
				state.nextTry = time.Now().Add(state.backoff)
				c.warnf("DNS resolver failed: server=%s error=%v retry_after=%s", result.server, result.err, state.backoff)
				continue
			}
			state.valid, state.rtt = true, result.rtt
			state.backoff, state.nextTry = 0, time.Time{}
			if winner == nil {
				winner = result // Results are ordered by completion.
			}
		}
		if winner != nil {
			if discover {
				c.lastDiscovery = time.Now()
			}
			c.infof("DNS query complete: server=%s rtt=%s", winner.server, winner.rtt)
			return winner.response, nil
		}
	}
	if lastErr != nil {
		return nil, fmt.Errorf("all resolvers failed: %w", lastErr)
	}
	return nil, errors.New("all resolvers failed: no DNS resolvers ready; endpoints are cooling down or already tried")
}

func (c *Client) contextError(ctx context.Context) error {
	if c.ctx.Err() != nil {
		return ErrClosed
	}
	return ctx.Err()
}

func (c *Client) selectServers(seen map[string]bool, now time.Time, discover bool) []string {
	var servers []string
	for _, resolver := range c.pool {
		if !seen[resolver.Server] && !now.Before(c.servers[resolver.Server].nextTry) {
			servers = append(servers, resolver.Server)
		}
	}
	slices.SortStableFunc(servers, func(a, b string) int {
		left, right := c.servers[a], c.servers[b]
		if discover {
			return left.lastProbe.Compare(right.lastProbe)
		}
		if left.valid != right.valid {
			if left.valid {
				return -1
			}
			return 1
		}
		return cmp.Compare(left.rtt, right.rtt)
	})
	return servers
}

func (c *Client) probe(ctx context.Context, servers []string, msg *dns.Msg) []probeResult {
	batchCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	resultCh := make(chan probeResult, len(servers))
	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			started := time.Now()
			queryCtx, cancelQuery := queryContext(batchCtx, c.timeout)
			defer cancelQuery()
			backend := c.servers[server].server
			err := queryCtx.Err()
			if err == nil {
				err = backend.OpenContext(queryCtx, c.timeout)
			}
			var response *dns.Msg
			if err == nil {
				request := msg.Copy()
				response, err = backend.QueryContext(queryCtx, request, c.timeout)
				if err == nil {
					err = validateResponse(request, response)
				}
			}
			resultCh <- probeResult{
				server: server, response: response, err: err, rtt: time.Since(started),
				canceled: err != nil && batchCtx.Err() != nil,
			}
		}()
	}
	results := make([]probeResult, 0, len(servers))
	for range servers {
		result := <-resultCh
		results = append(results, result)
		if result.err == nil {
			cancel()
		}
	}
	wg.Wait()
	return results
}

func validateResponse(request, response *dns.Msg) error {
	if response == nil || !response.Response || response.Id != request.Id || response.Opcode != request.Opcode {
		return errors.New("invalid DNS response header")
	}
	if len(response.Question) != 1 {
		return errors.New("invalid DNS response question count")
	}
	question, received := request.Question[0], response.Question[0]
	if !strings.EqualFold(question.Name, received.Name) || question.Qtype != received.Qtype || question.Qclass != received.Qclass {
		return errors.New("DNS response question does not match query")
	}
	if response.Truncated {
		return errors.New("truncated DNS response")
	}
	if response.Rcode == dns.RcodeNameError {
		return ErrNoAnswer
	}
	if response.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("unusable DNS response: truncated=%t rcode=%d", response.Truncated, response.Rcode)
	}
	owners := map[string]bool{strings.ToLower(question.Name): true}
	for range response.Answer {
		before := len(owners)
		for _, answer := range response.Answer {
			if alias, ok := answer.(*dns.CNAME); ok && alias.Hdr.Class == question.Qclass && owners[strings.ToLower(alias.Hdr.Name)] {
				owners[strings.ToLower(alias.Target)] = true
			}
		}
		if len(owners) == before {
			break
		}
	}
	found := false
	for _, answer := range response.Answer {
		if answer.Header().Rrtype != question.Qtype {
			continue
		}
		if answer.Header().Class != question.Qclass || !owners[strings.ToLower(answer.Header().Name)] {
			return errors.New("DNS answer owner or class does not match query")
		}
		switch record := answer.(type) {
		case *dns.A:
			if record.A.To4() == nil {
				continue
			}
		case *dns.AAAA:
			if record.AAAA.To16() == nil {
				continue
			}
		case *dns.TXT:
			if len(record.Txt) == 0 {
				continue
			}
		}
		found = true
	}
	if !found {
		return ErrNoAnswer
	}
	return nil
}
