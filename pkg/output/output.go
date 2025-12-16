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
		buf.WriteString(files.Cat(extRepo.Files...) + "\n")
	}
	return buf.String()
}

func ToJSON(repo *ls.RepoContent) string {
	ctx := makeContext(repo)
	content, err := json.Marshal(ctx)
	if err != nil {
		log.Error("failed to marshal output", "error", err)
		return ""
	}
	return string(content)
}

// context represents a code repository in a structured form
// Each path is mapped to its content.
// Additionally, metadata such as listed file names and extensions are included.
type context struct {
	Files         []string          `json:"files"`
	Exts          []string          `json:"exts"`
	PathToContent map[string]string `json:"pathToContent"`
}

func makeContext(repo *ls.RepoContent) context {
	allFileExt := files.DiscoverExt(repo)
	ctx := context{
		Files:         repo.Files,
		Exts:          allFileExt,
		PathToContent: make(map[string]string, len(repo.Files)),
	}
	for _, ext := range allFileExt {
		extRepo := files.MatchExt(repo, ext)
		for _, filename := range extRepo.Files {
			ctx.PathToContent[filename] = files.Cat(filename)
		}
	}
	return ctx
}
