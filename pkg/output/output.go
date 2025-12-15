package output

import (
	"encoding/json"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/ls"
)

const (
	Text = "text"
	JSON = "json"
)

func ToText(repo *ls.RepoContent) string {
	allFileExt := files.DiscoverExt(repo)
	var buf strings.Builder
	for _, ext := range allFileExt {
		extRepo := files.MatchExt(repo, ext)
		for _, line := range files.Cat(extRepo.Files) {
			buf.WriteString(line + "\n")
		}
	}
	return buf.String()
}

func ToJSON(repo *ls.RepoContent) string {
	allFileExt := files.DiscoverExt(repo)
	output := make(map[string][]string, len(allFileExt)+1)
	output["files"] = repo.Files
	for _, ext := range allFileExt {
		extRepo := files.MatchExt(repo, ext)
		for _, line := range files.Cat(extRepo.Files) {
			output[ext] = append(output[ext], line)
		}
	}
	content, err := json.Marshal(output)
	if err != nil {
		log.Error("failed to marshal output", "error", err)
		return ""
	}
	return string(content)
}
