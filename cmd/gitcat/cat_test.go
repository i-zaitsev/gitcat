package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/i-zaitsev/gitcat/internal/testutil"
	"github.com/i-zaitsev/gitcat/pkg/gitpath"
	"github.com/i-zaitsev/gitcat/pkg/output"
	"github.com/i-zaitsev/pkg/assert"
)

// sampleRepo creates a plain directory, not a git repository, with a mix of
// source, extensionless, hidden, binary, and large files.
func sampleRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	testutil.WriteFiles(t, dir, map[string]string{
		"main.go":          "package main\n",
		"LICENSE":          "MIT\n",
		"docs/guide.md":    "# Guide\n",
		".env":             "SECRET=1\n",
		"logo.png":         "\x89PNG\r\n\x1a\n\x00\x00",
		"testdata/big.txt": strings.Repeat("x", 3000),
	})
	return dir
}

func TestCatParse_Defaults(t *testing.T) {
	dir := t.TempDir()
	c := newCat()

	err := c.Parse([]string{dir})

	assert.NotErr(t, err)
	assert.Eq(t, c.location, &gitpath.GitPath{Path: dir, Kind: gitpath.Local})
	assert.Eq(t, c.outFmt, output.FormatJSONL)
	assert.Eq(t, c.maxSize, size(-1))
	assert.Eq(t, c.headLines, 0)
	assert.False(t, c.dotFiles)
	assert.False(t, c.tmpClone)
}

func TestCatParse_Flags(t *testing.T) {
	c := newCat()

	err := c.Parse([]string{
		"-fmt", "md", "-keep", "go,md", "-path", "pkg/,cmd", "-exclude", "testdata",
		"-minsize", "1", "-maxsize", "4", "-head", "5", "-dot", "-tmp", "-out", "ctx",
		"git@github.com:user/repo.git",
	})

	assert.NotErr(t, err)
	assert.Eq(t, c.outFmt, output.FormatMarkdown)
	assert.Eq(t, c.keepExt, extensions{".go", ".md"})
	assert.Eq(t, c.includePaths, paths{"pkg", "cmd"})
	assert.Eq(t, c.excludePaths, paths{"testdata"})
	assert.Eq(t, c.minSize, size(1024))
	assert.Eq(t, c.maxSize, size(4096))
	assert.Eq(t, c.headLines, 5)
	assert.True(t, c.dotFiles)
	assert.True(t, c.tmpClone)
	assert.Eq(t, c.outFile, "ctx")
	assert.Eq(t, c.location, &gitpath.GitPath{Path: "git@github.com:user/repo.git", Kind: gitpath.SSH})
}

func TestCatParse_Errors(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			"missing_repository",
			nil,
			"cat: repository is required",
		},
		{
			"extra_argument",
			[]string{dir, "extra"},
			"cat: invalid argument: extra",
		},
		{
			"unsupported_repository",
			[]string{"not-a-repo"},
			"cat: unsupported repository",
		},
		{
			"invalid_format",
			[]string{"-fmt", "xml", dir},
			`invalid format "xml"`,
		},
		{
			"invalid_size",
			[]string{"-maxsize", "big", dir},
			`invalid size "big"`,
		},
		{
			"negative_head",
			[]string{"-head", "-1", dir},
			"cat: invalid -head: -1",
		},
		{
			"dir_with_tmp",
			[]string{"-tmp", "-dir", "x", dir},
			"cat: -dir and -tmp cannot be used together",
		},
		{
			"unknown_flag",
			[]string{"-nope", dir},
			"flag provided but not defined: -nope",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newCat().Parse(tt.args)

			assert.ErrContains(t, err, tt.wantErr)
		})
	}
}

func TestCatParse_Help(t *testing.T) {
	err := newCat().Parse([]string{"-h"})

	assert.ErrIs(t, err, flag.ErrHelp)
}

func TestCat_Markdown(t *testing.T) {
	dir := sampleRepo(t)

	stdout, _, err := runGitcat(t, "-fmt", "md", dir)

	assert.NotErr(t, err)
	assert.ContainsAll(t, stdout,
		"## LICENSE\n\n```\nMIT\n```",
		"## main.go\n*Extension: .go*\n\n```go\npackage main\n```",
		"## docs/guide.md",
		"## testdata/big.txt",
	)
	assert.ContainsNone(t, stdout, ".env", "SECRET", "logo.png", dir)
}

func TestCat_Filters(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			"keep_extensions",
			[]string{"-keep", "go"},
			"package main\n\n",
		},
		{
			"include_path",
			[]string{"-path", "docs"},
			"# Guide\n\n",
		},
		{
			"exclude_path_and_dotfiles",
			[]string{"-exclude", "docs,testdata", "-dot", "-keep", "env,go"},
			"SECRET=1\n\npackage main\n\n",
		},
		{
			"max_size",
			[]string{"-maxsize", "1", "-keep", "txt,go"},
			"package main\n\n",
		},
		{
			"min_size",
			[]string{"-minsize", "1"},
			strings.Repeat("x", 3000) + "\n\n",
		},
		{
			"head",
			[]string{"-head", "1", "-keep", "go,md"},
			"package main\n\n# Guide\n\n",
		},
	}
	dir := sampleRepo(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"-fmt", "text"}, tt.args...)

			stdout, _, err := runGitcat(t, append(args, dir)...)

			assert.NotErr(t, err)
			assert.Eq(t, stdout, tt.want)
		})
	}
}

func TestCat_GitRepo_RespectsGitignore(t *testing.T) {
	testutil.RequireGit(t)
	dir := testutil.GitRepo(t, map[string]string{
		".gitignore":   "build/\n",
		"main.go":      "package main\n",
		"build/gen.go": "package build\n",
	})

	stdout, _, err := runGitcat(t, "-fmt", "jsonl", dir)

	assert.NotErr(t, err)
	assert.Eq(t, stdout, `{"file":"main.go","ext":".go","content":"package main\n"}`+"\n")
}

func TestCat_OutFile(t *testing.T) {
	dir := sampleRepo(t)
	out := filepath.Join(t.TempDir(), "context")

	stdout, stderr, err := runGitcat(t, "-fmt", "text", "-keep", "go", "-out", out, dir)

	assert.NotErr(t, err)
	assert.Eq(t, stdout, "")
	assert.Contains(t, stderr, "output written to file")
	data, err := os.ReadFile(out + ".text")
	assert.NotErr(t, err)
	assert.Eq(t, string(data), "package main\n\n")
}

func TestCat_DryRun_DoesNothing(t *testing.T) {
	t.Chdir(t.TempDir())

	stdout, stderr, err := runGitcat(t, "-dryrun", "git@github.com:user/repo.git")

	assert.NotErr(t, err)
	assert.Eq(t, stdout, "")
	assert.ContainsAll(t, stderr, "would clone repository", "dir=repo")
	assert.NoPathExists(t, "repo")
}

func TestCat_TmpClone_RemovesDirectory(t *testing.T) {
	testutil.RequireGit(t)
	src := testutil.GitRepo(t, map[string]string{"main.go": "package main\n"})
	testutil.Commit(t, src)
	missing := filepath.Join(t.TempDir(), "missing")
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	tests := []struct {
		name    string
		url     string
		wantOut string
		wantErr string
	}{
		{
			"success",
			"file://" + src,
			`{"file":"main.go","ext":".go","content":"package main\n"}` + "\n",
			"",
		},
		{
			"failed_clone",
			"file://" + missing,
			"",
			"git clone file://" + missing,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCat()
			c.tmpClone = true
			c.location = &gitpath.GitPath{Path: tt.url, Kind: gitpath.SSH}
			e, stdout, _ := newTestEnv()

			err := c.Run(e)

			if tt.wantErr != "" {
				assert.ErrContains(t, err, tt.wantErr)
			} else {
				assert.NotErr(t, err)
			}
			assert.Eq(t, stdout.String(), tt.wantOut)
			entries, err := os.ReadDir(tmp)
			assert.NotErr(t, err)
			assert.Len(t, entries, 0)
		})
	}
}
