// Package files reads the files of a repository into documents.
package files

import (
	"bytes"
	"cmp"
	"path"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/repo"
)

// binaryProbe is the number of leading bytes searched for a NUL byte to tell
// a binary file from a text file, the same heuristic that git uses.
const binaryProbe = 8000

// File is a repository file read into memory.
// Path is slash-separated and relative to the repository root.
// Ext is the extension of Path including the dot, or empty when the name has none.
type File struct {
	Path    string
	Ext     string
	Content string
}

// Load reads the files of r concurrently.
//
// If head > 0, only the first head number of lines of each file is kept.
// Every non-empty content ends with a newline.
//
// The files are returned grouped by extension, and sorted by path within a group.
// The files without an extension come first.
// Binary files and files that cannot be read are skipped, and the reason is logged.
func Load(r *repo.Repo, head int) []File {
	paths := slices.Collect(r.Paths())
	loaded := make([]*File, len(paths))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range min(runtime.NumCPU(), len(paths)) {
		wg.Go(func() {
			for i := range jobs {
				loaded[i] = load(r, paths[i], head)
			}
		})
	}
	for i := range paths {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	files := make([]File, 0, len(loaded))
	for _, f := range loaded {
		if f != nil {
			files = append(files, *f)
		}
	}
	slices.SortFunc(files, func(a, b File) int {
		return cmp.Or(cmp.Compare(a.Ext, b.Ext), cmp.Compare(a.Path, b.Path))
	})
	return files
}

// load reads the file p of r and keeps its first head lines.
// It returns nil when the file is binary or cannot be read.
func load(r *repo.Repo, p string, head int) *File {
	data, err := r.ReadFile(p)
	if err != nil {
		log.Warn("failed to read file", "file", p, "error", err)
		return nil
	}
	if isBinary(data) {
		log.Debug("skipping binary file", "file", p)
		return nil
	}
	log.Debug("read file", "file", p, "size", len(data))
	return &File{Path: p, Ext: path.Ext(p), Content: headLines(string(data), head)}
}

// isBinary reports whether data looks like the content of a binary file,
// that is, whether its first [binaryProbe] bytes contain a NUL byte.
func isBinary(data []byte) bool {
	return bytes.IndexByte(data[:min(len(data), binaryProbe)], 0) >= 0
}

// headLines returns the first n lines of s, or all of s when n <= 0.
// A non-empty result always ends with a newline.
func headLines(s string, n int) string {
	if n > 0 {
		end := 0
		for range n {
			i := strings.IndexByte(s[end:], '\n')
			if i < 0 {
				end = len(s)
				break
			}
			end += i + 1
		}
		s = s[:end]
	}
	if s != "" && !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}
