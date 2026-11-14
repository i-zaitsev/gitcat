package git

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/i-zaitsev/gitcat/internal/testutil"
	"github.com/i-zaitsev/pkg/assert"
)

func TestLsFiles_RespectsGitignore(t *testing.T) {
	testutil.RequireGit(t)
	dir := testutil.GitRepo(t, map[string]string{
		".gitignore":    "build/\n*.log\n",
		"main.go":       "package main\n",
		"pkg/a.go":      "package pkg\n",
		"build/out.txt": "artifact\n",
		"debug.log":     "log\n",
	})
	testutil.Commit(t, dir)
	testutil.WriteFiles(t, dir, map[string]string{"untracked.go": "package main\n"})

	files, err := LsFiles(context.Background(), dir)

	assert.NotErr(t, err)
	assert.Eq(t, files, []string{".gitignore", "main.go", "pkg/a.go", "untracked.go"})
}

func TestLsFiles_NotWorkTree(t *testing.T) {
	testutil.RequireGit(t)

	_, err := LsFiles(context.Background(), t.TempDir())

	assert.ErrContains(t, err, "git ls-files")
}

func TestIsWorkTree(t *testing.T) {
	testutil.RequireGit(t)
	repo := testutil.GitRepo(t, map[string]string{"sub/a.go": "package sub\n"})
	ctx := context.Background()

	assert.True(t, IsWorkTree(ctx, repo))
	assert.True(t, IsWorkTree(ctx, filepath.Join(repo, "sub")))
	assert.False(t, IsWorkTree(ctx, t.TempDir()))
}

func TestClone_IsShallow(t *testing.T) {
	testutil.RequireGit(t)
	src := testutil.GitRepo(t, map[string]string{"a.go": "package a\n"})
	testutil.Commit(t, src)
	testutil.WriteFiles(t, src, map[string]string{"b.go": "package a\n"})
	testutil.Commit(t, src)
	dst := filepath.Join(t.TempDir(), "clone")

	err := Clone(context.Background(), "file://"+src, dst)

	assert.NotErr(t, err)
	assert.PathExists(t, filepath.Join(dst, "b.go"))
	assert.Eq(t, testutil.Git(t, dst, "rev-list", "--count", "HEAD"), "1")
}

func TestClone_Failure_IncludesGitOutput(t *testing.T) {
	testutil.RequireGit(t)
	missing := filepath.Join(t.TempDir(), "missing")

	err := Clone(context.Background(), "file://"+missing, filepath.Join(t.TempDir(), "clone"))

	assert.ErrContains(t, err, "git clone file://"+missing)
	assert.ErrContains(t, err, "fatal:")
}
