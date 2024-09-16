package scheduler

import (
	"context"
	"time"
)

type Job interface {
	Execute()
}

type Scheduler struct {
	JobQueue Job
	Interval time.Duration
}

func NewJobScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		Interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)

	for {
		select {
		case <-ticker.C:
			go func() {
				s.JobQueue.Execute()
			}()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
