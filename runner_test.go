package jobber

import (
	"context"
	"testing"
	"time"
)

type TestJob struct{}

func (s TestJob) Handle(ctx context.Context) error {
	return nil
}

func (s TestJob) Timer() *time.Timer {
	return time.NewTimer(time.Minute * 30) // execute immediately
}

func (s TestJob) ResetTimer(handleErr error, timer *time.Timer) {
	timer.Reset(time.Second) // execute every second
}

func TestRunner(t *testing.T) {
	runner := NewRunner(TestJob{})
	go func() {
		if err := runner.Start(context.Background()); err != nil {
			t.Error("starting runner should not return error")
		}
	}()

	if err := runner.WaitForStatus(StatusRunning); err != nil {
		t.Error(err)
	}

	ready := make(chan struct{})
	go func() {
		if err := runner.Start(context.Background()); err == nil {
			t.Error("starting twice should return error")
		}
		close(ready)
	}()
	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Error("looks like second Start didn't return error")
	}

	if err := runner.Close(); err != nil {
		t.Error("closing runner should not return error")
	}

	if err := runner.Close(); err != nil {
		t.Error("closing runner again should not return error")
	}
}
