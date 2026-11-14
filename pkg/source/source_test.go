package source

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/i-zaitsev/gitcat/internal/testutil"
	"github.com/i-zaitsev/gitcat/pkg/repo"
	"github.com/i-zaitsev/pkg/assert"
)

// paths returns the paths of the files of r.
func paths(r *repo.Repo) []string {
	return slices.Collect(r.Paths())
}

var tree = map[string]string{
	".env":           "SECRET=1\n",
	".github/ci.yml": "name: ci\n",
	"README":         "readme\n",
	"main.go":        "package main\n",
	"pkg/a/a.go":     "package a\n",
}

func TestLocal_Directory_ScansEveryFile(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteFiles(t, dir, tree)

	r, err := Local(context.Background(), dir)

	assert.NotErr(t, err)
	assert.Eq(t, paths(r), []string{".env", ".github/ci.yml", "README", "main.go", "pkg/a/a.go"})
	data, err := r.ReadFile("pkg/a/a.go")
	assert.NotErr(t, err)
	assert.Eq(t, string(data), "package a\n")
}

func TestLocal_TrailingSlashRoot(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteFiles(t, dir, map[string]string{"pkg/a.go": "package pkg\n"})

	r, err := Local(context.Background(), dir+string(filepath.Separator))

	assert.NotErr(t, err)
	assert.Eq(t, paths(r), []string{"pkg/a.go"})
}

func TestLocal_SkipsSymlinks(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteFiles(t, dir, map[string]string{"a.go": "package a\n"})
	if err := os.Symlink(filepath.Join(dir, "a.go"), filepath.Join(dir, "link.go")); err != nil {
		t.Skip("symlinks are not supported:", err)
	}

	r, err := Local(context.Background(), dir)

	assert.NotErr(t, err)
	assert.Eq(t, paths(r), []string{"a.go"})
}

func TestLocal_NotDirectory(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteFiles(t, dir, map[string]string{"a.go": "package a\n"})

	_, err := Local(context.Background(), filepath.Join(dir, "a.go"))

	assert.ErrContains(t, err, "not a directory")
}

func TestLocal_GitWorkTree_RespectsGitignore(t *testing.T) {
	testutil.RequireGit(t)
	files := map[string]string{
		".gitignore":    "build/\n",
		"build/out.txt": "artifact\n",
	}
	for rel, content := range tree {
		files[rel] = content
	}
	dir := testutil.GitRepo(t, files)

	r, err := Local(context.Background(), dir)

	assert.NotErr(t, err)
	assert.Eq(t, paths(r), []string{".env", ".github/ci.yml", ".gitignore", "README", "main.go", "pkg/a/a.go"})
}

func TestLocal_GitWorkTree_SkipsSymlinks(t *testing.T) {
	testutil.RequireGit(t)
	dir := testutil.GitRepo(t, map[string]string{"a.go": "package a\n"})
	if err := os.Symlink(filepath.Join(dir, "a.go"), filepath.Join(dir, "link.go")); err != nil {
		t.Skip("symlinks are not supported:", err)
	}

	r, err := Local(context.Background(), dir)

	assert.NotErr(t, err)
	assert.Eq(t, paths(r), []string{"a.go"})
}

func TestRemote_ClonesAndLists(t *testing.T) {
	testutil.RequireGit(t)
	src := testutil.GitRepo(t, map[string]string{"main.go": "package main\n", "pkg/a.go": "package pkg\n"})
	testutil.Commit(t, src)
	dst := filepath.Join(t.TempDir(), "clone")

	r, err := Remote(context.Background(), "file://"+src, dst)

	assert.NotErr(t, err)
	assert.Eq(t, paths(r), []string{"main.go", "pkg/a.go"})
	data, err := r.ReadFile("main.go")
	assert.NotErr(t, err)
	assert.Eq(t, string(data), "package main\n")
}
