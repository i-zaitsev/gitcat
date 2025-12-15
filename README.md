# gitcat

A CLI tool that concatenates git repository contents into a single output file. Supports both local and remote repositories with flexible output formats.

## Features

- Clone and process remote git repositories (SSH and HTTPS)
- Process local git repositories
- Multiple output formats (JSON and text)
- File extension-based grouping
- Temporary clone support with automatic cleanup
- Colored logging output
- Dry-run mode for testing

## Installation

```bash
go install github.com/i-zaitsev/gitcat/cmd/gitcat@latest
```

Or build from source:

```bash
git clone https://github.com/i-zaitsev/gitcat.git
cd gitcat
go build -o gitcat ./cmd/gitcat
```

## Usage

```bash
gitcat [options] <repository-url>
```

### Examples

Clone and concatenate a remote repository via SSH:
```bash
gitcat git@github.com:user/repo.git
```

Clone and concatenate via HTTPS:
```bash
gitcat https://github.com/user/repo.git
```

Process a local repository:
```bash
gitcat /path/to/local/repo
```

Clone to a specific directory:
```bash
gitcat -dir ./myrepo git@github.com:user/repo.git
```

Use temporary clone with automatic cleanup:
```bash
gitcat -tmp https://github.com/user/repo.git
```

Output in text format instead of JSON:
```bash
gitcat -fmt text git@github.com:user/repo.git
```

Dry run to see what would happen:
```bash
gitcat -dryrun -dir ./myrepo git@github.com:user/repo.git
```

Enable debug logging:
```bash
gitcat -debug git@github.com:user/repo.git
```

## Command-line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-dryrun` | false | Dry run mode - log actions without executing them |
| `-debug` | false | Enable debug logging |
| `-tmp` | false | Clone into a temporary directory which is deleted after execution |
| `-fmt` | json | Output format: `json` or `text` |
| `-dir` | (repo name) | Local directory to clone into |

## Output Formats

### JSON Format

The JSON output includes:
- `files`: Array of all file paths in the repository
- File extension keys (e.g., `.go`, `.md`, `.json`): Arrays of concatenated content for each file type

Example output structure:
```json
{
  "files": ["main.go", "README.md", "LICENSE"],
  ".go": ["package main\n", "import \"fmt\"\n"],
  ".md": ["# Project\n", "## Description\n"]
}
```

### Text Format

The text format outputs all file contents concatenated together, grouped by file extension.

## Supported Repository Types

- **SSH**: `git@github.com:user/repo.git`
- **HTTPS**: `https://github.com/user/repo.git`
- **Local**: `/path/to/local/repository`

## Requirements

- Go 1.25 or later
- Git (for cloning remote repositories)
- SSH keys configured (for SSH repositories)

## License

MIT License - see [LICENSE](LICENSE) file for details.
