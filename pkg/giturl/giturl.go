package giturl

import (
	"fmt"
	"strings"
)

// GitUrl represents a git repository URL with protocol information.
type GitUrl struct {
	Path    string // URL to clone the repo from
	UseHttp bool   // true for SSH, false for HTTPS
}

// Parse validates and parses a git repository URL.
// Supports SSH format (git@host:user/repo.git) and HTTPS format (https://host/user/repo.git).
// Returns nil if the URL format is not recognized.
func Parse(rawURL string) (*GitUrl, error) {
	if strings.HasPrefix(rawURL, "https://") || strings.HasPrefix(rawURL, "http://") {
		if !strings.Contains(rawURL, ".git") && !strings.HasSuffix(rawURL, ".git") {
			return nil, fmt.Errorf("invalid HTTPS git URL: %s", rawURL)
		}
		return &GitUrl{
			Path:    rawURL,
			UseHttp: true,
		}, nil
	}

	if strings.Contains(rawURL, "@") && strings.Contains(rawURL, ":") {
		parts := strings.SplitN(rawURL, "@", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid SSH git URL: %s", rawURL)
		}
		hostPath := parts[1]
		if !strings.Contains(hostPath, ":") {
			return nil, fmt.Errorf("invalid SSH git URL: %s", rawURL)
		}
		return &GitUrl{
			Path:    rawURL,
			UseHttp: false,
		}, nil
	}

	return nil, fmt.Errorf("unsupported URL format: must be SSH (git@host:user/repo.git) or HTTPS (https://host/user/repo.git)")
}
