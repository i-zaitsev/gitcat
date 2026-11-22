# gitcat

Concatenates the files of a git repository into a single document. Works with local
directories and remote repositories. Writes plain text, JSON Lines, or Markdown to
stdout or a file.

## Features

* Shallow-clones remote repositories over SSH or HTTPS
* Reads local directories and git work trees
* Respects `.gitignore` in git work trees
* Skips binary files and dotfiles (opt in with `-dot`)
* Filters by path, extension, and file size, and can keep only the first lines of each file
* Clones into a temporary directory with automatic cleanup
* Dry-run mode and colored logging

## Requirements

* Go 1.25 or later
* Git
* SSH keys configured (for SSH repositories)

## Installation

With Go:
```bash
go install github.com/i-zaitsev/gitcat/cmd/gitcat@latest
```

From source:
```bash
git clone https://github.com/i-zaitsev/gitcat.git
cd gitcat
make build  # writes bin/gitcat
```

Check the installed version with `gitcat version`.

## Usage

```bash
gitcat [flags] <repo>
gitcat <command> [flags] [args]
gitcat help
```

## Examples

```bash
# clone via SSH into ./gitcat
gitcat git@github.com:i-zaitsev/gitcat.git

# temporary clone, Markdown output
gitcat -tmp -fmt md https://github.com/i-zaitsev/gitcat

# current directory
gitcat -fmt text .

# remote, with only some paths included and testdata excluded
gitcat -tmp -path pkg/files,cmd -exclude testdata https://github.com/i-zaitsev/gitcat

# Go files up to 500 KB
gitcat -tmp -maxsize 500 -keep .go https://github.com/i-zaitsev/gitcat

# take first 50 lines of each file, written to gitcat.jsonl
gitcat -tmp -head 50 -out gitcat https://github.com/i-zaitsev/gitcat

# only log what would happen when the command is executed
gitcat -dryrun -dir ./gitcat-src git@github.com:i-zaitsev/gitcat.git
```

See `gitcat -h` for more options.

## Notes

* In a git work tree, the files come from `git ls-files` 
* Elsewhere, the directory is walked
* Only regular files are read: no symbolic links or submodules
* Binary files (a NUL byte in the first 8000 bytes) are skipped
* Files are grouped by extension and sorted by path within a group
* Files without an extension come first
* Paths in the output are relative to the repository root

## License

MIT. See [LICENSE](LICENSE).
