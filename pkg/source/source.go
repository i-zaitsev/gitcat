// Package source reads repositories from the disk into a [repo.Repo].
//
// The files of a git work tree are listed by git, so files ignored by
// .gitignore are left out. Any other directory is scanned. Only regular
// files are kept: symbolic links and submodules are skipped. The contents
// are not read here: the returned Repo reads them from disk on demand.
package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/i-zaitsev/gitcat/pkg/git"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/repo"
)

// Local returns the repository in the directory dir.
func Local(ctx context.Context, dir string) (*repo.Repo, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve %s: %w", dir, err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", dir)
	}
	if git.IsWorkTree(ctx, root) {
		return gitWorkTree(ctx, root)
	}
	log.Debug("scanning directory", "dir", root)
	r, err := repo.Scan(os.DirFS(root))
	if err != nil {
		return nil, fmt.Errorf("scan %s: %w", root, err)
	}
	return r, nil
}

// Remote clones the repository at url into the directory dir and returns
// it. The clone must stay on disk while the returned Repo is in use.
func Remote(ctx context.Context, url, dir string) (*repo.Repo, error) {
	if err := git.Clone(ctx, url, dir); err != nil {
		return nil, err
	}
	return Local(ctx, dir)
}

// gitWorkTree returns the repository of the git work tree at root, with
// the regular files that git lists.
func gitWorkTree(ctx context.Context, root string) (*repo.Repo, error) {
	log.Debug("listing git work tree", "dir", root)
	listed, err := git.LsFiles(ctx, root)
	if err != nil {
		return nil, err
	}
	regular := slices.DeleteFunc(listed, func(rel string) bool {
		return !isRegular(root, rel)
	})
	return repo.New(os.DirFS(root), regular), nil
}

// isRegular reports whether the file rel under root is a regular file,
// without following symbolic links.
func isRegular(root, rel string) bool {
	info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil && info.Mode().IsRegular()
}
