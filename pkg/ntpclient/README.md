# NTP synchronization

`Client` caches source quality and polling deadlines, not clock offsets. Every
successful synchronization uses fresh replies and the configured `TimeFunc`.

## Logging

Logging is optional and uses the package's `Logger` interface:

```go
type Logger interface {
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
}
```

Pass an implementation at construction; existing calls without options remain
valid and silent. `WithLogger(nil)` also leaves logging disabled.

```go
client, err := ntpclient.New(pool, retries, timeoutSeconds, time.Now,
    ntpclient.WithLogger(appLogger))
```

Info messages report waiting, discovery or preferred-source selection, live
per-endpoint completion counts, fallback batches, consensus checking, the final
offset and agreeing source count, and cancellation. Warnings identify rejected
replies and cooldowns, KoD responses, disagreeing sources, single-source
fallback, and unsuccessful synchronization. Offsets use duration formatting
without rounding to milliseconds. These logs describe the measured correction;
the caller remains responsible for applying it to its time source.

Implementations must be safe for concurrent use and return promptly without
calling back into the client. No logging goroutine or global logger is created.
Explorer v1/v2/v3 pass their existing `*logger.Adapter`, which already implements
this interface. The standalone `Probe()` remains silent.

## Lifecycle

The client has no background polling loop. Queries launch bounded probe workers
and join them before returning. `QueryContext` and `QueryAverageContext` support
cancellation while waiting for another query, resolving DNS, or reading UDP.
Canceling a query stops further probes and closes in-flight connections; it does
not close the client or count as a server failure.

`Close()` permanently cancels the client and waits for active probe workers to
exit. It is idempotent; subsequent queries return `ErrClosed`. The original
`Query` and `QueryAverage` APIs remain available and can also be interrupted by
`Close`. Explorer v1/v2/v3 use their session context for all NTP operations and
join their synchronization loop during shutdown, including failed startup.

## Polling

- The first query and the first query at least 30 minutes after discovery probe
  every eligible pool endpoint once, with at most five concurrent requests.
- Routine `Query()` calls target five preferred sources, one request each.
  Discovery replies are themselves measurements; there is no second burst to
  the selected server.
- `QueryAverage(n)` combines up to `n` distinct sources instead of running `n`
  complete synchronization rounds. Discovery can query more than `n` sources.
- `retries` bounds additional batches using other endpoints after failures or
  disagreement. An endpoint is never retried within the same operation.
- Each endpoint has a minimum 64-second interval after completion. Consecutive
  failures back off exponentially to one hour. `RATE` increases the persistent
  interval and honors longer server poll hints; `DENY` and `RSTR` disable the
  endpoint for the lifetime of the client.
- Concurrent calls share these limits. If no source is eligible, the call
  returns an error without sending requests or reusing stale clock offsets.

For a healthy pool of `P >= 5` endpoints, startup `QueryAverage(5)` sends `P`
requests instead of `5 * (P + 5)` before retries. Routine synchronization sends
five instead of `P + 5`; periodic discovery sends `P`. With the current
two-minute caller schedule, a preferred endpoint normally gets one request per
cycle instead of six.

## Accuracy

Replies must pass `ntp.Response.Validate()`. Synchronization distance includes
half the network RTT, half the server root delay, and root dispersion, with a
floor of server precision or one microsecond. Offset intervals must overlap
for a strict majority of usable responses before the best agreeing sources are
weighted by inverse squared distance. Sources outside the agreeing interval
are deprioritized until a later successful measurement.

The weighted estimate stays inside the agreeing interval. Calculations retain
nanosecond resolution even when the supplied monotonic clock uses an epoch
decades away from server time. This is numeric resolution, not a guarantee of
nanosecond network accuracy; path asymmetry and source quality still matter.

Fewer available sources are accepted if they agree. A single usable source is
accepted for availability, but cannot be cross-checked. Disagreement without a
majority returns an error, leaving the caller's current clock unchanged.

With one configured endpoint, each cycle uses one validated reply, without
cross-checking. With two, overlapping offset intervals are required when both
reply; a disagreement cannot be resolved by picking the lower RTT or taking a
blind average. If one fails validation or is unavailable, the other is accepted
as a single-source fallback. Neither case generates extra requests to reach
the default target of five samples. Three or more independent sources are
recommended when rejecting an incorrect clock is important.

Limits apply per client and configured endpoint. Equivalent host spelling,
default ports, and IP spellings are deduplicated, but different DNS aliases can
still reach the same physical server. Separate clients do not share limits.
The standalone `Probe()` is a stateless diagnostic API and is not rate-limited.

## Tests

Tests use injected NTP replies and `testing/synctest`; no external servers are
contacted and polling/backoff tests do not wait in real time. Run:

```sh
go test -race ./pkg/ntpclient ./internal/hardware/explorer
```
