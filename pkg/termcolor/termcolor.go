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

// colorize wraps s in the escape code and a reset.
// The function is explicitly configured via Init and doesn't do any coloring by default.
var colorize = plain

// plain returns s as is.
func plain(_, s string) string {
	return s
}

// wrap returns s in the escape code and a reset.
func wrap(code, s string) string {
	return code + s + reset
}

// ColorMode defines when the text is colored.
type ColorMode int

// Color modes.
const (
	ColorModeNoColor ColorMode = iota
	ColorModeTerminalOnly
	ColorModeAlways
)

// Config configures the coloring.
type Config struct {
	ColorMode ColorMode
}

// DefaultConfig enables coloring for terminal output only.
func DefaultConfig() *Config {
	return &Config{ColorMode: ColorModeTerminalOnly}
}

// Init picks the colorize function for the output strings.
// Note this is intentionally upper-cased to avoid automatic initialization.
func Init(cfg *Config) {
	if enabled(cfg.ColorMode, term.IsTerminal(int(os.Stderr.Fd()))) {
		colorize = wrap
	} else {
		colorize = plain
	}
}

// enabled reports whether the mode colors the output
// that goes to a terminal or not.
func enabled(mode ColorMode, isTerminal bool) bool {
	return mode == ColorModeAlways || (mode == ColorModeTerminalOnly && isTerminal)
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
