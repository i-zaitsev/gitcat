package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/ls"
)

const (
	FormatText     = "text"
	FormatJSONL    = "jsonl"
	FormatMarkdown = "md"
)

// Format represents an output format and implements flag.Value interface.
type Format string

// String returns the string representation of the format.
func (f *Format) String() string {
	if f == nil {
		return FormatJSONL
	}
	return string(*f)
}

// Set validates and sets the format value.
func (f *Format) Set(value string) error {
	switch value {
	case FormatText, FormatJSONL, FormatMarkdown:
		*f = Format(value)
		return nil
	default:
		return fmt.Errorf("invalid format %q: must be one of: text, jsonl, md", value)
	}
}

func ToText(repo *ls.RepoContent) string {
	var buf strings.Builder
	for _, ext := range files.DiscoverExt(repo) {
		extRepo := files.MatchExt(repo, ext)
		buf.WriteString(files.Cat(extRepo.Files...) + "\n")
	}
	return buf.String()
}

func ToJSONL(repo *ls.RepoContent) (string, error) {
	var buf strings.Builder

	for _, ext := range files.DiscoverExt(repo) {
		extRepo := files.MatchExt(repo, ext)
		for _, filename := range extRepo.Files {
			entry := outputEntry{
				File:    filename,
				Ext:     ext,
				Content: files.Cat(filename),
			}
			content, err := json.Marshal(entry)
			if err != nil {
				log.Error("failed to marshal output entry", "entry", entry, "error", err)
				continue
			}
			buf.Write(content)
			buf.WriteRune('\n')
		}
	}

	return buf.String(), nil
}

// ToMarkdown formats repository content as markdown with code blocks.
// Uses the same file iteration as JSONL but outputs markdown format.
func ToMarkdown(repo *ls.RepoContent) string {
	var buf strings.Builder

	for _, ext := range files.DiscoverExt(repo) {
		extRepo := files.MatchExt(repo, ext)
		for _, filename := range extRepo.Files {
			content := files.Cat(filename)

			buf.WriteString("## ")
			buf.WriteString(filename)
			buf.WriteString("\n")
			buf.WriteString("*Extension: ")
			buf.WriteString(ext)
			buf.WriteString("*\n\n")
			buf.WriteString("```")
			buf.WriteString(strings.TrimPrefix(ext, "."))
			buf.WriteString("\n")
			buf.WriteString(content)
			if !strings.HasSuffix(content, "\n") {
				buf.WriteString("\n")
			}
			buf.WriteString("```\n\n")
			buf.WriteString("---\n\n")
		}
	}

	return buf.String()
}

type outputEntry struct {
	File    string `json:"file"`
	Ext     string `json:"ext"`
	Content string `json:"content"`
}
