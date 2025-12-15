package log

import (
	"io"
	"log/slog"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/termcolor"
)

// colorWriter wraps an io.Writer and colorizes log lines based on their level.
type colorWriter struct {
	out io.Writer
}

// NewColorWriter creates a writer that colorizes entire log lines.
func NewColorWriter(out io.Writer) io.Writer {
	return &colorWriter{out: out}
}

func (w *colorWriter) Write(p []byte) (n int, err error) {
	line := string(p)
	level := extractLevel(line)
	colored := termcolor.ForLevel(level, strings.TrimSuffix(line, "\n"))
	if len(line) > 0 && line[len(line)-1] == '\n' {
		colored += "\n"
	}
	return w.out.Write([]byte(colored))
}

func extractLevel(line string) slog.Level {
	if strings.Contains(line, "level=DEBUG") {
		return slog.LevelDebug
	} else if strings.Contains(line, "level=INFO") {
		return slog.LevelInfo
	} else if strings.Contains(line, "level=WARN") {
		return slog.LevelWarn
	} else if strings.Contains(line, "level=ERROR") {
		return slog.LevelError
	}
	return slog.LevelInfo
}
