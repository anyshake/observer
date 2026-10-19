package jobtracker

import (
	"context"
	"sync"
	"time"
)

type Job struct {
	ID         string
	Kind       string
	Status     JobStatus
	_startedAt time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Error      error
}

type Tracker struct {
	mu   sync.Mutex
	ctx  context.Context
	done chan struct{}
	job  Job
}

func New(ctx context.Context, kind string) *Tracker {
	return &Tracker{
		ctx: ctx,
		job: Job{
			ID:     "",
			Kind:   kind,
			Status: JobStatusIdle,
		},
	}
}
