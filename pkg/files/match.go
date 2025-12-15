package files

import (
	"path/filepath"
	"sort"

	"github.com/i-zaitsev/gitcat/pkg/ls"
)

// DiscoverExt returns a list of all extensions matching the given pattern.
// It only checks the last suffixed part of the filename.
func DiscoverExt(content *ls.RepoContent) []string {
	found := make(map[string]bool)
	for _, filename := range content.Files {
		if ext := filepath.Ext(filename); ext != "" {
			found[ext] = true
		}
	}
	var exts []string
	for ext := range found {
		exts = append(exts, ext)
	}
	sort.Strings(exts)
	return exts
}

// MatchExt checks if a file extension matches the given pattern.
// It returns a new RepoContent with matching files only.
func MatchExt(content *ls.RepoContent, ext ...string) *ls.RepoContent {
	var matched ls.RepoContent
	matched.Root = content.Root
	for _, filename := range content.Files {
		for _, e := range ext {
			if filepath.Ext(filename) == e {
				matched.Files = append(matched.Files, filename)
			}
		}
	}
	return &matched
}
