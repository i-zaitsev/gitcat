package ls

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/gitclone"
	"github.com/i-zaitsev/gitcat/pkg/giturl"
)

type RepoContent struct {
	Root    string
	Content []string
}

func RemoteRepo(repoUrl *giturl.GitUrl, cloneDir string) (*RepoContent, error) {
	var err error

	if err = gitclone.Clone(repoUrl, cloneDir); err != nil {
		return nil, err
	}

	status, err := gitclone.Status(cloneDir)
	if err != nil && strings.Contains(status, "fatal") {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var paths []string
	err = filepath.WalkDir(cloneDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.Contains(path, ".git") {
			return nil
		}

		relPath := strings.TrimPrefix(path, cloneDir)
		relPath = strings.TrimPrefix(relPath, "/")

		paths = append(paths, relPath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &RepoContent{
		Root:    cloneDir,
		Content: paths,
	}, nil
}
