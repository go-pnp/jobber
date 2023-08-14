package jobber

import (
	"context"
	"time"
)

type jobFunc func(ctx context.Context) error

type InfinityLoopJob jobFunc

func (i InfinityLoopJob) Handle(ctx context.Context) error {
	return i(ctx)
}

func (i InfinityLoopJob) Timer() *time.Timer {
	return time.NewTimer(0)
}

func (i InfinityLoopJob) ResetTimer(handleErr error, timer *time.Timer) {
	timer.Reset(0)
}

type IntervalJob struct {
	StartImmediately bool
	Interval         time.Duration
	OnErrorInterval  time.Duration
	Job              func(ctx context.Context) error
}

func NewIntervalJob(
	startImmediately bool,
	interval time.Duration,
	onErrorInterval time.Duration,
	job jobFunc,
) IntervalJob {
	return IntervalJob{
		StartImmediately: startImmediately,
		Interval:         interval,
		OnErrorInterval:  onErrorInterval,
		Job:              job,
	}
}

func (i IntervalJob) Handle(ctx context.Context) error {
	return i.Job(ctx)
}

func (i IntervalJob) Timer() *time.Timer {
	if i.StartImmediately {
		return time.NewTimer(0)
	}

	return time.NewTimer(i.Interval)
}

func (i IntervalJob) ResetTimer(handleErr error, timer *time.Timer) {
	if handleErr != nil {
		timer.Reset(i.OnErrorInterval)
		return
	}

	timer.Reset(i.Interval)
}
