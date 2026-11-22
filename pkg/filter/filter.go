// Package filter selects the files of a [repo.Repo] by path, extension, and
// size.
package filter

import (
	"path"
	"slices"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/repo"
)

// Filter is a set of rules that select the files of a repository. A new
// Filter keeps every file, and its builder methods add rules:
//
//	filter.New().IgnoreDotFiles().WithPaths("cmd", "pkg").KeepExt(".go").Apply(r)
//
// A file is kept when it passes every rule.
type Filter struct {
	dotIgnore    bool
	includePaths []string
	excludePaths []string
	exts         []string
	minSize      int64
	maxSize      int64
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

// WithPaths keeps only the files under the given slash-separated paths,
// relative to the repository root. With no paths, every file is kept.
func (f *Filter) WithPaths(paths ...string) *Filter {
	f.includePaths = cleanPaths(paths)
	return f
}

// ExcludePaths leaves out the files under the given slash-separated paths,
// relative to the repository root. Exclusion wins over [Filter.WithPaths].
func (f *Filter) ExcludePaths(paths ...string) *Filter {
	f.excludePaths = cleanPaths(paths)
	return f
}

// KeepExt keeps only the files whose extension, including the dot, is one
// of exts. With no extensions, every file is kept.
func (f *Filter) KeepExt(exts ...string) *Filter {
	f.exts = exts
	return f
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
		return f.keepPath(p) && f.keepExt(p) && f.keepSize(r, p)
	})
}

// keepPath reports whether p passes the dotfile and path rules.
func (f *Filter) keepPath(p string) bool {
	if f.dotIgnore && hasDotPart(p) {
		return false
	}
	under := func(dir string) bool { return isUnder(p, dir) }
	if slices.ContainsFunc(f.excludePaths, under) {
		return false
	}
	return len(f.includePaths) == 0 || slices.ContainsFunc(f.includePaths, under)
}

// keepExt reports whether p passes the extension rule.
func (f *Filter) keepExt(p string) bool {
	return len(f.exts) == 0 || slices.Contains(f.exts, path.Ext(p))
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
