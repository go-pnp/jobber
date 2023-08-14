package jobber

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

const (
	StatusIdle = iota
	StatusStarting
	StatusRunning
	StatusClosing
)

type Job interface {
	Handle(ctx context.Context) error
	Timer() *time.Timer
	ResetTimer(timer *time.Timer)
}

type Runner struct {
	job     Job
	cancel  context.CancelFunc
	done    chan struct{}
	status  atomic.Int32
	options *options
}

func NewRunner(job Job, optFuncs ...OptionFunc) *Runner {
	opts := &options{
		errorsNotifyTimeout: time.Second,
		errorsChan:          make(chan error),
	}
	for _, optFunc := range optFuncs {
		optFunc(opts)
	}

	return &Runner{
		job:     job,
		done:    make(chan struct{}),
		options: opts,
	}
}

func (r *Runner) Errors() <-chan error {
	return r.options.errorsChan
}

// Start starts the daemon in non-blocking way.
func (r *Runner) Start(ctx context.Context) error {
	if !r.status.CompareAndSwap(StatusIdle, StatusStarting) {
		return errors.New("daemon is not in idle state")
	}
	defer r.status.Store(StatusIdle)

	ctx, r.cancel = context.WithCancel(ctx)
	defer r.cancel() // This is line is important, otherwise the goroutines can leak

	timer := r.job.Timer()
	defer timer.Stop()

	r.status.Store(StatusRunning)

	for {
		select {
		case <-timer.C:
			err := r.job.Handle(ctx)
			if err != nil {
				go r.notifyJobError(err)
			}
			r.job.ResetTimer(timer)
		case <-r.done:
			return nil
		}
	}
}

// Close stops the daemon and waits until goroutines spawned by the daemon are finished.
func (r *Runner) Close() error {
	if r.status.Load() == StatusStarting {
		if err := r.WaitForStatus(StatusRunning); err != nil {
			return err
		}
	}

	if !r.status.CompareAndSwap(StatusRunning, StatusClosing) {
		if r.options.strictMode {
			return errors.New("daemon is not in running state")
		} else {
			return nil
		}
	}

	go r.cancel()
	r.done <- struct{}{}

	return nil
}

func (r *Runner) WaitForStatus(status int32) error {
	for i := 0; i < 100; i++ {
		if r.status.Load() == status {
			return nil
		}

		time.Sleep(time.Millisecond * 50)
	}

	return errors.New("status was not reached")
}

func (r *Runner) notifyJobError(err error) {
	select {
	case r.options.errorsChan <- err:
	case <-time.After(r.options.errorsNotifyTimeout):
	}
}
