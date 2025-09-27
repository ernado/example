package instrumentationdb

import (
	"context"
	"fmt"
	"time"

	"github.com/ernado/example"
	"github.com/ernado/example/internal/semconv"
	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var _ example.DB = (*Instrumentation)(nil)

// Instrumentation layer for [example.DB] interface.
type Instrumentation struct {
	tracer trace.Tracer
	meter  metric.Meter
	db     example.DB

	methodCounter metric.Int64Counter
	methodTime    metric.Float64Histogram
}

// instrument operation call.
func (i Instrumentation) instrument(ctx context.Context, operation string, err *error) context.CancelFunc {
	ctx, span := i.tracer.Start(
		ctx,
		fmt.Sprintf("DB.%s", operation), // TODO: Do we really need "DB." prefix?
		trace.WithAttributes(
			semconv.DBOperation(operation),
		),
	)
	start := time.Now()
	return func() {
		attributes := []attribute.KeyValue{
			semconv.DBOperation(operation),
		}
		if err != nil && *err != nil {
			e := *err
			// TODO: Add to semconv
			span.AddEvent("error", trace.WithAttributes(
				attribute.String("error.message", e.Error()),
				attribute.String("error.type", fmt.Sprintf("%T", e)),
				attribute.String("error.detailed", fmt.Sprintf("%+v", e)),
			))
			span.SetStatus(codes.Error, "Operation failed")
			attributes = append(attributes, semconv.DBResultError())
		} else {
			span.SetStatus(codes.Ok, "Operation succeeded")
			attributes = append(attributes, semconv.DBResultOk())
		}
		span.End()
		i.methodCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		i.methodTime.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			semconv.DBOperation(operation),
		))
	}
}

// TODO: Generate instrumentation layer from interface definition.

func (i Instrumentation) CreateTask(ctx context.Context, title string) (ret *example.Task, err error) {
	defer i.instrument(ctx, "CreateTask", &err)()
	return i.db.CreateTask(ctx, title)
}

func (i Instrumentation) ListTasks(ctx context.Context) (ret []*example.Task, err error) {
	defer i.instrument(ctx, "ListTasks", &err)()
	return i.db.ListTasks(ctx)
}

func (i Instrumentation) DeleteTask(ctx context.Context, id int64) (err error) {
	defer i.instrument(ctx, "DeleteTask", &err)()
	return i.db.DeleteTask(ctx, id)
}

func (i Instrumentation) GenerateError(ctx context.Context) (err error) {
	defer i.instrument(ctx, "GenerateError", &err)()
	return i.db.GenerateError(ctx)
}

func New(
	db example.DB,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*Instrumentation, error) {
	tracer := tracerProvider.Tracer(semconv.SystemDatabase)
	meter := meterProvider.Meter(semconv.SystemDatabase)

	methodCounter, err := meter.Int64Counter(
		semconv.DatabaseOperationCount,
		metric.WithDescription("Counts number of calls to each DB method"),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "create %s counter", semconv.DatabaseOperationCount)
	}

	methodTime, err := meter.Float64Histogram(
		semconv.DatabaseOperationTime,
		metric.WithDescription("Records the time spent in each DB method"),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "create %s histogram", semconv.DatabaseOperationTime)
	}

	return &Instrumentation{
		tracer: tracer,
		meter:  meter,
		db:     db,

		methodCounter: methodCounter,
		methodTime:    methodTime,
	}, nil
}
