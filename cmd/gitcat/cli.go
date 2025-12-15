package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/gitpath"
	"github.com/i-zaitsev/gitcat/pkg/log"
)

type Cli struct {
	location *gitpath.GitPath
	localDir string
	dryRun   bool
	debug    bool
	tmpClone bool
}

func NewCLI() *Cli {
	return &Cli{}
}

// Parse takes args without the program's name and parses the flags.
func (c *Cli) Parse(args []string) error {
	fs := flag.NewFlagSet("gitcat", flag.ContinueOnError)

	fs.Usage = c.usage(fs)

	fs.BoolVar(&c.dryRun, "dryrun", false, "dry run mode - log actions without executing them")
	fs.BoolVar(&c.debug, "debug", false, "enable debug logging")
	fs.BoolVar(&c.tmpClone, "tmp", false, "clone into a temporary directory which is deleted after execution")
	fs.StringVar(&c.localDir, "dir", "", "local directory to clone into (defaults to repo name)")

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
