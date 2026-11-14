package termcolor

import (
	"log/slog"
	"testing"

	"github.com/i-zaitsev/pkg/assert"
)

func TestForLevel_KeepsText(t *testing.T) {
	for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		t.Run(level.String(), func(t *testing.T) {
			assert.Contains(t, ForLevel(level, "message"), "message")
		})
	}
}

func TestForLevel_Info_IsUncolored(t *testing.T) {
	assert.Eq(t, ForLevel(slog.LevelInfo, "message"), "message")
}
