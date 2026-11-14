// Package testutil builds file trees and git repositories for tests.
package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// WriteFiles creates files in dir, keyed by slash-separated relative path,
// together with any missing parent directories. It stops the test on failure.
func WriteFiles(t testing.TB, dir string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// RequireGit skips the test when git is not installed. Otherwise, it isolates
// git from the user and system configuration of the machine.
func RequireGit(t testing.TB) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
}

// Git runs git with args in dir and returns its trimmed standard output.
// It stops the test on failure.
func Git(t testing.TB, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

// GitRepo creates a git repository in a new temporary directory, writes
// files to its work tree without committing them, and returns its path.
// The test must call [RequireGit] first.
func GitRepo(t testing.TB, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	Git(t, dir, "init", "-q")
	WriteFiles(t, dir, files)
	return dir
}

// Commit stages every file in the repository at dir and commits it.
func Commit(t testing.TB, dir string) {
	t.Helper()
	Git(t, dir, "add", "-A")
	Git(t, dir, "-c", "user.name=test", "-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false", "commit", "-q", "-m", "commit")
}
