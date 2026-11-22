# Changelog

## v0.2.0 (unreleased)

### Features

- `-preset go|python|c|cpp|all` selects the sources and docs of a project and skips its build artifacts and caches
- `-keep` and `-ignore` take file names and globs like `go.mod`, `*.py[cod]`, or `__pycache__`

### Breaking

- `-path` is now `-only`, and `-exclude` is now `-ignore`
- `-ignore` without a slash matches at any depth: `-ignore testdata` also skips `pkg/x/testdata`
- `-keep go` matches a file named `go`: use `-keep .go` for the extension
- Library API: `filter.Filter` has `Only`, `Ignore`, and `Keep` in place of `WithPaths`, `ExcludePaths`, and `KeepExt`

## v0.1.0 (2026-11-10)

The first tagged release.

### Features

- `gitcat [flags] <repo>` concatenates a local directory or a remote repository into one document: `text`, `jsonl`, or `md`
- `gitcat help` and `gitcat version`
- `.gitignore` is respected in git work trees
- Remote repositories are cloned shallow; with `-tmp`, into a temporary directory that is always removed
- Binary files are skipped, and files without an extension are included
- `-dot` includes dotfiles and dot directories
- `-path`, `-exclude`, and `-keep` take comma-separated lists and can be repeated
- `ssh://` URLs and https URLs without the `.git` suffix
- Release archives for linux and darwin on amd64 and arm64, with checksums

### Fixes

- Files with lines longer than 64 KB are read in full
- Markdown output stays valid when a file contains code fences
- A failed clone shows the error reported by git

### Breaking

- Output paths are relative to the repository root
- `-dir` and `-tmp` cannot be combined
- A local directory named `cat`, `help`, or `version` must be passed as `./name`
- Library API: `ls` is replaced by `repo` (the model of a repository), `source` (reading it from disk), and
  `filter` (selecting files); `gitclone` is now `git`; `files.Load` and `output.Write` replace `files.Cat` and
  `output.ToText`, `ToJSONL`, `ToMarkdown`

### Other

- golangci-lint with depguard rules for the package layering
- CI runs lint, tests, build, and a release dry run
- Releases are built by GoReleaser from `v*` tags
