package observ

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{Logger: l.Logger.With("request_id", getRequestID(ctx))}
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{Logger: l.Logger.With(args...)}
}

func getRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value("request_id").(string); ok {
		return reqID
	}
	return ""
}

type LogEntry struct {
	logger *Logger
	fields map[string]any
}

func (l *Logger) Entry() *LogEntry {
	return &LogEntry{
		logger: l,
		fields: make(map[string]any),
	}
}

func (e *LogEntry) WithField(key string, value any) *LogEntry {
	e.fields[key] = value
	return e
}

func (e *LogEntry) WithError(err error) *LogEntry {
	if err != nil {
		e.fields["error"] = err.Error()
	}
	return e
}

func (e *LogEntry) Info(msg string) {
	e.logger.WithFields(e.fields).Info(msg)
}

func (e *LogEntry) Error(msg string) {
	e.logger.WithFields(e.fields).Error(msg)
}

func (e *LogEntry) Warn(msg string) {
	e.logger.WithFields(e.fields).Warn(msg)
}

func (e *LogEntry) Debug(msg string) {
	e.logger.WithFields(e.fields).Debug(msg)
}
