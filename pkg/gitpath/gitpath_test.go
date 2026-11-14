package gitpath

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/i-zaitsev/pkg/assert"
)

func TestFromURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		kind    string
		wantErr string
	}{
		{
			"https_with_git_suffix",
			"https://github.com/user/repo.git",
			HTTPS,
			"",
		},
		{
			"https_without_git_suffix",
			"https://github.com/user/repo",
			HTTPS,
			"",
		},
		{
			"http",
			"http://example.com/user/repo",
			HTTPS,
			"",
		},
		{
			"ssh_url",
			"ssh://git@github.com/user/repo.git",
			SSH,
			"",
		},
		{
			"scp_like",
			"git@github.com:user/repo.git",
			SSH,
			"",
		},
		{
			"no_repository_path",
			"https://github.com/",
			"",
			"host and repository path are required",
		},
		{
			"no_host",
			"https:///repo.git",
			"",
			"host and repository path are required",
		},
		{
			"unsupported_scheme",
			"ftp://example.com/repo.git",
			"",
			`unsupported URL scheme "ftp"`,
		},
		{
			"plain_word",
			"repo",
			"",
			"unsupported repository",
		},
		{
			"scp_without_path",
			"git@github.com:",
			"",
			"unsupported repository",
		},
		{
			"scp_without_user",
			"@github.com:repo",
			"",
			"unsupported repository",
		},
		{
			"scp_with_slash_in_host",
			"git@host/x:repo",
			"",
			"unsupported repository",
		},
		{
			"scp_option_as_user",
			"-oProxyCommand=x@host:repo",
			"",
			"unsupported repository",
		},
		{
			"scp_option_as_host",
			"git@-oProxyCommand=x:repo",
			"",
			"unsupported repository",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromURL(tt.url)

			if tt.wantErr != "" {
				assert.ErrContains(t, err, tt.wantErr)
				return
			}
			assert.NotErr(t, err)
			assert.Eq(t, got, &GitPath{Path: tt.url, Kind: tt.kind})
		})
	}
}

func TestParse_LocalDir_IsAbsolute(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}

	sep := string(filepath.Separator)

	got, err := Parse(dir + sep + "sub" + sep + "..")

	assert.NotErr(t, err)
	assert.Eq(t, got, &GitPath{Path: dir, Kind: Local})
	assert.True(t, got.IsLocal())
}

func TestParse_Remote(t *testing.T) {
	got, err := Parse("git@github.com:user/repo.git")

	assert.NotErr(t, err)
	assert.Eq(t, got.Kind, SSH)
	assert.False(t, got.IsLocal())
}

func TestParse_MissingDir_IsNotLocal(t *testing.T) {
	_, err := Parse(filepath.Join(t.TempDir(), "missing"))

	assert.ErrContains(t, err, "unsupported repository")
}

func TestFromDir_Missing(t *testing.T) {
	_, err := FromDir(filepath.Join(t.TempDir(), "missing"))

	assert.ErrContains(t, err, "does not exist")
}

func TestName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{
			"https://github.com/user/repo.git",
			"repo",
		},
		{
			"https://github.com/user/repo/",
			"repo",
		},
		{
			"ssh://git@github.com/user/repo",
			"repo",
		},
		{
			"git@github.com:user/repo.git",
			"repo",
		},
		{
			"git@github.com:repo.git",
			"repo",
		},
		{
			"/home/user/project",
			"project",
		},
		{
			".git",
			"repo",
		},
		{
			"",
			"repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			g := &GitPath{Path: tt.path}

			assert.Eq(t, g.Name(), tt.want)
		})
	}
}
