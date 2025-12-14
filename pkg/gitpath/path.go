package gitpath

import (
	"fmt"

	"github.com/i-zaitsev/gitcat/pkg/internal/utils"
)

const (
	SSH   = "ssh"
	HTTPS = "https"
	Local = "local"
)

// GitPath represents a git repository URL with protocol information.
type GitPath struct {
	Path string
	Kind string
}

func FromDir(localDir string) (*GitPath, error) {
	if !utils.DirExists(localDir) {
		return nil, fmt.Errorf("directory %s does not exist", localDir)
	}
	return &GitPath{Path: localDir, Kind: Local}, nil
}

func (g *GitPath) IsLocal() bool {
	return g.Kind == Local
}
