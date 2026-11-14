package repo

import (
	"io/fs"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/i-zaitsev/pkg/assert"
)

var fsys = fstest.MapFS{
	"README":     {Data: []byte("readme\n")},
	"main.go":    {Data: []byte("package main\n")},
	"pkg/a/a.go": {Data: []byte("package a\n")},
	"link.go":    {Data: []byte("main.go"), Mode: fs.ModeSymlink},
}

func TestNew_SortsAndDeduplicates(t *testing.T) {
	paths := []string{"pkg/a/a.go", "README", "main.go", "README"}

	r := New(fsys, paths)

	assert.Eq(t, slices.Collect(r.Paths()), []string{"README", "main.go", "pkg/a/a.go"})
	assert.Eq(t, r.Len(), 3)
	assert.Eq(t, paths, []string{"pkg/a/a.go", "README", "main.go", "README"})
}

func TestScan_RegularFilesOnly(t *testing.T) {
	r, err := Scan(fsys)

	assert.NotErr(t, err)
	assert.Eq(t, slices.Collect(r.Paths()), []string{"README", "main.go", "pkg/a/a.go"})
}

func TestScan_Empty(t *testing.T) {
	r, err := Scan(fstest.MapFS{})

	assert.NotErr(t, err)
	assert.Eq(t, r.Len(), 0)
}

func TestRepo_Select(t *testing.T) {
	r := New(fsys, []string{"README", "main.go", "pkg/a/a.go"})

	got := r.Select(func(p string) bool { return p != "main.go" })

	assert.Eq(t, slices.Collect(got.Paths()), []string{"README", "pkg/a/a.go"})
	assert.Eq(t, r.Len(), 3)
}

func TestRepo_ReadFile(t *testing.T) {
	r := New(fsys, []string{"main.go", "missing.go"})

	data, err := r.ReadFile("main.go")
	_, errMissing := r.ReadFile("missing.go")

	assert.NotErr(t, err)
	assert.Eq(t, string(data), "package main\n")
	assert.ErrIs(t, errMissing, fs.ErrNotExist)
}

func TestRepo_Stat(t *testing.T) {
	r := New(fsys, []string{"pkg/a/a.go"})

	info, err := r.Stat("pkg/a/a.go")

	assert.NotErr(t, err)
	assert.Eq(t, info.Size(), int64(len("package a\n")))
}
