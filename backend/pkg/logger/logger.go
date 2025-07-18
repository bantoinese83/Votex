package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup configures the global logger with the specified level and format
func Setup(level string, isDevelopment bool) {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	if isDevelopment {
		// Development: Human-readable format
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		// Production: JSON format for log aggregation
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// WithContext creates a logger with additional context
func WithContext(ctx map[string]interface{}) *slog.Logger {
	attrs := make([]any, 0, len(ctx)*2)
	for k, v := range ctx {
		attrs = append(attrs, k, v)
	}
	return slog.With(attrs...)
}
