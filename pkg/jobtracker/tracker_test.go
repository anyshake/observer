package jobtracker_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/jobtracker"
)

func TestTrackerPreventsOverlappingJobsAndCopiesSnapshots(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		tracker := jobtracker.New(t.Context(), "purge")
		idle := tracker.Get()
		if idle.Status != jobtracker.JobStatusIdle || idle.Kind != "purge" || idle.ID != "" || idle.StartedAt != nil || idle.FinishedAt != nil {
			t.Fatalf("initial job = %+v", idle)
		}
		now := time.Unix(1700000000, 0)
		release := make(chan struct{})
		unblock := sync.OnceFunc(func() { close(release) })
		defer unblock()
		started := tracker.Start(now, "purge", func(context.Context) error { <-release; return nil })
		if started.Status != jobtracker.JobStatusRunning || started.ID != "1700000000000" || started.StartedAt == nil || !started.StartedAt.Equal(now) {
			t.Fatalf("started job = %+v", started)
		}
		synctest.Wait()
		var duplicateRan atomic.Bool
		duplicate := tracker.Start(now.Add(time.Second), "other", func(context.Context) error { duplicateRan.Store(true); return nil })
		if duplicate.ID != started.ID || duplicate.Kind != "purge" {
			t.Fatalf("overlapping Start() replaced active job: %+v", duplicate)
		}
		*started.StartedAt = time.Time{}
		snapshot := tracker.Get()
		if !snapshot.StartedAt.Equal(now) {
			t.Fatal("mutating Start() result changed tracker")
		}
		*snapshot.StartedAt = time.Time{}
		time.Sleep(3 * time.Second)
		unblock()
		synctest.Wait()
		finished := tracker.Get()
		if duplicateRan.Load() || finished.Status != jobtracker.JobStatusSucceeded || finished.Error != nil || finished.FinishedAt == nil || !finished.FinishedAt.Equal(now.Add(3*time.Second)) {
			t.Fatalf("completed job = %+v, duplicate ran = %v", finished, duplicateRan.Load())
		}
		*finished.FinishedAt = time.Time{}
		if !tracker.Get().FinishedAt.Equal(now.Add(3 * time.Second)) {
			t.Fatal("mutating Get() result changed tracker")
		}
	})
}

func TestTrackerFailureAndRestart(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		tracker := jobtracker.New(t.Context(), "purge")
		wantErr := errors.New("database unavailable")
		failed := tracker.Start(time.Now(), "purge", func(context.Context) error { return wantErr })
		synctest.Wait()
		job := tracker.Get()
		if job.Status != jobtracker.JobStatusFailed || !errors.Is(job.Error, wantErr) || job.FinishedAt == nil {
			t.Fatalf("failed job = %+v", job)
		}
		time.Sleep(time.Second)
		restarted := tracker.Start(time.Now(), "retry", func(context.Context) error { return nil })
		if restarted.ID == failed.ID || restarted.Kind != "retry" || restarted.Error != nil || restarted.FinishedAt != nil || restarted.Status != jobtracker.JobStatusRunning {
			t.Fatalf("restart retained old job state: %+v", restarted)
		}
		synctest.Wait()
		if got := tracker.Get(); got.Status != jobtracker.JobStatusSucceeded || got.Error != nil {
			t.Fatalf("restarted job = %+v", got)
		}
	})
}

func TestTrackerCancelsRunningJob(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		tracker := jobtracker.New(ctx, "purge")
		started := make(chan struct{})
		tracker.Start(time.Now(), "purge", func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			return ctx.Err()
		})
		<-started
		cancel()
		waitCtx, cancelWait := context.WithTimeout(t.Context(), time.Second)
		defer cancelWait()
		if err := tracker.Wait(waitCtx); err != nil {
			t.Fatal(err)
		}
		if job := tracker.Get(); job.Status != jobtracker.JobStatusFailed || !errors.Is(job.Error, context.Canceled) {
			t.Fatalf("canceled job = %+v", job)
		}
	})
}

func TestTrackerSkipsJobAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	tracker := jobtracker.New(ctx, "purge")
	var ran atomic.Bool
	tracker.Start(time.Now(), "purge", func(context.Context) error {
		ran.Store(true)
		return nil
	})
	if err := tracker.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if ran.Load() || !errors.Is(tracker.Get().Error, context.Canceled) {
		t.Fatal("job ran after application cancellation")
	}
}

func TestTrackerWaitDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		tracker := jobtracker.New(t.Context(), "purge")
		if err := tracker.Wait(t.Context()); err != nil {
			t.Fatalf("wait for idle tracker: %v", err)
		}
		release := make(chan struct{})
		defer close(release)
		tracker.Start(time.Now(), "purge", func(context.Context) error {
			<-release
			return nil
		})
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		if err := tracker.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("wait = %v; want deadline exceeded", err)
		}
	})
}
