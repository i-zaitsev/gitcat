// Package output renders repository files in the formats gitcat supports:
// plain text, JSON Lines, and Markdown.
//
// The formatters only write to an [io.Writer]. Reading the files and opening
// the destination are up to the caller.
package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/files"
)

// Supported output formats.
const (
	FormatText     = "text"
	FormatJSONL    = "jsonl"
	FormatMarkdown = "md"
)

// Format is an output format. It implements [flag.Value].
type Format string

// String returns the name of the format, or the default format [FormatJSONL]
// for a nil Format.
func (f *Format) String() string {
	if f == nil {
		return FormatJSONL
	}
	return string(*f)
}

// Set validates value and sets the format.
func (f *Format) Set(value string) error {
	switch value {
	case FormatText, FormatJSONL, FormatMarkdown:
		*f = Format(value)
		return nil
	default:
		return fmt.Errorf("invalid format %q: must be one of: text, jsonl, md", value)
	}
}

// Write renders the files to w in the format f.
func Write(w io.Writer, f Format, loaded []files.File) error {
	switch f {
	case FormatText:
		return writeText(w, loaded)
	case FormatJSONL:
		return writeJSONL(w, loaded)
	case FormatMarkdown:
		return writeMarkdown(w, loaded)
	default:
		return fmt.Errorf("unknown output format %q", string(f))
	}
}

// writeText writes the contents of the files one after another, with an
// empty line after each group of files that share an extension.
func writeText(w io.Writer, loaded []files.File) error {
	bw := bufio.NewWriter(w)
	for i, f := range loaded {
		_, _ = bw.WriteString(f.Content)
		if i == len(loaded)-1 || loaded[i+1].Ext != f.Ext {
			_ = bw.WriteByte('\n')
		}
	}
	return bw.Flush()
}

// entry is a line of the JSON Lines output.
type entry struct {
	File    string `json:"file"`
	Ext     string `json:"ext"`
	Content string `json:"content"`
}

// writeJSONL writes each file as a JSON object on its own line.
// HTML characters in the content are not escaped.
func writeJSONL(w io.Writer, loaded []files.File) error {
	bw := bufio.NewWriter(w)
	enc := json.NewEncoder(bw)
	enc.SetEscapeHTML(false)
	for _, f := range loaded {
		if err := enc.Encode(entry{File: f.Path, Ext: f.Ext, Content: f.Content}); err != nil {
			return fmt.Errorf("encode %s: %w", f.Path, err)
		}
	}
	return bw.Flush()
}

// writeMarkdown writes a section for each file: a heading with the path,
// the extension if any, and the content in a fenced code block tagged with
// the extension. The fence is longer than any run of backticks in the
// content, so code blocks inside the content do not end it early.
func writeMarkdown(w io.Writer, loaded []files.File) error {
	bw := bufio.NewWriter(w)
	for _, f := range loaded {
		_, _ = fmt.Fprintf(bw, "## %s\n", f.Path)
		if f.Ext != "" {
			_, _ = fmt.Fprintf(bw, "*Extension: %s*\n", f.Ext)
		}
		fence := fenceFor(f.Content)
		content := f.Content
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		_, _ = fmt.Fprintf(bw, "\n%s%s\n%s%s\n\n---\n\n", fence, strings.TrimPrefix(f.Ext, "."), content, fence)
	}
	return bw.Flush()
}

// fenceFor returns a code fence for content: at least three backticks, and
// more than the longest run of backticks in content.
func fenceFor(content string) string {
	longest, run := 0, 0
	for i := range len(content) {
		if content[i] == '`' {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	return strings.Repeat("`", max(3, longest+1))
}
