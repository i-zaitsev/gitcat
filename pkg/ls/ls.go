package ls

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/gitclone"
	"github.com/i-zaitsev/gitcat/pkg/gitpath"
)

type RepoContent struct {
	Root    string
	Content []string
}

func RemoteRepo(repoUrl *gitpath.GitPath, cloneDir string) (*RepoContent, error) {
	if err := gitclone.Clone(repoUrl, cloneDir); err != nil {
		return nil, err
	}

	if _, err := gitclone.Status(cloneDir); err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	if content, err := walkGitRepo(cloneDir); err != nil {
		return nil, err
	} else {
		return content, nil
	}
}

func LocalRepo(repoDir string) (*RepoContent, error) {
	state, err := os.Stat(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}
	if !state.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", repoDir)
	}
	return walkGitRepo(repoDir)
}

func walkGitRepo(repoDir string) (*RepoContent, error) {
	content := RepoContent{
		Root: repoDir,
	}

	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.Contains(path, ".git") {
			return nil
		}

		relPath := strings.TrimPrefix(path, repoDir)
		relPath = strings.TrimPrefix(relPath, "/")

		content.Content = append(content.Content, relPath)
		return nil
	}); err != nil {
		return nil, err
	}
	return &content, nil
}
