package jobber

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type jobFunc = func(ctx context.Context) error

type InfinityJob jobFunc

func (i InfinityJob) Init(ctx context.Context) error {
	return nil
}

func (i InfinityJob) Handle(ctx context.Context) error {
	return i(ctx)
}

func (i InfinityJob) Timer() *time.Timer {
	return time.NewTimer(0)
}

func (i InfinityJob) ResetTimer(timer *time.Timer) {
	timer.Reset(0)
}

type IntervalJob struct {
	startImmediately bool
	interval         time.Duration
	job              jobFunc
}

func NewIntervalJob(
	startImmediately bool,
	interval time.Duration,
	job jobFunc,
) IntervalJob {
	return IntervalJob{
		startImmediately: startImmediately,
		interval:         interval,
		job:              job,
	}
}

func (i IntervalJob) Init(ctx context.Context) error {
	return nil
}

func (i IntervalJob) Handle(ctx context.Context) error {
	return i.job(ctx)
}

func (i IntervalJob) Timer() *time.Timer {
	if i.startImmediately {
		return time.NewTimer(0)
	}

	return time.NewTimer(i.interval)
}

func (i IntervalJob) ResetTimer(timer *time.Timer) {
	timer.Reset(i.interval)
}

type CronJob struct {
	startImmediately bool
	schedule         cron.Schedule
	job              jobFunc
}

func NewCronJob(
	startImmediately bool,
	cronStr string,
	job jobFunc,
) (CronJob, error) {
	schedule, err := cron.ParseStandard(cronStr)
	if err != nil {
		return CronJob{}, err
	}

	return CronJob{
		startImmediately: startImmediately,
		schedule:         schedule,
		job:              job,
	}, nil
}

func (c CronJob) Init(ctx context.Context) error {
	return nil
}

func (c CronJob) Handle(ctx context.Context) error {
	return c.job(ctx)
}

func (c CronJob) Timer() *time.Timer {
	if c.startImmediately {
		return time.NewTimer(0)
	}

	return time.NewTimer(c.durationToNextRun())
}

func (c CronJob) ResetTimer(timer *time.Timer) {
	timer.Reset(c.durationToNextRun())
}

func (c CronJob) durationToNextRun() time.Duration {
	return c.schedule.Next(time.Now()).Sub(time.Now())
}
