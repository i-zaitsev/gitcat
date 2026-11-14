// Package repo models the content of a repository: a set of files, each
// mapped from a slash-separated path to its content.
//
// A [Repo] does no I/O of its own. It reads its files from an [fs.FS], such
// as the directory of a work tree or an in-memory map, and reads each file
// only when asked to.
package repo

import (
	"io/fs"
	"iter"
	"slices"
)

// Repo is a read-only set of files of a repository. The paths are
// slash-separated, relative to the root of the file system, sorted, and
// unique.
type Repo struct {
	fsys  fs.FS
	paths []string
}

// New returns a Repo with the given paths of fsys. The paths are not
// checked against fsys.
func New(fsys fs.FS, paths []string) *Repo {
	sorted := slices.Clone(paths)
	slices.Sort(sorted)
	return &Repo{fsys: fsys, paths: slices.Compact(sorted)}
}

// Scan returns a Repo with every regular file of fsys. Symbolic links and
// other non-regular files are left out.
func Scan(fsys fs.FS) (*Repo, error) {
	var paths []string
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			paths = append(paths, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return New(fsys, paths), nil
}

// Len returns the number of files.
func (r *Repo) Len() int {
	return len(r.paths)
}

// Paths returns the paths of the files in sorted order.
func (r *Repo) Paths() iter.Seq[string] {
	return slices.Values(r.paths)
}

// Select returns a Repo with the files whose path satisfies keep.
// Both repos read from the same file system.
func (r *Repo) Select(keep func(path string) bool) *Repo {
	var kept []string
	for _, p := range r.paths {
		if keep(p) {
			kept = append(kept, p)
		}
	}
	return &Repo{fsys: r.fsys, paths: kept}
}

// ReadFile returns the content of the file at path.
func (r *Repo) ReadFile(path string) ([]byte, error) {
	return fs.ReadFile(r.fsys, path)
}

// Stat returns the description of the file at path.
func (r *Repo) Stat(path string) (fs.FileInfo, error) {
	return fs.Stat(r.fsys, path)
}
