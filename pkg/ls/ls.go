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
	dotIgnore bool
}

func NewList() *List {
	return &List{}
}

func (l *List) IgnoreDotFiles() *List {
	l.dotIgnore = true
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

func (l *List) walkGitRepo(repoDir string) (*RepoContent, error) {
	log.Debug("walking repository directory", "dir", repoDir)
	content := RepoContent{
		Root: repoDir,
	}

	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if l.dotIgnore && strings.HasPrefix(path, ".") {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(path, repoDir+"/")
		content.Files = append(content.Files, relPath)
		return nil
	}); err != nil {
		return nil, err
	}

	log.Debug("directory walk completed", "files", len(content.Files))
	return &content, nil
}
