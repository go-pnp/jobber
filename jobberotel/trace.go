package jobberotel

import (
	"context"

	"github.com/go-pnp/jobber"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/go-pnp/jobber/jobberotel"

type TracedJob struct {
	jobber.Job

	tracer trace.Tracer
}

func NewTracedJob(job jobber.Job, opts ...Option) TracedJob {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return TracedJob{
		Job:    job,
		tracer: options.tracerProvider.Tracer(instrumentationName),
	}
}

func (j TracedJob) Handle(ctx context.Context) error {
	ctx, span := j.tracer.Start(ctx, "job "+j.Job.Name())
	defer span.End()

	span.SetAttributes(attribute.String("job.name", j.Job.Name()))

	err := j.Job.Handle(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
