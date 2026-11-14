package main

import (
	"fmt"
	"math"
	"path"
	"strconv"
	"strings"
)

// extensions is a comma-separated list of file extensions given as a flag.
// Each extension is stored with a leading dot.
type extensions []string

// String returns the extensions separated by commas.
func (es *extensions) String() string {
	if es == nil {
		return ""
	}
	return strings.Join(*es, ",")
}

// Set appends the comma-separated list in value to the extensions, adding a
// leading dot where it is missing and skipping empty items. The flag can be
// repeated.
func (es *extensions) Set(value string) error {
	for item := range strings.SplitSeq(value, ",") {
		if ext := strings.TrimPrefix(strings.TrimSpace(item), "."); ext != "" {
			*es = append(*es, "."+ext)
		}
	}
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
