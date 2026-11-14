package files

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/i-zaitsev/gitcat/pkg/repo"
	"github.com/i-zaitsev/pkg/assert"
)

func TestLoad(t *testing.T) {
	fsys := fstest.MapFS{
		"b.go":      {Data: []byte("package b")},
		"a.go":      {Data: []byte("package a\n")},
		"LICENSE":   {Data: []byte("MIT\n")},
		"docs/x.md": {Data: []byte("# X\n")},
		"empty.txt": {Data: []byte("")},
		"logo.png":  {Data: []byte("\x89PNG\r\n\x1a\n\x00\x00")},
	}
	r := repo.New(fsys, []string{"docs/x.md", "b.go", "LICENSE", "logo.png", "a.go", "missing.go", "empty.txt"})

	got := Load(r, 0)

	assert.Eq(t, got, []File{
		{Path: "LICENSE", Ext: "", Content: "MIT\n"},
		{Path: "a.go", Ext: ".go", Content: "package a\n"},
		{Path: "b.go", Ext: ".go", Content: "package b\n"},
		{Path: "docs/x.md", Ext: ".md", Content: "# X\n"},
		{Path: "empty.txt", Ext: ".txt", Content: ""},
	})
}

func TestLoad_Empty(t *testing.T) {
	got := Load(repo.New(fstest.MapFS{}, nil), 0)

	assert.Len(t, got, 0)
}

func TestLoad_Head(t *testing.T) {
	fsys := fstest.MapFS{"a.go": {Data: []byte("package a\n\nfunc A() {}\n")}}

	got := Load(repo.New(fsys, []string{"a.go"}), 1)

	assert.Eq(t, got, []File{{Path: "a.go", Ext: ".go", Content: "package a\n"}})
}

func TestLoad_LongLine_IsKept(t *testing.T) {
	long := strings.Repeat("x", 100_000) + "\n"
	fsys := fstest.MapFS{"min.js": {Data: []byte(long)}}

	got := Load(repo.New(fsys, []string{"min.js"}), 0)

	assert.Len(t, got, 1)
	assert.Eq(t, got[0].Content, long)
}

func TestHeadLines(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{
			"all_lines",
			"1\n2\n3\n",
			0,
			"1\n2\n3\n",
		},
		{
			"first_two",
			"1\n2\n3\n",
			2,
			"1\n2\n",
		},
		{
			"more_than_available",
			"1\n2\n",
			5,
			"1\n2\n",
		},
		{
			"adds_final_newline",
			"1\n2",
			0,
			"1\n2\n",
		},
		{
			"cut_before_unterminated_line",
			"1\n2",
			1,
			"1\n",
		},
		{
			"empty",
			"",
			3,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Eq(t, headLines(tt.s, tt.n), tt.want)
		})
	}
}

func TestIsBinary(t *testing.T) {
	assert.True(t, isBinary([]byte("PK\x03\x04\x00")))
	assert.False(t, isBinary([]byte("plain text\n")))
	assert.False(t, isBinary([]byte(strings.Repeat("a", binaryProbe)+"\x00")))
}
