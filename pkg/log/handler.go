package log

import (
	"log/slog"

	"github.com/i-zaitsev/gitcat/pkg/termcolor"
)

// ColorizeLevel returns a ReplaceAttr function that colorizes log levels.
func ColorizeLevel() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level)
			colorized := termcolor.ForLevel(level, a.Value.String())
			a.Value = slog.StringValue(colorized)
		}
		return a
	}
}
