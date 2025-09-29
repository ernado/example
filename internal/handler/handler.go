package handler

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/ernado/example"
	"github.com/ernado/example/internal/oas"
	"github.com/ernado/example/internal/semconv"
	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/zctx"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var _ oas.Handler = (*Handler)(nil)

type Handler struct {
	svc    example.Service
	tracer trace.Tracer
	meter  metric.Meter

	tasksCreated  metric.Int64Counter
	tasksDeleted  metric.Int64Counter
	tasksReturned metric.Int64Counter
}

func (h *Handler) GenerateError(ctx context.Context) (*oas.Error, error) {
	ctx, span := h.tracer.Start(ctx, "GenerateError")
	defer span.End()

	err := h.svc.GenerateError(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "generate error")
	}

	return nil, errors.New("unreachable")
}

func (h *Handler) CreateTask(ctx context.Context, req *oas.CreateTaskRequest) (*oas.Task, error) {
	ctx, span := h.tracer.Start(ctx, "CreateTask")
	defer span.End()

	lg := zctx.From(ctx)
	lg.Info("Creating task", semconv.TaskTitle(req.Title))

	task, err := h.svc.CreateTask(ctx, req.Title)
	if err != nil {
		return nil, errors.Wrap(err, "create task")
	}

	h.tasksCreated.Add(ctx, 1)

	return convertToOASTask(task), nil
}

func (h *Handler) DeleteTask(ctx context.Context, params oas.DeleteTaskParams) (oas.DeleteTaskRes, error) {
	lg := zctx.From(ctx).With(semconv.TaskID(params.ID))
	lg.Info("Deleting task")

	err := h.svc.DeleteTask(ctx, params.ID)
	if errors.Is(err, example.ErrTaskNotFound) {
		return nil, &oas.ErrorStatusCode{
			StatusCode: 200,
			Response: oas.Error{
				ErrorMessage: "task not found",
			},
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "delete task")
	}

	lg.Info("Task deleted")

	h.tasksDeleted.Add(ctx, 1)

	return &oas.DeleteTaskNoContent{}, nil
}

func (h *Handler) ListTasks(ctx context.Context) (*oas.TaskList, error) {
	tasks, err := h.svc.ListTasks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "list tasks")
	}

	h.tasksReturned.Add(ctx, int64(len(tasks)))

	return &oas.TaskList{
		Tasks: convertToOASTasks(tasks),
	}, nil
}

func New(
	svc example.Service,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*Handler, error) {
	tracer := tracerProvider.Tracer(semconv.SystemHandler)
	meter := meterProvider.Meter(semconv.SystemHandler)

	// TODO: Generate metrics initialization code from semantic conventions.
	tasksCreated, err := meter.Int64Counter(
		semconv.MetricTasksCreated,
		metric.WithDescription("Number of tasks created"),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "create %s counter", semconv.MetricTasksCreated)
	}
	tasksReturned, err := meter.Int64Counter(
		semconv.MetricTasksReturned,
		metric.WithDescription("Number of tasks returned"),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "create %s counter", semconv.MetricTasksReturned)
	}
	tasksDeleted, err := meter.Int64Counter(
		semconv.MetricTasksDeleted,
		metric.WithDescription("Number of tasks deleted"),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "create %s counter", semconv.MetricTasksDeleted)
	}

	h := &Handler{
		svc:    svc,
		tracer: tracer,
		meter:  meter,

		tasksCreated:  tasksCreated,
		tasksReturned: tasksReturned,
		tasksDeleted:  tasksDeleted,
	}

	return h, nil
}

func (h *Handler) GetHealth(ctx context.Context) (*oas.Health, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errors.New("failed to read build info")
	}

	var commit string
	var buildDate time.Time

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			commit = setting.Value
		case "vcs.time":
			buildDate, _ = time.Parse(time.RFC3339, setting.Value)
		case "vcs.modified":
			if setting.Value == "true" {
				commit += "-modified"
			}
		}
	}

	return &oas.Health{
		Status:    "ok",
		Version:   buildInfo.Main.Version,
		BuildDate: buildDate,
		Commit:    commit,
	}, nil
}

func (h *Handler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	var (
		traceID oas.OptTraceID
		spanID  oas.OptSpanID
	)
	if span := trace.SpanFromContext(ctx).SpanContext(); span.HasTraceID() {
		// Extract trace/span IDs from context if available.
		traceID = oas.NewOptTraceID(oas.TraceID(span.TraceID().String()))
		spanID = oas.NewOptSpanID(oas.SpanID(span.SpanID().String()))
	}
	if v, ok := errors.Into[*oas.ErrorStatusCode](err); ok {
		// Error is already *oas.ErrorStatusCode, just fill fields if needed.
		v.Response.TraceID = traceID
		v.Response.SpanID = spanID
		if v.StatusCode == 0 {
			v.StatusCode = 500
		}
		if v.Response.ErrorMessage == "" {
			v.Response.ErrorMessage = "internal error"
		}
		return v
	}
	return &oas.ErrorStatusCode{
		StatusCode: 500,
		Response: oas.Error{
			ErrorMessage: err.Error(), // NB: Exposing internal error messages may be a security risk.
			TraceID:      traceID,
			SpanID:       spanID,
		},
	}
}
