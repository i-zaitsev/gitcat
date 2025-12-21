package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/gitpath"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/output"
)

type Cli struct {
	location     *gitpath.GitPath
	localDir     string
	outFile      string
	dryRun       bool
	debug        bool
	tmpClone     bool
	outFmt       output.Format
	keepExt      gitpath.Extensions
	includePaths gitpath.Paths
	excludePaths gitpath.Paths
	minSize      gitpath.Size
	maxSize      gitpath.Size
	headLines    int
}

func NewCLI() *Cli {
	return &Cli{
		outFmt:  output.FormatJSONL,
		maxSize: -1,
	}
}

// Parse takes args without the program's name and parses the flags.
func (c *Cli) Parse(args []string) error {
	fs := flag.NewFlagSet("gitcat", flag.ContinueOnError)

	fs.Usage = c.usage(fs)

	fs.BoolVar(&c.dryRun, "dryrun", false, "dry run mode - log actions without executing them")
	fs.BoolVar(&c.debug, "debug", false, "enable debug logging")
	fs.BoolVar(&c.tmpClone, "tmp", false, "clone into a temporary directory which is deleted after execution")
	fs.StringVar(&c.outFile, "out", "", "output file (without extension, uses -fmt for extension)")
	fs.StringVar(&c.localDir, "dir", "", "local directory to clone into (defaults to repo name)")
	fs.Var(&c.outFmt, "fmt", "output format (text, jsonl, or md)")
	fs.Var(&c.keepExt, "keep", "comma-separated list of file extensions to keep (default: none)")
	fs.Var(&c.includePaths, "path", "comma-separated paths to include")
	fs.Var(&c.excludePaths, "exclude", "comma-separated paths to exclude")
	fs.Var(&c.minSize, "minsize", "minimum file size in KB (e.g., 100)")
	fs.Var(&c.maxSize, "maxsize", "maximum file size in KB (e.g., 500)")
	fs.IntVar(&c.headLines, "head", 0, "number of lines to read from each file (0 = all)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	c.setLog()

	remaining := fs.Args()
	if len(remaining) == 0 {
		fs.Usage()
		return fmt.Errorf("repository URL is required")
	}

	if len(remaining) > 1 {
		fs.Usage()
		return fmt.Errorf("too many arguments")
	}

	if location, err := gitpath.Parse(remaining[0]); err != nil {
		return err
	} else {
		c.location = location
	}

	if c.localDir == "" {
		c.localDir = c.inferLocalDir(c.location.Path)
	}

	return nil
}

func (c *Cli) usage(fs *flag.FlagSet) func() {
	return func() {
		old := fs.Output()
		var b strings.Builder
		fs.SetOutput(&b)
		b.WriteString("gitcat - concatenates a git repo into a single file\n\n")
		b.WriteString("usage: gitcat [options] <repository-url>\n\n")
		b.WriteString("options:\n")
		fs.PrintDefaults()
		b.WriteString("\nexamples:\n")
		b.WriteString("  gitcat git@github.com:user/repo.git\n")
		b.WriteString("  gitcat https://github.com/user/repo.git\n")
		b.WriteString("  gitcat -dryrun -dir ./myrepo git@github.com:user/repo.git\n")
		b.WriteString("  gitcat -path pkg/files,cmd -exclude testdata https://github.com/user/repo.git\n")
		b.WriteString("  gitcat -maxsize 500 -keep .go https://github.com/user/repo.git\n")
		b.WriteString("  gitcat -head 50 https://github.com/user/repo.git\n")
		fs.SetOutput(old)
		_, _ = fmt.Fprintln(fs.Output(), b.String())
	}
}

func (c *Cli) inferLocalDir(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	if len(parts) == 0 {
		return "repo"
	}
	last := parts[len(parts)-1]
	return strings.TrimSuffix(last, ".git")
}

func (c *Cli) setLog() {
	var logLevel slog.Level

	if c.debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(
		slog.NewTextHandler(
			log.NewColorWriter(os.Stderr),
			&slog.HandlerOptions{
				Level: logLevel,
			},
		),
	).With("app", "gitcat")

	log.SetLogger(logger)
}
