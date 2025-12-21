package ls

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/gitclone"
	"github.com/i-zaitsev/gitcat/pkg/gitpath"
	"github.com/i-zaitsev/gitcat/pkg/log"
)

type RepoContent struct {
	Root  string
	Files []string
}

type List struct {
	dotIgnore    bool
	includePaths []string
	excludePaths []string
}

func NewList() *List {
	return &List{}
}

func (l *List) IgnoreDotFiles() *List {
	l.dotIgnore = true
	return l
}

// WithPaths configures the list to only include files under the specified paths.
func (l *List) WithPaths(paths ...string) *List {
	l.includePaths = paths
	return l
}

// ExcludePaths configures the list to exclude files under the specified paths.
func (l *List) ExcludePaths(paths ...string) *List {
	l.excludePaths = paths
	return l
}

func (l *List) RemoteRepo(repoUrl *gitpath.GitPath, cloneDir string) (*RepoContent, error) {
	if err := gitclone.Clone(repoUrl, cloneDir); err != nil {
		return nil, err
	}

	if _, err := gitclone.Status(cloneDir); err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	if content, err := l.walkGitRepo(cloneDir); err != nil {
		return nil, err
	} else {
		return content, nil
	}
}

func (l *List) LocalRepo(repoDir string) (*RepoContent, error) {
	state, err := os.Stat(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}
	if !state.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", repoDir)
	}
	return l.walkGitRepo(repoDir)
}

// shouldIncludePath determines if a relative path should be included based on
// include and exclude path filters. For directories, set isDirectory=true to
// also check if the directory is a parent of an include path.
func (l *List) shouldIncludePath(relPath string, isDirectory bool) bool {
	for _, exclude := range l.excludePaths {
		if relPath == exclude || strings.HasPrefix(relPath, exclude+"/") {
			return false
		}
	}

	if len(l.includePaths) == 0 {
		return true
	}

	for _, include := range l.includePaths {
		if relPath == include || strings.HasPrefix(relPath, include+"/") {
			return true
		}
		if isDirectory && strings.HasPrefix(include, relPath+"/") {
			return true
		}
	}

	return false
}

func (l *List) walkGitRepo(repoDir string) (*RepoContent, error) {
	log.Debug("walking repository directory", "dir", repoDir)
	content := RepoContent{
		Root: repoDir,
	}

	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == repoDir {
			return nil
		}

		relPath := strings.TrimPrefix(path, repoDir+"/")
		isDotFile := l.dotIgnore && strings.HasPrefix(filepath.Base(path), ".")

		if d.IsDir() {
			if isDotFile {
				log.Debug("skipping dot directory", "dir", relPath)
				return filepath.SkipDir
			}

			if !l.shouldIncludePath(relPath, true) {
				log.Debug("skipping filtered directory", "dir", relPath)
				return filepath.SkipDir
			}

			return nil
		}

		if isDotFile {
			log.Debug("skipping dot file", "file", relPath)
			return nil
		}

		if !l.shouldIncludePath(relPath, false) {
			log.Debug("skipping filtered file", "file", relPath)
			return nil
		}

		content.Files = append(content.Files, relPath)
		return nil
	}); err != nil {
		return nil, err
	}

	log.Debug("directory walk completed", "files", len(content.Files))
	return &content, nil
}
