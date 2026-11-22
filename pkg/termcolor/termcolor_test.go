package termcolor

import (
	"log/slog"
	"testing"

	"github.com/i-zaitsev/pkg/assert"
)

// useConfig applies the config and restores the plain output after the test.
func useConfig(t *testing.T, mode ColorMode) {
	t.Helper()
	t.Cleanup(func() { Init(&Config{ColorMode: ColorModeNoColor}) })
	Init(&Config{ColorMode: mode})
}

func TestEnabled(t *testing.T) {
	tests := []struct {
		name       string
		mode       ColorMode
		isTerminal bool
		want       bool
	}{
		{
			"no_color_terminal",
			ColorModeNoColor,
			true,
			false,
		},
		{
			"no_color_pipe",
			ColorModeNoColor,
			false,
			false,
		},
		{
			"terminal_only_terminal",
			ColorModeTerminalOnly,
			true,
			true,
		},
		{
			"terminal_only_pipe",
			ColorModeTerminalOnly,
			false,
			false,
		},
		{
			"always_terminal",
			ColorModeAlways,
			true,
			true,
		},
		{
			"always_pipe",
			ColorModeAlways,
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Eq(t, enabled(tt.mode, tt.isTerminal), tt.want)
		})
	}
}

func TestInit_Always_ColorsText(t *testing.T) {
	useConfig(t, ColorModeAlways)

	assert.Eq(t, Red("text"), "\033[31mtext\033[0m")
	assert.Eq(t, Yellow("text"), "\033[33mtext\033[0m")
	assert.Eq(t, Blue("text"), "\033[34mtext\033[0m")
}

func TestInit_NoColor_KeepsText(t *testing.T) {
	useConfig(t, ColorModeNoColor)

	assert.Eq(t, Red("text"), "text")
}

func TestInit_NoColor_TurnsColorOff(t *testing.T) {
	useConfig(t, ColorModeAlways)

	Init(&Config{ColorMode: ColorModeNoColor})

	assert.Eq(t, Red("text"), "text")
}

func TestForLevel_Colored(t *testing.T) {
	useConfig(t, ColorModeAlways)
	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{
			"debug",
			slog.LevelDebug,
			"\033[34mmessage\033[0m",
		},
		{
			"info",
			slog.LevelInfo,
			"message",
		},
		{
			"warn",
			slog.LevelWarn,
			"\033[33mmessage\033[0m",
		},
		{
			"error",
			slog.LevelError,
			"\033[31mmessage\033[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Eq(t, ForLevel(tt.level, "message"), tt.want)
		})
	}
}

func TestForLevel_Default_IsPlain(t *testing.T) {
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		t.Run(level.String(), func(t *testing.T) {
			assert.Eq(t, ForLevel(level, "message"), "message")
		})
	}
}
