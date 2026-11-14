// Package log is the package-level structured logger of gitcat.
//
// The functions log through [log/slog], using [slog.Default] if
// it is not overridden. The command-line tool can install its
// logger with [SetLogger].
package log

import "log/slog"

// logger is the logger used by the package-level functions.
var logger = slog.Default()

// SetLogger sets the logger used by the package-level functions.
// A nil logger is ignored.
func SetLogger(l *slog.Logger) {
	if l != nil {
		logger = l
	}
}

// Debug logs a message at the debug level with key-value pairs in args.
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Info logs a message at the info level with key-value pairs in args.
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Warn logs a message at the warn level with key-value pairs in args.
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Error logs a message at the error level with key-value pairs in args.
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
