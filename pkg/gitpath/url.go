package gitpath

import (
	"fmt"
	"net/url"
	"strings"
)

// FromURL validates the address of a remote repository and returns its
// location. It accepts http, https, and ssh URLs with a host and a path:
//
//	https://github.com/i-zaitsev/gitcat
//	ssh://git@github.com/i-zaitsev/gitcat.git
//
// and scp-like addresses of the form user@host:path:
//
//	git@github.com:i-zaitsev/gitcat.git
func FromURL(rawURL string) (*GitPath, error) {
	if strings.Contains(rawURL, "://") {
		return fromStdURL(rawURL)
	}
	return fromSCP(rawURL)
}

// fromStdURL parses an address with an explicit scheme.
func fromStdURL(rawURL string) (*GitPath, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid git URL %q: %w", rawURL, err)
	}
	if u.Host == "" || strings.Trim(u.Path, "/") == "" {
		return nil, fmt.Errorf("invalid git URL %q: host and repository path are required", rawURL)
	}
	switch u.Scheme {
	case "http", "https":
		return &GitPath{Path: rawURL, Kind: HTTPS}, nil
	case "ssh":
		return &GitPath{Path: rawURL, Kind: SSH}, nil
	default:
		return nil, fmt.Errorf("unsupported URL scheme %q: must be http, https, or ssh", u.Scheme)
	}
}

// fromSCP parses an scp-like address of the form user@host:path.
func fromSCP(rawURL string) (*GitPath, error) {
	user, rest, hasUser := strings.Cut(rawURL, "@")
	host, path, hasPath := strings.Cut(rest, ":")
	if !hasUser || !hasPath || user == "" || host == "" || path == "" || strings.Contains(host, "/") ||
		strings.HasPrefix(user, "-") || strings.HasPrefix(host, "-") {
		return nil, fmt.Errorf("unsupported repository %q: must be a local directory, "+
			"an http(s) or ssh URL, or an address like git@host:user/repo.git", rawURL)
	}
	return &GitPath{Path: rawURL, Kind: SSH}, nil
}
