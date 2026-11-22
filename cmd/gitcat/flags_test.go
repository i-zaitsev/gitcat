package main

import (
	"testing"

	"github.com/i-zaitsev/gitcat/pkg/preset"
	"github.com/i-zaitsev/pkg/assert"
)

func TestPatterns_Set(t *testing.T) {
	tests := []struct {
		value   string
		want    patterns
		wantStr string
	}{
		{
			"go",
			patterns{"go"},
			"go",
		},
		{
			".go,go.mod",
			patterns{".go", "go.mod"},
			".go,go.mod",
		},
		{
			" *.py[cod] , __pycache__ ,",
			patterns{"*.py[cod]", "__pycache__"},
			"*.py[cod],__pycache__",
		},
		{
			"",
			nil,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			var ps patterns

			assert.NotErr(t, ps.Set(tt.value))
			assert.Eq(t, ps, tt.want)
			assert.Eq(t, ps.String(), tt.wantStr)
		})
	}
}

func TestPatterns_Set_Invalid(t *testing.T) {
	var ps patterns

	assert.ErrContains(t, ps.Set("*.go,["), `invalid pattern "["`)
}

func TestPresets_Set(t *testing.T) {
	var ps presets

	assert.NotErr(t, ps.Set("c|go"))
	assert.NotErr(t, ps.Set("python"))

	assert.Eq(t, ps, presets{preset.C, preset.Go, preset.Python})
	assert.Eq(t, ps.String(), "c|go|python")
	assert.ErrContains(t, ps.Set("cobol"), `unknown preset "cobol"`)
}

func TestPaths_Set(t *testing.T) {
	tests := []struct {
		value string
		want  paths
	}{
		{
			"pkg",
			paths{"pkg"},
		},
		{
			"pkg/, ./cmd/gitcat",
			paths{"pkg", "cmd/gitcat"},
		},
		{
			"pkg//files/",
			paths{"pkg/files"},
		},
		{
			", ., ./",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			var p paths

			assert.NotErr(t, p.Set(tt.value))
			assert.Eq(t, p, tt.want)
		})
	}
}

func TestListFlags_RepeatedFlagsAccumulate(t *testing.T) {
	var (
		ps patterns
		p  paths
	)

	assert.NotErr(t, ps.Set(".go"))
	assert.NotErr(t, ps.Set(".md,go.mod"))
	assert.NotErr(t, p.Set("pkg"))
	assert.NotErr(t, p.Set("cmd/"))

	assert.Eq(t, ps, patterns{".go", ".md", "go.mod"})
	assert.Eq(t, p, paths{"pkg", "cmd"})
}

func TestSize_Set(t *testing.T) {
	var s size

	assert.NotErr(t, s.Set("2"))
	assert.Eq(t, s, size(2048))
	assert.Eq(t, s.String(), "2")
	assert.ErrContains(t, s.Set("1.5"), `invalid size "1.5"`)
	assert.ErrContains(t, s.Set("-1"), `invalid size "-1"`)
	assert.ErrContains(t, s.Set("9223372036854775"), `invalid size "9223372036854775"`)
}

func TestFlagValues_NilString(t *testing.T) {
	var (
		ps *patterns
		p  *paths
		pr *presets
		s  *size
	)

	assert.Eq(t, ps.String(), "")
	assert.Eq(t, p.String(), "")
	assert.Eq(t, pr.String(), "")
	assert.Eq(t, s.String(), "0")
}
