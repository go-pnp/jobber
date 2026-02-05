package jobberotel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type options struct {
	tracerProvider trace.TracerProvider
}

func defaultOptions() options {
	return options{
		tracerProvider: otel.GetTracerProvider(),
	}
}

type Option func(o *options)

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(o *options) {
		o.tracerProvider = tp
	}
}
