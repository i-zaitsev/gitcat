# gitcat

A CLI tool that concatenates the files of a git repository into a single document, ready to be loaded into an LLM's
context. Works with local directories and remote repositories, and writes plain text, JSON Lines, or Markdown.

## Features

- Shallow-clones remote repositories over SSH or HTTPS
- Reads local directories and git work trees
- Respects `.gitignore` in git work trees
- Skips binary files and dotfiles (opt in with `-dot`)
- Filters by path, extension, and file size, and can keep only the first lines of each file
- Clones into a temporary directory with automatic cleanup
- Dry-run mode and colored logging

## Installation

Download the archive for your platform (linux or darwin, amd64 or arm64) from the
[releases](https://github.com/i-zaitsev/gitcat/releases) page, check it against `checksums.txt`, and put the
`gitcat` binary on your `PATH`.

Or install with Go:

```bash
go install github.com/i-zaitsev/gitcat/cmd/gitcat@latest
```

Or build from source:

```bash
git clone https://github.com/i-zaitsev/gitcat.git
cd gitcat
make build      # writes bin/gitcat
```

Check the installed version with `gitcat version`.

## Usage

```bash
gitcat [flags] <repo>
gitcat <command> [flags] [args]
gitcat help
```

`cat` is the default command, so `gitcat [flags] <repo>` is the same as `gitcat cat [flags] <repo>`.

### Examples

The examples read gitcat's own repository, so they can be run as they are.

```bash
# clone via SSH into ./gitcat
gitcat git@github.com:i-zaitsev/gitcat.git

# temporary clone, Markdown output
gitcat -tmp -fmt md https://github.com/i-zaitsev/gitcat

# current directory, for example a checkout of gitcat
gitcat -fmt text .

# only some paths
gitcat -tmp -path pkg/files,cmd -exclude testdata https://github.com/i-zaitsev/gitcat

# Go files up to 500 KB
gitcat -tmp -maxsize 500 -keep .go https://github.com/i-zaitsev/gitcat

# first 50 lines of each file, written to gitcat.jsonl
gitcat -tmp -head 50 -out gitcat https://github.com/i-zaitsev/gitcat

# log what would happen
gitcat -dryrun -dir ./gitcat-src git@github.com:i-zaitsev/gitcat.git
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-fmt` | `jsonl` | Output format: `text`, `jsonl`, or `md` |
| `-out` | stdout | Output file without extension; the format is appended, e.g. `context.md` |
| `-tmp` | false | Clone into a temporary directory which is deleted after execution |
| `-dir` | repo name | Local directory to clone into; cannot be combined with `-tmp` |
| `-keep` | all | Comma-separated extensions to keep, e.g. `go,md` |
| `-path` | all | Comma-separated paths to include |
| `-exclude` | none | Comma-separated paths to exclude |
| `-minsize` | 0 | Minimum file size in KB |
| `-maxsize` | none | Maximum file size in KB |
| `-head` | 0 | Number of lines to keep from each file (0 = all) |
| `-dot` | false | Include dotfiles and dot directories, such as `.github` |
| `-dryrun` | false | Log actions without executing them |
| `-debug` | false | Enable debug logging |

`-keep`, `-path`, and `-exclude` can be repeated; the values add up.

## Which files are included

- In a git work tree, the files come from `git ls-files`: tracked files plus untracked files that are not ignored.
  Elsewhere, the directory is walked.
- Only regular files are read; symbolic links and submodules are skipped.
- Binary files (a NUL byte in the first 8000 bytes) are skipped.
- Files are grouped by extension and sorted by path within a group; files without an extension come first.
- Paths in the output are relative to the repository root.

## Output formats

### JSON Lines (`jsonl`)

One object per file:

```json
{"file":"cmd/gitcat/gitcat.go","ext":".go","content":"package main\n..."}
```

### Markdown (`md`)

A section per file with its path and the content in a fenced code block tagged with the extension.
The fence grows when the content itself contains code fences.

### Text (`text`)

The contents of all files concatenated, with an empty line after each extension group.

## Supported repositories

- **SSH**: `git@github.com:i-zaitsev/gitcat.git`, `ssh://git@github.com/i-zaitsev/gitcat.git`
- **HTTPS**: `https://github.com/i-zaitsev/gitcat.git`, `https://github.com/i-zaitsev/gitcat`
- **Local**: `/path/to/local/repository`, `.`

## Development

```bash
make test       # go test ./...
make lint       # go fmt, go vet, golangci-lint
make build      # bin/gitcat
make snapshot   # release archives in dist/, nothing is published
```

`make lint` runs golangci-lint with the Go toolchain declared in `go.mod`.
`make build` stamps the output of `git describe` as the version; `make snapshot` needs [GoReleaser](https://goreleaser.com).

## Releasing

A release is a `v*` tag. The [Release workflow](.github/workflows/release.yml) runs GoReleaser, which builds the
archives for linux and darwin on amd64 and arm64, writes `checksums.txt`, and publishes a GitHub release.
The same tag is the version that `go install github.com/i-zaitsev/gitcat/cmd/gitcat@v0.1.0` installs.

1. Add the changes to [CHANGELOG.md](CHANGELOG.md) and commit them as `docs: release v0.1.0`.
2. Tag and push:

   ```bash
   git tag v0.1.0
   git push origin main v0.1.0
   ```

Versions follow [semantic versioning](https://semver.org). From v2 on, the module path in `go.mod` and the import
paths must end in `/v2`. CI dry-runs the release on every pull request and push to main.

## Requirements

- Go 1.25 or later
- Git
- SSH keys configured (for SSH repositories)

## License

MIT License - see [LICENSE](LICENSE) file for details.
