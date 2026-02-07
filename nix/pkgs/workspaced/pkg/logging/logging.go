package logging

import (
	"context"
	"encoding/json"
	"log/slog"
	"workspaced/pkg/types"
)

// GetLogger retrieves the logger instance from the context.
// It returns the default slog logger if no logger is found in the context.
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(types.LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// ChannelLogHandler is a custom slog.Handler that broadcasts log records to a channel.
// This is used to stream server-side logs to the client via the daemon connection.
type ChannelLogHandler struct {
	Out    chan<- types.StreamPacket
	Parent slog.Handler
	Ctx    context.Context
}

// Enabled reports whether the handler handles records at the given level.
func (h *ChannelLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

// Handle processes a log record, marshals it to JSON, and sends it as a StreamPacket.
// It also delegates to the parent handler if one is configured.
func (h *ChannelLogHandler) Handle(ctx context.Context, r slog.Record) error {
	entry := types.LogEntry{
		Level:   r.Level.String(),
		Message: r.Message,
		Attrs:   make(map[string]any),
	}
	r.Attrs(func(a slog.Attr) bool {
		entry.Attrs[a.Key] = a.Value.Any()
		return true
	})
	payload, _ := json.Marshal(entry)

	select {
	case h.Out <- types.StreamPacket{Type: "log", Payload: payload}:
	case <-h.Ctx.Done():
		return h.Ctx.Err()
	}

	if h.Parent != nil {
		return h.Parent.Handle(ctx, r)
	}
	return nil
}

// WithAttrs returns a new ChannelLogHandler with the given attributes added.
func (h *ChannelLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ChannelLogHandler{Out: h.Out, Parent: h.Parent.WithAttrs(attrs), Ctx: h.Ctx}
}

// WithGroup returns a new ChannelLogHandler with the given group name.
func (h *ChannelLogHandler) WithGroup(name string) slog.Handler {
	return &ChannelLogHandler{Out: h.Out, Parent: h.Parent.WithGroup(name), Ctx: h.Ctx}
}

// ReportError logs an unexpected error using the logger from the context.
// It serves as the centralized error reporting function.
func ReportError(ctx context.Context, err error, attrs ...slog.Attr) {
	if err == nil {
		return
	}
	args := make([]any, len(attrs)+1)
	args[0] = slog.Any("error", err)
	for i, a := range attrs {
		args[i+1] = a
	}
	GetLogger(ctx).Error("unexpected error", args...)
}
