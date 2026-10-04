package ntpclient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/beevik/ntp"
)

type recordedLog struct {
	level   string
	message string
}

type recordingLogger struct {
	mu      sync.Mutex
	entries []recordedLog
}

func (l *recordingLogger) Infof(format string, args ...any) { l.append("info", format, args...) }
func (l *recordingLogger) Warnf(format string, args ...any) { l.append("warn", format, args...) }

func (l *recordingLogger) append(level, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, recordedLog{level, fmt.Sprintf(format, args...)})
}

func (l *recordingLogger) contains(level, message string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, entry := range l.entries {
		if entry.level == level && strings.Contains(entry.message, message) {
			return true
		}
	}
	return false
}

func TestLoggerIsOptional(t *testing.T) {
	for _, options := range [][]Option{nil, {WithLogger(nil)}} {
		client, _ := newTestClient(t, 1, 0, options...)
		if client.logger != nil {
			t.Fatal("a default logger was installed")
		}
		if _, _, err := client.Query(); err != nil {
			t.Fatalf("query without logger failed: %v", err)
		}
	}
}

func TestLoggerReportsProgressBeforeBatchCompletes(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log := &recordingLogger{}
		client, recorder := newTestClient(t, 2, 0, WithLogger(log))
		if client.logger != log {
			t.Fatal("WithLogger did not install the supplied logger")
		}
		release := make(chan struct{})
		respond := recorder.respond
		recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
			if server == client.pool[1] {
				<-release
			}
			return respond(server, attempt)
		}
		done := make(chan error, 1)
		go func() {
			_, err := client.QueryAverage(5)
			done <- err
		}()
		synctest.Wait()
		for _, message := range []string{
			"NTP synchronization requested:",
			"NTP pool discovery:",
			"NTP probe batch 1: querying 2 endpoints",
			"NTP probe progress: batch=1 completed=1/2 server=" + client.pool[0],
			"offset=1.234567ms rtt=2ms",
		} {
			if !log.contains("info", message) {
				t.Errorf("missing progress message before batch completes: %q", message)
			}
		}
		if log.contains("info", "NTP synchronization complete:") {
			t.Error("synchronization reported complete while a probe was pending")
		}
		close(release)
		if err := <-done; err != nil {
			t.Fatal(err)
		}
		for _, message := range []string{
			"NTP probe progress: batch=1 completed=2/2",
			"NTP sample collection: usable=2 target=2",
			"NTP checking consensus: usable_samples=2",
			"NTP synchronization complete:",
			"agreeing=2/2 weighted_sources=2",
		} {
			if !log.contains("info", message) {
				t.Errorf("missing completion message: %q", message)
			}
		}
	})
}

func TestLoggerReportsFallbackAndBackoff(t *testing.T) {
	log := &recordingLogger{}
	client, recorder := newTestClient(t, 7, 1, WithLogger(log))
	client.lastDiscovery = time.Now()
	respond := recorder.respond
	recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
		if server == client.pool[0] {
			return nil, errors.New("timeout")
		}
		return respond(server, attempt)
	}
	if _, _, err := client.Query(); err != nil {
		t.Fatal(err)
	}
	if !log.contains("warn", "NTP sample rejected: server="+client.pool[0]+" error=timeout retry_after=1m4s") {
		t.Error("missing failed endpoint and cooldown")
	}
	for _, message := range []string{
		"NTP preferred-source sampling: target=5",
		"NTP fallback: querying other endpoints, usable_samples=4 target=5",
		"NTP probe batch 2: querying 1 endpoints",
		"NTP synchronization complete:",
	} {
		if !log.contains("info", message) {
			t.Errorf("missing fallback message: %q", message)
		}
	}
}

func TestLoggerReportsSourceProblems(t *testing.T) {
	for _, tt := range []struct {
		name    string
		count   int
		warning string
		wantErr bool
	}{
		{"single source", 1, "NTP single-source synchronization:", false},
		{"disagreement", 2, "NTP consensus not reached:", true},
		{"outlier", 3, "NTP source disagrees:", false},
		{"empty response", 1, "error=empty NTP response", true},
		{"source timeout", 1, "error=context deadline exceeded", true},
		{"RATE", 1, "NTP source backoff:", true},
		{"DENY", 1, "NTP source disabled:", true},
		{"RSTR", 1, "NTP source disabled:", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			log := &recordingLogger{}
			client, recorder := newTestClient(t, tt.count, 0, WithLogger(log))
			respond := recorder.respond
			recorder.respond = func(server string, attempt int) (*ntp.Response, error) {
				switch tt.name {
				case "RATE", "DENY", "RSTR":
					return &ntp.Response{KissCode: tt.name, Poll: 10 * time.Minute}, nil
				case "empty response":
					return nil, nil
				case "source timeout":
					return nil, context.DeadlineExceeded
				case "disagreement", "outlier":
					if server == client.pool[0] {
						return testResponse(time.Hour, time.Millisecond), nil
					}
				}
				return respond(server, attempt)
			}
			_, _, err := client.Query()
			if (err != nil) != tt.wantErr {
				t.Fatalf("query error = %v, wantErr %v", err, tt.wantErr)
			}
			if !log.contains("warn", tt.warning) {
				t.Errorf("missing source warning: %q", tt.warning)
			}
			if tt.wantErr {
				if !log.contains("warn", "NTP synchronization failed:") {
					t.Error("missing final failure message")
				}
				if log.contains("info", "NTP synchronization complete:") {
					t.Error("failed query reported success")
				}
			}
		})
	}
}

func TestLoggerReportsCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log := &recordingLogger{}
		client, _ := newTestClient(t, 1, 0, WithLogger(log))
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client.query = func(string, ntp.QueryOptions) (*ntp.Response, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		}
		done := make(chan error, 1)
		go func() {
			_, _, err := client.QueryContext(ctx)
			done <- err
		}()
		synctest.Wait()
		cancel()
		if err := <-done; !errors.Is(err, context.Canceled) {
			t.Fatalf("query error = %v", err)
		}
		if !log.contains("info", "NTP synchronization stopped: context canceled") {
			t.Error("missing cancellation message")
		}
		if log.contains("info", "NTP synchronization complete:") || log.contains("warn", "NTP synchronization failed:") {
			t.Error("cancellation was reported as success or source failure")
		}
	})
}
