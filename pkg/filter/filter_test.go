package filter

import (
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/i-zaitsev/gitcat/pkg/repo"
	"github.com/i-zaitsev/pkg/assert"
)

// newRepo returns a Repo with every file of a map from paths to contents.
func newRepo(t *testing.T, files map[string]string) *repo.Repo {
	t.Helper()
	fsys := fstest.MapFS{}
	for p, content := range files {
		fsys[p] = &fstest.MapFile{Data: []byte(content)}
	}
	r, err := repo.Scan(fsys)
	assert.NotErr(t, err)
	return r
}

var tree = map[string]string{
	".env":           "SECRET=1\n",
	".github/ci.yml": "name: ci\n",
	"README":         "readme\n",
	"main.go":        "package main\n",
	"pkg/a/a.go":     "package a\n",
	"pkg/b/b.go":     "package b\n",
	"testdata/x.txt": "x\n",
}

func TestApply_PathRules(t *testing.T) {
	tests := []struct {
		name   string
		filter *Filter
		want   []string
	}{
		{
			"everything",
			New(),
			[]string{".env", ".github/ci.yml", "README", "main.go", "pkg/a/a.go", "pkg/b/b.go", "testdata/x.txt"},
		},
		{
			"ignore_dot_files",
			New().IgnoreDotFiles(),
			[]string{"README", "main.go", "pkg/a/a.go", "pkg/b/b.go", "testdata/x.txt"},
		},
		{
			"include_paths",
			New().Only("pkg/a", "main.go"),
			[]string{"main.go", "pkg/a/a.go"},
		},
		{
			"exclude_paths",
			New().IgnoreDotFiles().Ignore("pkg", "testdata"),
			[]string{"README", "main.go"},
		},
		{
			"exclude_wins_over_include",
			New().Only("pkg").Ignore("pkg/b"),
			[]string{"pkg/a/a.go"},
		},
		{
			"include_is_not_a_name_prefix",
			New().Only("pkg/a/a"),
			nil,
		},
	}
	r := newRepo(t, tree)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Apply(r)

			assert.Eq(t, slices.Collect(got.Paths()), tt.want)
		})
	}
}

func TestApply_PathRules_AreCleaned(t *testing.T) {
	tests := []struct {
		name   string
		filter *Filter
		want   []string
	}{
		{
			"include_trailing_slash",
			New().Only("pkg/"),
			[]string{"pkg/a.go"},
		},
		{
			"include_dot_slash",
			New().Only("./pkg"),
			[]string{"pkg/a.go"},
		},
		{
			"include_root",
			New().Only("."),
			[]string{"cmd/b.go", "pkg/a.go"},
		},
		{
			"exclude_trailing_slash",
			New().Ignore("pkg/"),
			[]string{"cmd/b.go"},
		},
	}
	r := newRepo(t, map[string]string{"pkg/a.go": "package pkg\n", "cmd/b.go": "package cmd\n"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Apply(r)

			assert.Eq(t, slices.Collect(got.Paths()), tt.want)
		})
	}
}

func TestApply_Keep_Extensions(t *testing.T) {
	r := newRepo(t, map[string]string{"a.go": "", "b.md": "", "c": "", "d.txt": "", "pkg/e.go": ""})

	got := New().Keep(".go", ".md", ".go").Apply(r)

	assert.Eq(t, slices.Collect(got.Paths()), []string{"a.go", "b.md", "pkg/e.go"})
}

var globTree = map[string]string{
	"a.py":                          "",
	"cmd/pkg/files/y.go":            "",
	"docs/a.md":                     "",
	"docs/sub/b.md":                 "",
	"go.mod":                        "",
	"pkg/__pycache__/a.cpython.pyc": "",
	"pkg/b.py":                      "",
	"pkg/b.pyo":                     "",
	"pkg/files/x.go":                "",
	"proj.egg-info/PKG-INFO":        "",
	"sub/go.mod":                    "",
	"tools/go.mod.bak":              "",
}

func TestApply_Keep_Patterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		want     []string
	}{
		{
			"exact_name_at_any_depth",
			[]string{"go.mod"},
			[]string{"go.mod", "sub/go.mod"},
		},
		{
			"name_glob",
			[]string{"*.py"},
			[]string{"a.py", "pkg/b.py"},
		},
		{
			"character_class",
			[]string{"*.py[co]"},
			[]string{"pkg/__pycache__/a.cpython.pyc", "pkg/b.pyo"},
		},
		{
			"anchored_path",
			[]string{"docs/*.md"},
			[]string{"docs/a.md"},
		},
		{
			"bare_word_is_a_name",
			[]string{"md"},
			nil,
		},
	}
	r := newRepo(t, globTree)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New().Keep(tt.patterns...).Apply(r)

			assert.Eq(t, slices.Collect(got.Paths()), tt.want)
		})
	}
}

func TestApply_Ignore_Patterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		drop     []string
	}{
		{
			"dir_name_at_any_depth",
			[]string{"__pycache__"},
			[]string{"pkg/__pycache__/a.cpython.pyc"},
		},
		{
			"character_class",
			[]string{"*.py[cod]"},
			[]string{"pkg/__pycache__/a.cpython.pyc", "pkg/b.pyo"},
		},
		{
			"dir_glob",
			[]string{"*.egg-info"},
			[]string{"proj.egg-info/PKG-INFO"},
		},
		{
			"anchored_path",
			[]string{"pkg/files"},
			[]string{"pkg/files/x.go"},
		},
		{
			"anchored_glob_covers_subdirs",
			[]string{"docs/*"},
			[]string{"docs/a.md", "docs/sub/b.md"},
		},
		{
			"trailing_slash_is_not_anchored",
			[]string{"files/"},
			[]string{"cmd/pkg/files/y.go", "pkg/files/x.go"},
		},
	}
	r := newRepo(t, globTree)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New().Ignore(tt.patterns...).Apply(r)

			want := slices.Collect(r.Paths())
			want = slices.DeleteFunc(want, func(p string) bool { return slices.Contains(tt.drop, p) })
			assert.Eq(t, slices.Collect(got.Paths()), want)
		})
	}
}

func TestApply_Ignore_WinsOverKeepAndOnly(t *testing.T) {
	r := newRepo(t, globTree)

	got := New().Only("pkg").Keep(".py", ".pyo").Ignore("__pycache__", "*.pyo").Apply(r)

	assert.Eq(t, slices.Collect(got.Paths()), []string{"pkg/b.py"})
}

func TestCheckPattern(t *testing.T) {
	assert.NotErr(t, CheckPattern("*.py[cod]"))
	assert.NotErr(t, CheckPattern("docs/*"))
	assert.ErrContains(t, CheckPattern("["), `invalid pattern "["`)
	assert.ErrContains(t, CheckPattern("a[b"), `invalid pattern "a[b"`)
}

func TestApply_SizeRange(t *testing.T) {
	fsys := fstest.MapFS{
		"small.txt": {Data: []byte(strings.Repeat("s", 10))},
		"big.txt":   {Data: []byte(strings.Repeat("b", 2000))},
	}
	r := repo.New(fsys, []string{"big.txt", "small.txt", "missing.txt"})
	tests := []struct {
		name     string
		min, max int64
		want     []string
	}{
		{
			"no_limits",
			0,
			-1,
			[]string{"big.txt", "missing.txt", "small.txt"},
		},
		{
			"min_size",
			1024,
			-1,
			[]string{"big.txt"},
		},
		{
			"max_size",
			0,
			1024,
			[]string{"small.txt"},
		},
		{
			"both_limits",
			100,
			1000,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New().SizeRange(tt.min, tt.max).Apply(r)

			assert.Eq(t, slices.Collect(got.Paths()), tt.want)
		})
	}
}

func TestApply_AllRules(t *testing.T) {
	r := newRepo(t, map[string]string{
		".github/ci.go": "package ci\n",
		"cmd/main.go":   "package main\n",
		"pkg/a.go":      "package a\n",
		"pkg/a.md":      "# A\n",
		"pkg/big.go":    strings.Repeat("x", 2000),
	})

	got := New().
		IgnoreDotFiles().
		Ignore("cmd").
		Keep(".go").
		SizeRange(0, 1024).
		Apply(r)

	assert.Eq(t, slices.Collect(got.Paths()), []string{"pkg/a.go"})
}
