package main

import (
	"errors"
	"flag"
	"io"
	"os"
	"strconv"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/gitcat/pkg/filter"
	"github.com/i-zaitsev/gitcat/pkg/gitpath"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/output"
	"github.com/i-zaitsev/gitcat/pkg/repo"
	"github.com/i-zaitsev/gitcat/pkg/source"
	"github.com/i-zaitsev/pkg/cmd"
)

// catCmd concatenates the files of a repository into a single document.
type catCmd struct {
	fs           *flag.FlagSet
	location     *gitpath.GitPath
	localDir     string
	outFile      string
	dryRun       bool
	debug        bool
	tmpClone     bool
	dotFiles     bool
	outFmt       output.Format
	keepExt      extensions
	includePaths paths
	excludePaths paths
	minSize      size
	maxSize      size
	headLines    int
}

// newCat returns a cat command with the default options and its flags defined.
func newCat() *catCmd {
	c := &catCmd{outFmt: output.FormatJSONL, maxSize: -1}
	c.fs = cmd.NewFlagSet(c)
	c.fs.BoolVar(&c.dryRun, "dryrun", false, "dry run mode: log actions without executing them")
	c.fs.BoolVar(&c.debug, "debug", false, "enable debug logging")
	c.fs.BoolVar(&c.tmpClone, "tmp", false, "clone into a temporary directory which is deleted after execution")
	c.fs.BoolVar(&c.dotFiles, "dot", false, "include dotfiles and dot directories")
	c.fs.StringVar(&c.outFile, "out", "", "output file without extension")
	c.fs.StringVar(&c.localDir, "dir", "", "local directory to clone into (defaults to repo name)")
	c.fs.Var(&c.outFmt, "fmt", "output format: text, jsonl, or md")
	c.fs.Var(&c.keepExt, "keep", "comma-separated list of file extensions to keep (default: all)")
	c.fs.Var(&c.includePaths, "path", "comma-separated paths to include")
	c.fs.Var(&c.excludePaths, "exclude", "comma-separated paths to exclude")
	c.fs.Var(&c.minSize, "minsize", "minimum file size in KB (e.g., 100)")
	c.fs.Var(&c.maxSize, "maxsize", "maximum file size in KB (e.g., 500)")
	c.fs.IntVar(&c.headLines, "head", 0, "number of lines to read from each file (0 = all)")
	return c
}

// Name returns the name of the command.
func (c *catCmd) Name() string { return "cat" }

// Desc returns a one-line description of the command.
func (c *catCmd) Desc() string {
	return "Concatenate the files of a repository into a single document"
}

// Args returns the positional arguments of the command.
func (c *catCmd) Args() string { return "<repo>" }

// PrintFlags writes the flags of the command and their defaults to w.
func (c *catCmd) PrintFlags(w io.Writer) {
	c.fs.SetOutput(w)
	c.fs.PrintDefaults()
	c.fs.SetOutput(io.Discard)
}

// Parse sets the options of the command from the flags in args, and the
// repository location from the single positional argument.
func (c *catCmd) Parse(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return cmd.Errorf(c, "%w", err)
	}
	if c.fs.NArg() > 1 {
		return cmd.Errorf(c, "%w", &cmd.ArgInvalidError{Name: "argument", Value: c.fs.Arg(1)})
	}
	location, err := cmd.RequiredArg(c, c.fs, "repository")
	if err != nil {
		return err
	}
	if c.tmpClone && c.localDir != "" {
		return cmd.Errorf(c, "-dir and -tmp cannot be used together")
	}
	if c.headLines < 0 {
		return cmd.Errorf(c, "%w", &cmd.ArgInvalidError{Name: "-head", Value: strconv.Itoa(c.headLines)})
	}
	if c.location, err = gitpath.Parse(location); err != nil {
		return cmd.Errorf(c, "%w", err)
	}
	return nil
}

// Run writes the files of the repository to the output file,
// or to the standard output of e when no output file is set.
// A remote repository is cloned first.
// With -tmp, into a temporary directory removed when Run returns.
func (c *catCmd) Run(e *env) error {
	setupLog(e.stderr, c.debug)
	if c.dryRun {
		c.logDryRun()
		return nil
	}
	if c.location.IsLocal() || !c.tmpClone {
		return c.cat(e, c.cloneDir())
	}
	tmpDir, err := os.MkdirTemp("", "gitcat-*")
	if err != nil {
		return cmd.Errorf(c, "create temporary directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()
	return c.cat(e, tmpDir)
}

// cloneDir returns the directory to clone a remote repository into:
// the -dir flag, or the repository name in the current directory.
func (c *catCmd) cloneDir() string {
	if c.localDir != "" {
		return c.localDir
	}
	return c.location.Name()
}

// cat opens, filters, reads, and writes the files of the repository.
// A remote repository is cloned into cloneDir first.
func (c *catCmd) cat(e *env, cloneDir string) error {
	r, err := c.open(e, cloneDir)
	if err != nil {
		return cmd.Errorf(c, "%w", err)
	}
	log.Info("listed repository files", "count", r.Len())

	r = c.fileFilter().Apply(r)
	log.Info("files after filters", "count", r.Len())

	loaded := files.Load(r, c.headLines)
	log.Info("files loaded", "count", len(loaded))

	if err = c.write(e.stdout, loaded); err != nil {
		return cmd.Errorf(c, "%w", err)
	}
	return nil
}

// open returns the repository, cloning a remote repository into cloneDir
// first.
func (c *catCmd) open(e *env, cloneDir string) (*repo.Repo, error) {
	if c.location.IsLocal() {
		log.Info("reading local repository", "path", c.location.Path)
		return source.Local(e.ctx, c.location.Path)
	}
	log.Info("cloning repository", "url", c.location.Path, "protocol", c.location.Kind, "dir", cloneDir)
	return source.Remote(e.ctx, c.location.Path, cloneDir)
}

// fileFilter returns the filter that the flags describe.
func (c *catCmd) fileFilter() *filter.Filter {
	f := filter.New().
		Only(c.includePaths...).
		Ignore(c.excludePaths...).
		Keep(c.keepExt...).
		SizeRange(int64(c.minSize), int64(c.maxSize))
	if !c.dotFiles {
		f.IgnoreDotFiles()
	}
	return f
}

// write renders the files to the output file named by -out with the format
// as its extension, or to stdout when the -out flag is not set.
func (c *catCmd) write(stdout io.Writer, loaded []files.File) (err error) {
	if c.outFile == "" {
		return output.Write(stdout, c.outFmt, loaded)
	}
	name := c.outFile + "." + string(c.outFmt)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	if err = output.Write(f, c.outFmt, loaded); err != nil {
		return err
	}
	log.Info("output written to file", "file", name)
	return nil
}

// logDryRun logs what Run would do, without doing it.
func (c *catCmd) logDryRun() {
	log.Info("dry run mode enabled: no actions will be executed")
	if c.location.IsLocal() {
		log.Info("would list local repository", "path", c.location.Path)
		return
	}
	if c.tmpClone {
		log.Info("would clone repository", "url", c.location.Path, "protocol", c.location.Kind, "dir", "<tmp>")
		log.Warn("cloning to a tmp directory: deleted after execution")
		return
	}
	log.Info("would clone repository", "url", c.location.Path, "protocol", c.location.Kind, "dir", c.cloneDir())
}
