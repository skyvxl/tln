// Package logger provides structured logging initialization using slog.
package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Format specifies the output format for log lines.
type Format string

// Supported logging formats.
const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// New initializes and returns a configures *slog.Logger.
func New(w io.Writer, format Format, level string) *slog.Logger {
	if w == nil {
		w = os.Stderr
	}

	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var h slog.Handler
	if format == FormatJSON {
		h = slog.NewJSONHandler(w, opts)
	} else {
		h = slog.NewTextHandler(w, opts)
	}

	return slog.New(h)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
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
