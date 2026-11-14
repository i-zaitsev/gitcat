package main

import (
	"testing"

	"github.com/i-zaitsev/pkg/assert"
)

func TestExtensions_Set(t *testing.T) {
	tests := []struct {
		value   string
		want    extensions
		wantStr string
	}{
		{
			"go",
			extensions{".go"},
			".go",
		},
		{
			".go,md",
			extensions{".go", ".md"},
			".go,.md",
		},
		{
			" .go , md ,",
			extensions{".go", ".md"},
			".go,.md",
		},
		{
			"",
			nil,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			var es extensions

			assert.NotErr(t, es.Set(tt.value))
			assert.Eq(t, es, tt.want)
			assert.Eq(t, es.String(), tt.wantStr)
		})
	}
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
		es extensions
		p  paths
	)

	assert.NotErr(t, es.Set("go"))
	assert.NotErr(t, es.Set("md,txt"))
	assert.NotErr(t, p.Set("pkg"))
	assert.NotErr(t, p.Set("cmd/"))

	assert.Eq(t, es, extensions{".go", ".md", ".txt"})
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
		es *extensions
		p  *paths
		s  *size
	)

	assert.Eq(t, es.String(), "")
	assert.Eq(t, p.String(), "")
	assert.Eq(t, s.String(), "0")
}
