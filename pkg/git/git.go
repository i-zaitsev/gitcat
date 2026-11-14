// Package git runs the git commands gitcat needs: shallow clones of remote
// repositories and listing the files of a work tree.
//
// Commands run with terminal prompts disabled, so a repository that needs
// credentials fails instead of waiting for input.
package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/log"
)

// Clone makes a shallow clone of the repository at url into dir. Only the
// latest commit is fetched, since gitcat reads file contents, not history.
// For ssh, the key must be configured; for https, only public repositories
// are supported. The error includes the message that git printed.
func Clone(ctx context.Context, url, dir string) error {
	log.Debug("cloning repository", "url", url, "dir", dir)
	if _, err := run(ctx, "", "clone", "--depth", "1", "--quiet", "--", url, dir); err != nil {
		return fmt.Errorf("git clone %s: %w", url, err)
	}
	return nil
}

// IsWorkTree reports whether dir is inside a git work tree.
func IsWorkTree(ctx context.Context, dir string) bool {
	out, err := run(ctx, dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// LsFiles returns the files of the work tree at dir: the tracked files and
// the untracked files that are not excluded by .gitignore or the other
// standard exclude files. The paths are slash-separated, relative to dir,
// and sorted.
func LsFiles(ctx context.Context, dir string) ([]string, error) {
	out, err := run(ctx, dir, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	var files []string
	for f := range strings.SplitSeq(string(out), "\x00") {
		if f != "" {
			files = append(files, f)
		}
	}
	slices.Sort(files)
	return slices.Compact(files), nil
}

// run executes git with args in dir, or in the current directory when dir
// is empty, and returns its standard output. The error of a failed command
// includes its standard error.
func run(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return nil, fmt.Errorf("%w: %s", err, msg)
		}
		return nil, err
	}
	return out, nil
}
