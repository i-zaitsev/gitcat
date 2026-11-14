package log

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/i-zaitsev/pkg/assert"
)

func TestExtractLevel(t *testing.T) {
	tests := []struct {
		line string
		want slog.Level
	}{
		{
			"time=x level=DEBUG msg=a",
			slog.LevelDebug,
		},
		{
			"time=x level=INFO msg=a",
			slog.LevelInfo,
		},
		{
			"time=x level=WARN msg=a",
			slog.LevelWarn,
		},
		{
			"time=x level=ERROR msg=a",
			slog.LevelError,
		},
		{
			"no level",
			slog.LevelInfo,
		},
		{
			`time=x level=INFO msg="level=ERROR in message"`,
			slog.LevelInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			assert.Eq(t, extractLevel(tt.line), tt.want)
		})
	}
}

func TestColorWriter_Write(t *testing.T) {
	var buf strings.Builder
	w := NewColorWriter(&buf)
	line := []byte("level=WARN msg=careful\n")

	n, err := w.Write(line)

	assert.NotErr(t, err)
	assert.Eq(t, n, len(line))
	assert.Contains(t, buf.String(), "level=WARN msg=careful")
	assert.True(t, strings.HasSuffix(buf.String(), "\n"))
}

func TestSetLogger(t *testing.T) {
	prev := logger
	t.Cleanup(func() { SetLogger(prev) })
	var buf strings.Builder

	SetLogger(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	SetLogger(nil)
	Debug("traced", "key", "value")

	assert.Contains(t, buf.String(), "msg=traced key=value")
}
