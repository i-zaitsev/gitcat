// Package filter selects the files of a [repo.Repo] by path, name pattern, and size.
// The patterns use the syntax of [path.Match] and follow the rules of .gitignore
// where they can.
package filter

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/repo"
)

// Filter is a set of rules that select the files of a repository. A new
// Filter keeps every file, and its builder methods add rules:
//
//	filter.New().IgnoreDotFiles().Only("cmd", "pkg").Keep(".go", "go.mod").Ignore("testdata").Apply(r)
//
// A file is kept when it passes every rule.
type Filter struct {
	dotIgnore bool
	only      []string
	ignore    []string
	keep      []string
	minSize   int64
	maxSize   int64
}

// New returns a Filter that keeps every file.
func New() *Filter {
	return &Filter{maxSize: -1}
}

// IgnoreDotFiles leaves out the files and directories starting with a dot.
func (f *Filter) IgnoreDotFiles() *Filter {
	f.dotIgnore = true
	return f
}

// Only keeps only the files under the given slash-separated paths,
// relative to the repository root. With no paths, every file is kept.
func (f *Filter) Only(paths ...string) *Filter {
	f.only = cleanPaths(paths)
	return f
}

// Ignore leaves out the files that match any of the patterns. A pattern
// without a slash matches any part of the path and matches at any depth.
// A pattern with a slash matches from the repository root, together with
// everything under it. Ignoring wins over [Filter.Only] and [Filter.Keep].
// Malformed patterns match nothing, see [CheckPattern].
func (f *Filter) Ignore(patterns ...string) *Filter {
	f.ignore = cleanPaths(patterns)
	return f
}

// Keep keeps only the files that match one of the patterns. A pattern
// without a slash matches the file name, and a pattern with a slash matches
// the whole path. An extension like ".go" is shorthand for "*.go". With no
// patterns, every file is kept. Malformed patterns match nothing, see
// [CheckPattern].
func (f *Filter) Keep(patterns ...string) *Filter {
	f.keep = nil
	for _, p := range cleanPaths(patterns) {
		f.keep = append(f.keep, extGlob(p))
	}
	return f
}

// CheckPattern returns an error when the pattern is malformed.
func CheckPattern(pattern string) error {
	if _, err := path.Match(pattern, ""); err != nil {
		return fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	return nil
}

// SizeRange keeps only the files whose size in bytes is at least minSize
// and at most maxSize. Use 0 for no minimum and -1 for no maximum. With a
// size limit set, the files whose size cannot be read are left out.
func (f *Filter) SizeRange(minSize, maxSize int64) *Filter {
	f.minSize, f.maxSize = minSize, maxSize
	return f
}

// Apply returns a Repo with the files of r that pass the rules.
func (f *Filter) Apply(r *repo.Repo) *repo.Repo {
	return r.Select(func(p string) bool {
		return f.keepPath(p) && f.keepName(p) && f.keepSize(r, p)
	})
}

// keepPath reports whether p passes the dotfile, ignore, and only rules.
func (f *Filter) keepPath(p string) bool {
	if f.dotIgnore && hasDotPart(p) {
		return false
	}
	ignored := func(pattern string) bool { return matchIgnore(pattern, p) }
	if slices.ContainsFunc(f.ignore, ignored) {
		return false
	}
	under := func(dir string) bool { return isUnder(p, dir) }
	return len(f.only) == 0 || slices.ContainsFunc(f.only, under)
}

// keepName reports whether p passes the keep rule.
func (f *Filter) keepName(p string) bool {
	kept := func(pattern string) bool { return matchKeep(pattern, p) }
	return len(f.keep) == 0 || slices.ContainsFunc(f.keep, kept)
}

// keepSize reports whether the file p of r passes the size rule.
func (f *Filter) keepSize(r *repo.Repo, p string) bool {
	if f.minSize <= 0 && f.maxSize < 0 {
		return true
	}
	info, err := r.Stat(p)
	if err != nil {
		log.Warn("failed to stat file for size filtering", "file", p, "error", err)
		return false
	}
	size := info.Size()
	return size >= f.minSize && (f.maxSize < 0 || size <= f.maxSize)
}

// cleanPaths returns the paths in their canonical form, so that "pkg/" and
// "./pkg" become "pkg". The repository root, ".", and empty paths are
// dropped, since they would select every file.
func cleanPaths(paths []string) []string {
	var cleaned []string
	for _, p := range paths {
		if p = path.Clean(p); p != "." {
			cleaned = append(cleaned, p)
		}
	}
	return cleaned
}

// extGlob returns "*" + p when p is an extension like ".go": a leading dot
// and no wildcards or slashes. Any other p is returned as is.
func extGlob(p string) string {
	if strings.HasPrefix(p, ".") && !strings.ContainsAny(p, `*?[\/`) {
		return "*" + p
	}
	return p
}

// matchKeep reports whether the keep pattern matches the path p: its file
// name, or the whole path when the pattern has a slash.
func matchKeep(pattern, p string) bool {
	if !strings.Contains(pattern, "/") {
		p = path.Base(p)
	}
	ok, _ := path.Match(pattern, p)
	return ok
}

// matchIgnore reports whether the ignore pattern matches the path p: any of
// its parts, or, when the pattern has a slash, p or any of its parent
// directories.
func matchIgnore(pattern, p string) bool {
	if !strings.Contains(pattern, "/") {
		for part := range strings.SplitSeq(p, "/") {
			if ok, _ := path.Match(pattern, part); ok {
				return true
			}
		}
		return false
	}
	for dir := p; dir != "."; dir = path.Dir(dir) {
		if ok, _ := path.Match(pattern, dir); ok {
			return true
		}
	}
	return false
}

// isUnder reports whether the slash-separated path p is dir or lies under it.
func isUnder(p, dir string) bool {
	return p == dir || strings.HasPrefix(p, dir+"/")
}

// hasDotPart reports whether any part of the slash-separated path p starts with a dot.
func hasDotPart(p string) bool {
	for part := range strings.SplitSeq(p, "/") {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}
