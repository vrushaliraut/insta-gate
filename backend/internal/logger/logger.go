package logger

import (
	"context"
	"io"
	"log/slog"
)

// contextKey is a custom type to prevent context key collisions.
type contextKey string

const (
	RequestIDKey contextKey = "request_id"
)

// contextHandler wraps a slog.Handler to automatically inject context values.
type contextHandler struct {
	slog.Handler
}

// Handle intercepts the log record and injects the request ID if present in the context.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.Handler.Handle(ctx, r)
}

// New initializes a new slog.Logger configured for JSON output.
func New(out io.Writer, env string) *slog.Logger {
	// Default to Debug level, switch to Info for production
	level := slog.LevelDebug
	if env == "production" {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// 1. Create a base standard JSON handler
	baseHandler := slog.NewJSONHandler(out, opts)

	// 2. Wrap it with our custom context handler to intercept logs
	handler := &contextHandler{Handler: baseHandler}

	// 3. Return the fully configured logger
	return slog.New(handler)
}
