package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Options configures the process-wide logger.
type Options struct {
	Level     string // debug|info|warn|error
	Format    string // text|json
	Component string
}

// New builds a slog logger with level/format defaults suitable for daemons.
func New(opts Options) *slog.Logger {
	level := parseLevel(opts.Level)
	format := strings.ToLower(strings.TrimSpace(opts.Format))
	if format == "" {
		format = "text"
	}

	handlerOpts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	logger := slog.New(handler)
	if opts.Component != "" {
		logger = logger.With("component", opts.Component)
	}
	return logger
}

// SetDefault installs logger as the package-level slog default.
func SetDefault(logger *slog.Logger) {
	slog.SetDefault(logger)
}

// Discard returns a logger that drops all output (useful in tests).
func Discard() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 4}))
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
