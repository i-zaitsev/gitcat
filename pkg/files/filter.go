package files

import (
	"os"
	"path/filepath"

	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/ls"
)

// FilterBySize returns a new RepoContent containing only files within the size range.
// minSize and maxSize are in bytes. Use 0 for no minimum, -1 for no maximum.
func FilterBySize(content *ls.RepoContent, minSize, maxSize int64) *ls.RepoContent {
	var filtered ls.RepoContent
	filtered.Root = content.Root

	for _, relPath := range content.Files {
		absPath := filepath.Join(content.Root, relPath)

		info, err := os.Stat(absPath)
		if err != nil {
			log.Warn("failed to stat file for size filtering", "file", relPath, "error", err)
			continue
		}

		size := info.Size()

		if minSize > 0 && size < minSize {
			log.Debug("file filtered by minsize", "file", relPath, "size", size, "minsize", minSize)
			continue
		}

		if maxSize >= 0 && size > maxSize {
			log.Debug("file filtered by maxsize", "file", relPath, "size", size, "maxsize", maxSize)
			continue
		}

		filtered.Files = append(filtered.Files, relPath)
	}

	log.Info("size filtering completed", "input", len(content.Files), "output", len(filtered.Files))
	return &filtered
}
