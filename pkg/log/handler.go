package log

import (
	"io"
	"log/slog"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/termcolor"
)

// colorWriter wraps an [io.Writer] and colors whole log lines written by a
// [slog.TextHandler] according to their level.
type colorWriter struct {
	out io.Writer
}

// NewColorWriter returns a writer that colors each log line by its level
// and writes it to out.
func NewColorWriter(out io.Writer) io.Writer {
	return &colorWriter{out: out}
}

// Write colors the log line in p and writes it to the underlying writer.
func (w *colorWriter) Write(p []byte) (int, error) {
	line := string(p)
	colored := termcolor.ForLevel(extractLevel(line), strings.TrimSuffix(line, "\n"))
	if strings.HasSuffix(line, "\n") {
		colored += "\n"
	}
	if _, err := io.WriteString(w.out, colored); err != nil {
		return 0, err
	}
	return len(p), nil
}

// extractLevel returns the level of a log line in the text format,
// or [slog.LevelInfo] when the line has no level.
func extractLevel(line string) slog.Level {
	_, rest, _ := strings.Cut(line, "level=")
	switch {
	case strings.HasPrefix(rest, "DEBUG"):
		return slog.LevelDebug
	case strings.HasPrefix(rest, "WARN"):
		return slog.LevelWarn
	case strings.HasPrefix(rest, "ERROR"):
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
