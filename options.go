package jobber

import "time"

type options struct {
	strictMode          bool
	errorsNotifyTimeout time.Duration
	errorsChan          chan error
}

type OptionFunc func(o *options)

func WithStrictMode() OptionFunc {
	return func(o *options) {
		o.strictMode = true
	}
}

func WithErrorsNotifyTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) {
		o.errorsNotifyTimeout = timeout
	}
}

func WithErrorsCh(errorsCh chan error) OptionFunc {
	return func(o *options) {
		o.errorsChan = errorsCh
	}
}
