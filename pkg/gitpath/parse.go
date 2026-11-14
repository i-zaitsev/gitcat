package gitpath

import "github.com/i-zaitsev/gitcat/pkg/internal/utils"

// Parse returns the repository location described by location.
// An existing directory is a local repository, see [FromDir].
// Anything else must be a remote address accepted by [FromURL].
func Parse(location string) (*GitPath, error) {
	if utils.DirExists(location) {
		return FromDir(location)
	}
	return FromURL(location)
}
