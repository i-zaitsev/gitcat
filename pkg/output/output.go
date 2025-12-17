package output

import (
	"encoding/json"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/ls"
)

const (
	FormatText  = "text"
	FormatJSONL = "jsonl"
)

func ToText(repo *ls.RepoContent) string {
	allFileExt := files.DiscoverExt(repo)
	var buf strings.Builder
	for _, ext := range allFileExt {
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

type outputEntry struct {
	File    string `json:"file"`
	Ext     string `json:"ext"`
	Content string `json:"content"`
}
