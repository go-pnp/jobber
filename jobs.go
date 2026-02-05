package jobber

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	_ Job = (*InfinityJob)(nil)
	_ Job = (*IntervalJob)(nil)
	_ Job = (*CronJob)(nil)
)

type jobFunc = func(ctx context.Context) error

type InfinityJob jobFunc

func (i InfinityJob) Name() string {
	return "infinity"
}

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

type IntervalJobParams struct {
	Name             string
	Job              jobFunc
	Interval         time.Duration
	StartImmediately bool
}

type IntervalJob struct {
	params IntervalJobParams
}

func NewIntervalJob(params IntervalJobParams) IntervalJob {
	return IntervalJob{params: params}
}

func (i IntervalJob) Name() string {
	return i.params.Name
}

func (i IntervalJob) Init(ctx context.Context) error {
	return nil
}

func (i IntervalJob) Handle(ctx context.Context) error {
	return i.params.Job(ctx)
}

func (i IntervalJob) Timer() *time.Timer {
	if i.params.StartImmediately {
		return time.NewTimer(0)
	}

	return time.NewTimer(i.params.Interval)
}

func (i IntervalJob) ResetTimer(timer *time.Timer) {
	timer.Reset(i.params.Interval)
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

func (c CronJob) Name() string {
	return "cron"
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
