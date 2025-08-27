package ntpclient

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/beevik/ntp"
	"github.com/samber/lo"
)

func (c *Client) Query() (time.Duration, string, error) {
	probes := c.parallelProbe(c.pool)

	filtered := lo.Filter(probes, func(p probeResult, _ int) bool {
		return p.err == nil && p.resp != nil
	})
	if len(filtered) == 0 {
		return 0, "", errors.New("no available servers in NTP pool")
	}

	type scoredServer struct {
		server       string
		rootDistance time.Duration
		rtt          time.Duration
	}
	var candidates []scoredServer
	for _, p := range filtered {
		rd := p.resp.RootDelay/2 + p.resp.RootDispersion + p.resp.RTT/2
		candidates = append(candidates, scoredServer{p.server, rd, p.resp.RTT})
	}

	bestCandidate := lo.MinBy(candidates, func(a, b scoredServer) bool { return a.rootDistance < b.rootDistance })
	bestServer := bestCandidate.server
	bestRD := bestCandidate.rootDistance
	bestRTT := bestCandidate.rtt

	type attemptResult struct {
		offset time.Duration
		rtt    time.Duration
		err    error
	}
	attempts := make([]attemptResult, QUERY_ATTEMPT)

	sem := make(chan struct{}, CONCURRENT_QUERIES)
	var wg sync.WaitGroup

	for i := 0; i < QUERY_ATTEMPT; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			retryCount := 0
			for {
				resp, err := ntp.QueryWithOptions(bestServer, ntp.QueryOptions{Timeout: c.readTimeout})
				if err != nil || resp == nil {
					retryCount++
					if retryCount > c.retries {
						attempts[idx] = attemptResult{err: fmt.Errorf("all retries failed")}
						return
					}
					continue
				}

				rd := resp.RootDelay/2 + resp.RootDispersion + resp.RTT/2
				if resp.RTT > time.Duration(float64(bestRTT)*1.2) || rd > time.Duration(float64(bestRD)*1.2) {
					retryCount++
					if retryCount > c.retries {
						attempts[idx] = attemptResult{err: fmt.Errorf("RTT/RootDistance too high after retries")}
						return
					}
					continue
				}

				attempts[idx] = attemptResult{offset: resp.ClockOffset, rtt: resp.RTT, err: nil}
				return
			}
		}(i)
	}
	wg.Wait()

	// Keep only valid attempts
	validAttempts := lo.Filter(attempts, func(a attemptResult, _ int) bool { return a.err == nil })
	if len(validAttempts) == 0 {
		return 0, bestServer, errors.New("all query attempts failed")
	}

	// Trim outliers
	slices.SortFunc(validAttempts, func(a, b attemptResult) int { return int(a.offset - b.offset) })
	trim := len(validAttempts) / 10 // remove 10%
	if trim > 0 && len(validAttempts) > 2*trim {
		validAttempts = validAttempts[trim : len(validAttempts)-trim]
	}

	if len(validAttempts) <= 2 {
		median := validAttempts[len(validAttempts)/2].offset
		return median, bestServer, nil
	}

	// Calculate RTT jitter
	var rtts []float64
	for _, a := range validAttempts {
		rtts = append(rtts, float64(a.rtt.Microseconds()))
	}
	mean := 0.0
	for _, r := range rtts {
		mean += r
	}
	mean /= float64(len(rtts))
	var variance float64
	for _, r := range rtts {
		variance += (r - mean) * (r - mean)
	}
	jitter := math.Sqrt(variance / float64(len(rtts))) // stddev (µs)

	// weight average offset by 1/RTT^exponent
	var (
		weightedSum float64
		weightSum   float64
	)
	for _, a := range validAttempts {
		rtt := float64(a.rtt.Microseconds() + 1)
		weight := 1.0 / math.Pow(rtt, lo.Ternary(jitter >= 5000, 8.0, 4.0))
		weightedSum += float64(a.offset) * weight
		weightSum += weight
	}
	if weightSum == 0 {
		return 0, bestServer, errors.New("weightSum=0, all weights dropped")
	}
	finalOffset := time.Duration(weightedSum / weightSum)

	return finalOffset, bestServer, nil
}

func (c *Client) QueryAverage(attempts int) (time.Duration, error) {
	var results []int64

	for i := 0; i < attempts; i++ {
		resp, server, err := c.Query()
		if err != nil {
			return 0, err
		}

		logger.GetLogger("ntp_client").Infof("%d of %d attempts: current server %s, clock offset %d ms", i+1, attempts, server, resp.Milliseconds())
		results = append(results, resp.Milliseconds())

		time.Sleep(time.Second)
	}

	var sum int64
	for _, v := range results {
		sum += v
	}
	avg := float64(sum) / float64(len(results))

	return time.Duration(math.Round(avg)) * time.Millisecond, nil
}
