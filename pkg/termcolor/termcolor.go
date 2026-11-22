// Package termcolor colors text with ANSI escape codes.
// Colors are applied only when the standard error is a terminal.
// Otherwise, the text is left as is.
package termcolor

import (
	"log/slog"
	"os"

	"golang.org/x/term"
)

// ANSI escape codes.
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

// colorize wraps s in the escape code and a reset. It is set at start-up
// and leaves s unchanged when standard error is not a terminal.
var colorize func(code, s string) string

// init picks the colorize function for the standard error stream.
func init() {
	if term.IsTerminal(int(os.Stderr.Fd())) {
		colorize = func(code, s string) string {
			return code + s + reset
		}
	} else {
		colorize = func(_, s string) string {
			return s
		}
	}
}

// Red returns s colored red.
func Red(s string) string {
	return colorize(red, s)
}

// Yellow returns s colored yellow.
func Yellow(s string) string {
	return colorize(yellow, s)
}

// Blue returns s colored blue.
func Blue(s string) string {
	return colorize(blue, s)
}

// ForLevel colors s by the log level: red for errors, yellow for warnings,
// blue for debug messages, and no color for info messages.
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
