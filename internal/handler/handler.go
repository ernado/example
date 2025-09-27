package handler

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/ernado/example/internal/oas"
	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/trace"
)

var _ oas.Handler = (*Handler)(nil)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
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
