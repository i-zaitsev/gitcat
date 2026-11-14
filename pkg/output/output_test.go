package output

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/pkg/assert"
)

var sample = []files.File{
	{Path: "LICENSE", Ext: "", Content: "MIT\n"},
	{Path: "a.go", Ext: ".go", Content: "package a\n"},
	{Path: "pkg/b.go", Ext: ".go", Content: "package b\n"},
	{Path: "README.md", Ext: ".md", Content: "# T\n\n```go\nx := 1\n```\n"},
}

func TestFormat_Set(t *testing.T) {
	var f Format

	assert.NotErr(t, f.Set("md"))
	assert.Eq(t, f, FormatMarkdown)
	assert.ErrContains(t, f.Set("xml"), `invalid format "xml"`)
	assert.Eq(t, f, FormatMarkdown)
}

func TestFormat_String_Nil(t *testing.T) {
	var f *Format

	assert.Eq(t, f.String(), FormatJSONL)
}

func TestWrite_Text(t *testing.T) {
	var buf strings.Builder

	err := Write(&buf, FormatText, sample)

	assert.NotErr(t, err)
	assert.Eq(t, buf.String(), "MIT\n\npackage a\npackage b\n\n# T\n\n```go\nx := 1\n```\n\n")
}

func TestWrite_JSONL(t *testing.T) {
	var buf strings.Builder

	err := Write(&buf, FormatJSONL, []files.File{
		{Path: "a.go", Ext: ".go", Content: "package a\n"},
		{Path: "web/i.html", Ext: ".html", Content: "<p>&</p>\n"},
	})

	assert.NotErr(t, err)
	assert.Eq(t, buf.String(),
		`{"file":"a.go","ext":".go","content":"package a\n"}`+"\n"+
			`{"file":"web/i.html","ext":".html","content":"<p>&</p>\n"}`+"\n")
}

func TestWrite_JSONL_RoundTrip(t *testing.T) {
	var buf strings.Builder
	assert.NotErr(t, Write(&buf, FormatJSONL, sample))

	dec := json.NewDecoder(strings.NewReader(buf.String()))
	var got []files.File
	for {
		var e entry
		err := dec.Decode(&e)
		if errors.Is(err, io.EOF) {
			break
		}
		assert.NotErr(t, err)
		got = append(got, files.File{Path: e.File, Ext: e.Ext, Content: e.Content})
	}

	assert.Eq(t, got, sample)
}

func TestWrite_Markdown(t *testing.T) {
	var buf strings.Builder

	err := Write(&buf, FormatMarkdown, sample)

	assert.NotErr(t, err)
	assert.Eq(t, buf.String(), ""+
		"## LICENSE\n\n```\nMIT\n```\n\n---\n\n"+
		"## a.go\n*Extension: .go*\n\n```go\npackage a\n```\n\n---\n\n"+
		"## pkg/b.go\n*Extension: .go*\n\n```go\npackage b\n```\n\n---\n\n"+
		"## README.md\n*Extension: .md*\n\n````md\n# T\n\n```go\nx := 1\n```\n````\n\n---\n\n")
}

func TestWrite_Markdown_UnterminatedContent(t *testing.T) {
	var buf strings.Builder

	err := Write(&buf, FormatMarkdown, []files.File{{Path: "a.go", Ext: ".go", Content: "x"}})

	assert.NotErr(t, err)
	assert.Eq(t, buf.String(), "## a.go\n*Extension: .go*\n\n```go\nx\n```\n\n---\n\n")
}

func TestWrite_UnknownFormat(t *testing.T) {
	err := Write(io.Discard, Format("xml"), sample)

	assert.ErrContains(t, err, `unknown output format "xml"`)
}

func TestWrite_NoFiles(t *testing.T) {
	for _, f := range []Format{FormatText, FormatJSONL, FormatMarkdown} {
		var buf strings.Builder

		assert.NotErr(t, Write(&buf, f, nil))
		assert.Eq(t, buf.String(), "")
	}
}

func TestFenceFor(t *testing.T) {
	tests := []struct {
		content string
		want    string
	}{
		{
			"plain\n",
			"```",
		},
		{
			"`inline` code\n",
			"```",
		},
		{
			"```go\n```\n",
			"````",
		},
		{
			"`````\n",
			"``````",
		},
	}
	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			assert.Eq(t, fenceFor(tt.content), tt.want)
		})
	}
}
