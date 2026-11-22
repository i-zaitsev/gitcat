package main

import (
	"fmt"
	"math"
	"path"
	"strconv"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/filter"
	"github.com/i-zaitsev/gitcat/pkg/preset"
)

// patterns is a comma-separated list of file patterns given as a flag,
// in the syntax described by the [filter] package.
type patterns []string

// String returns the patterns separated by commas.
func (ps *patterns) String() string {
	if ps == nil {
		return ""
	}
	return strings.Join(*ps, ",")
}

// Set appends the comma-separated list in value to the patterns, skipping
// empty items. A malformed pattern is an error. The flag can be repeated.
func (ps *patterns) Set(value string) error {
	for item := range strings.SplitSeq(value, ",") {
		pattern := strings.TrimSpace(item)
		if pattern == "" {
			continue
		}
		if err := filter.CheckPattern(pattern); err != nil {
			return err
		}
		*ps = append(*ps, pattern)
	}
	return nil
}

// presets is a list of presets given as a flag, such as "c|cpp|go".
type presets []preset.Option

// String returns the names of the presets joined with |.
func (ps *presets) String() string {
	if ps == nil {
		return ""
	}
	var b strings.Builder
	for i, p := range *ps {
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(p.String())
	}
	return b.String()
}

// Set appends the presets named in value, see [preset.Parse]. The flag can
// be repeated.
func (ps *presets) Set(value string) error {
	opts, err := preset.Parse(value)
	if err != nil {
		return err
	}
	*ps = append(*ps, opts...)
	return nil
}

// paths is a comma-separated list of slash-separated repository paths given
// as a flag.
type paths []string

// String returns the paths separated by commas.
func (p *paths) String() string {
	if p == nil {
		return ""
	}
	return strings.Join(*p, ",")
}

// Set appends the comma-separated list in value to the paths. Each path is
// cleaned, so "./pkg/" becomes "pkg". Empty items and the repository root
// "." are skipped. The flag can be repeated.
func (p *paths) Set(value string) error {
	for item := range strings.SplitSeq(value, ",") {
		if item = path.Clean(strings.TrimSpace(item)); item != "." {
			*p = append(*p, item)
		}
	}
	return nil
}

// size is a file size in bytes, given as a flag in kilobytes.
type size int64

// String returns the size in kilobytes.
func (s *size) String() string {
	if s == nil {
		return "0"
	}
	return strconv.FormatInt(int64(*s)/1024, 10)
}

// Set parses value as a non-negative whole number of kilobytes that fits
// into [size] when converted to bytes.
func (s *size) Set(value string) error {
	kb, err := strconv.ParseInt(value, 10, 64)
	if err != nil || kb < 0 || kb > math.MaxInt64/1024 {
		return fmt.Errorf("invalid size %q: must be a non-negative integer (KB)", value)
	}
	*s = size(kb * 1024)
	return nil
}
