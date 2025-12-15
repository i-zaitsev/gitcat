package termcolor

import (
	"log/slog"
	"os"

	"golang.org/x/term"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

var fmt func(code, s string) string

func init() {
	if term.IsTerminal(int(os.Stderr.Fd())) {
		fmt = func(code, s string) string {
			return code + s + reset
		}
	} else {
		fmt = func(code, s string) string {
			return s
		}
	}
}

// Red returns text in red color.
func Red(s string) string {
	return fmt(red, s)
}

// Yellow returns text in yellow color.
func Yellow(s string) string {
	return fmt(yellow, s)
}

// Blue returns text in blue color.
func Blue(s string) string {
	return fmt(blue, s)
}

// ForLevel colorizes text based on slog level.
func ForLevel(level slog.Level, s string) string {
	switch {
	case level >= slog.LevelError:
		return Red(s)
	case level >= slog.LevelWarn:
		return Yellow(s)
	case level >= slog.LevelInfo:
		return s
	default:
		return Blue(s)
	}
}
