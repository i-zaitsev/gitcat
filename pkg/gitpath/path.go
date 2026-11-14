// Package gitpath parses the location of a git repository given on the
// command line: a local directory, an http(s) or ssh URL, or an scp-like
// address such as git@github.com:i-zaitsev/gitcat.git.
package gitpath

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/internal/utils"
)

// Kinds of repository locations.
const (
	SSH   = "ssh"
	HTTPS = "https"
	Local = "local"
)

// GitPath is a repository location together with its kind,
// one of [SSH], [HTTPS], or [Local].
type GitPath struct {
	Path string
	Kind string
}

// FromDir returns the location of the local repository in localDir.
// The directory must exist, and its path is made absolute.
func FromDir(localDir string) (*GitPath, error) {
	if !utils.DirExists(localDir) {
		return nil, fmt.Errorf("directory %s does not exist", localDir)
	}
	abs, err := filepath.Abs(localDir)
	if err != nil {
		return nil, fmt.Errorf("resolve %s: %w", localDir, err)
	}
	return &GitPath{Path: abs, Kind: Local}, nil
}

// IsLocal reports whether the location is a local directory.
func (g *GitPath) IsLocal() bool {
	return g.Kind == Local
}

// Name returns the name of the repository: the last element of its path
// without the .git suffix. It returns "repo" when the path has no usable
// name. The name is the default directory to clone a remote repository into.
func (g *GitPath) Name() string {
	name := strings.TrimRight(g.Path, "/")
	if i := strings.LastIndexAny(name, "/:"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimSuffix(name, ".git")
	if name == "" || name == "." || name == ".." {
		return "repo"
	}
	return name
}
