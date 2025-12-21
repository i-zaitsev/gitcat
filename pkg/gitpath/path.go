package gitpath

import (
	"fmt"
	"strings"

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

// Extensions represent a list of file extensions to include in the output.
type Extensions []string

func (es *Extensions) String() string {
	if es == nil {
		return ""
	}
	return strings.Join(*es, ",")
}

func (es *Extensions) Set(value string) error {
	ret := strings.Split(value, ",")
	*es = make([]string, len(ret))
	for i, v := range ret {
		(*es)[i] = "." + strings.TrimSpace(strings.TrimPrefix(v, "."))
	}
	return nil
}
