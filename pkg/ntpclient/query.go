package ntpclient

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/beevik/ntp"
)

// Query combines fresh measurements from up to five preferred servers.
// The returned server is the lowest-distance member of the agreeing group.
func (c *Client) Query() (time.Duration, string, error) {
	return c.QueryContext(context.Background())
}

// QueryContext also cancels queued probes and in-flight network I/O when ctx ends.
func (c *Client) QueryContext(ctx context.Context) (time.Duration, string, error) {
	return c.measure(ctx, QUERY_ATTEMPT)
}

// QueryAverage combines up to attempts distinct sources, not repeated bursts
// against one server. Discovery and fallback may query additional sources.
func (c *Client) QueryAverage(attempts int) (time.Duration, error) {
	return c.QueryAverageContext(context.Background(), attempts)
}

func (c *Client) QueryAverageContext(ctx context.Context, attempts int) (time.Duration, error) {
	if attempts <= 0 {
		return 0, errors.New("NTP query attempts must be positive")
	}
	offset, _, err := c.measure(ctx, attempts)
	return offset, err
}

func (c *Client) measure(ctx context.Context, target int) (offset time.Duration, server string, err error) {
	started := time.Now()
	c.infof("NTP synchronization requested: target=%d pool=%d", target, len(c.pool))
	defer func() {
		if err == nil {
			return
		}
		// Per-server timeouts are wrapped failures, not cancellation of this operation.
		if err == context.Canceled || err == context.DeadlineExceeded || errors.Is(err, ErrClosed) {
			c.infof("NTP synchronization stopped: %v (elapsed=%s)", err, time.Since(started))
		} else {
			c.warnf("NTP synchronization failed: %v (elapsed=%s)", err, time.Since(started))
		}
	}()
	if len(c.pool) == 0 {
		return 0, "", errors.New("NTP pool is empty")
	}
	if c.ctx.Err() != nil {
		return 0, "", ErrClosed
	}
	ctx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(c.ctx, cancel)
	defer stop()
	defer cancel()

	// Serialize polling decisions, but allow callers waiting their turn to cancel.
	select {
	case c.gate <- struct{}{}:
	default:
		c.infof("NTP synchronization waiting for the active query")
		select {
		case c.gate <- struct{}{}:
		case <-ctx.Done():
			return 0, "", ctx.Err()
		}
	}
	defer func() { <-c.gate }()
	if c.ctx.Err() != nil {
		return 0, "", ErrClosed
	}
	if err := ctx.Err(); err != nil {
		return 0, "", err
	}
	target = min(target, len(c.pool))
	discover := c.lastDiscovery.IsZero() || time.Since(c.lastDiscovery) >= poolRefreshInterval
	if discover {
		c.infof("NTP pool discovery: checking eligible endpoints in pool of %d", len(c.pool))
	} else {
		c.infof("NTP preferred-source sampling: target=%d", target)
	}
	seen := make(map[string]bool, len(c.pool))
	var samples []clockSample
	var lastErr error

	for round := 0; round <= c.retries; round++ {
		if err := ctx.Err(); err != nil {
			return 0, "", err
		}
		servers := c.selectServers(seen, time.Now())
		if len(servers) == 0 {
			c.infof("NTP no eligible unqueried sources: endpoints are cooling down, disabled, or already tried")
			break
		}
		if !discover || round > 0 {
			needed := target - len(samples)
			if needed <= 0 {
				needed = target // Resolve disagreement using other sources.
			}
			servers = servers[:min(len(servers), needed)]
		}
		for _, server := range servers {
			seen[server] = true
		}
		if round > 0 {
			c.infof("NTP fallback: querying other endpoints, usable_samples=%d target=%d", len(samples), target)
		}
		c.infof("NTP probe batch %d: querying %d endpoints", round+1, len(servers))
		var onResult func(ProbeResult)
		if c.logger != nil {
			completed := 0
			onResult = func(result ProbeResult) {
				completed++
				if ctx.Err() != nil {
					return
				}
				if result.Err != nil || result.Resp == nil {
					c.infof("NTP probe progress: batch=%d completed=%d/%d server=%s no usable reply", round+1, completed, len(servers), result.Server)
				} else {
					c.infof("NTP probe progress: batch=%d completed=%d/%d server=%s offset=%s rtt=%s", round+1, completed, len(servers), result.Server, result.Resp.ClockOffset, result.Resp.RTT)
				}
			}
		}
		results := probeContext(ctx, servers, c.readTimeout, c.timeFunc, c.query, onResult)
		if err := ctx.Err(); err != nil {
			// Cancellation is not a source failure, but must not bypass polling limits.
			for _, server := range servers {
				c.servers[server].nextPoll = time.Now().Add(c.servers[server].interval)
			}
			return 0, "", err
		}
		for _, result := range results {
			sample, err := c.recordResult(result, time.Now())
			if err != nil {
				lastErr = fmt.Errorf("NTP server %s: %w", result.Server, err)
				state := c.servers[result.Server]
				if state.disabled {
					c.warnf("NTP source disabled: server=%s kiss_code=%s", result.Server, result.Resp.KissCode)
				} else if result.Resp != nil && result.Resp.IsKissOfDeath() {
					c.warnf("NTP source backoff: server=%s kiss_code=%s retry_after=%s", result.Server, result.Resp.KissCode, max(state.interval, state.backoff))
				} else {
					c.warnf("NTP sample rejected: server=%s error=%v retry_after=%s", result.Server, err, max(state.interval, state.backoff))
				}
				continue
			}
			samples = append(samples, sample)
		}
		if discover {
			c.lastDiscovery = time.Now()
		}
		c.infof("NTP sample collection: usable=%d target=%d", len(samples), target)
		if len(samples) >= target {
			if _, _, err := combineSamples(samples, target); err == nil {
				break
			} else {
				c.warnf("NTP consensus not reached: %v (samples=%d)", err, len(samples))
			}
		}
	}

	if len(samples) == 0 {
		if lastErr != nil {
			return 0, "", fmt.Errorf("no usable NTP samples: %w", lastErr)
		}
		return 0, "", errors.New("no NTP servers ready to poll; servers are cooling down or disabled")
	}
	c.infof("NTP checking consensus: usable_samples=%d", len(samples))
	offset, server, err = combineSamples(samples, target)
	if ctx.Err() != nil {
		return 0, "", ctx.Err()
	}
	if err == nil {
		agreeing := 0
		for _, sample := range samples {
			if offset < sample.offset-sample.distance || offset > sample.offset+sample.distance {
				c.servers[sample.server].valid = false
				c.warnf("NTP source disagrees: server=%s offset=%s distance=%s; deprioritized", sample.server, sample.offset, sample.distance)
			} else {
				agreeing++
			}
		}
		if len(samples) == 1 {
			c.warnf("NTP single-source synchronization: server=%s; no independent clock cross-check", server)
		}
		c.infof("NTP synchronization complete: server=%s offset=%s agreeing=%d/%d weighted_sources=%d elapsed=%s", server, offset, agreeing, len(samples), min(target, agreeing), time.Since(started))
	}
	return offset, server, err
}

func (c *Client) selectServers(seen map[string]bool, now time.Time) []string {
	var servers []string
	for _, server := range c.pool {
		state := c.servers[server]
		if !seen[server] && !state.disabled && !now.Before(state.nextPoll) {
			servers = append(servers, server)
		}
	}
	slices.SortStableFunc(servers, func(a, b string) int {
		left, right := c.servers[a], c.servers[b]
		if left.valid != right.valid {
			if left.valid {
				return -1
			}
			return 1
		}
		return cmp.Compare(left.distance, right.distance)
	})
	return servers
}

func (c *Client) recordResult(result ProbeResult, now time.Time) (clockSample, error) {
	state := c.servers[result.Server]
	err := result.Err
	if result.Resp != nil && result.Resp.IsKissOfDeath() {
		if result.Resp.KissCode == "DENY" || result.Resp.KissCode == "RSTR" {
			state.disabled = true
		}
		if result.Resp.KissCode == "RATE" {
			state.interval = max(state.interval, min(maxFailureBackoff, state.interval*2), result.Resp.Poll)
		}
		err = ntp.ErrKissOfDeath
	}
	if err == nil {
		if result.Resp == nil {
			err = errors.New("empty NTP response")
		} else {
			err = result.Resp.Validate()
		}
	}
	if err != nil {
		state.valid = false
		state.backoff = min(maxFailureBackoff, max(minPollInterval, state.backoff*2))
		state.nextPoll = now.Add(max(state.interval, state.backoff))
		return clockSample{}, err
	}

	response := result.Resp
	state.distance = max(response.RTT/2+response.RootDelay/2+response.RootDispersion, response.Precision, time.Microsecond)
	state.valid = true
	state.backoff = 0
	state.nextPoll = now.Add(state.interval)
	return clockSample{server: result.Server, offset: response.ClockOffset, distance: state.distance}, nil
}
