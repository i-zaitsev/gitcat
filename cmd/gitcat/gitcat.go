package main

import (
	"fmt"
	"os"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/ls"
	"github.com/i-zaitsev/gitcat/pkg/output"
)

// writeOutput writes content to the specified file or stdout.
// If outFile is empty, writes to stdout. Otherwise, appends the format extension.
func writeOutput(content, outFile string, format output.Format) error {
	if outFile == "" {
		fmt.Println(content)
		return nil
	}

	filename := outFile + "." + string(format)
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Info("output written to file", "file", filename)
	return nil
}

func main() {
	cli := NewCLI()

	if err := cli.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if cli.dryRun {
		log.Info("dry run mode enabled - no actions will be executed")
		if cli.location.IsLocal() {
			log.Info("would list local repository", "path", cli.location.Path)
		} else {
			log.Info("would clone repository",
				"url", cli.location.Path,
				"protocol", cli.location.Kind,
				"dir", cli.localDir)
			if cli.tmpClone {
				log.Warn("cloning to a tmp directory - deleted after execution")
			}
		}
		return
	}

	if cli.location.IsLocal() {
		log.Info("listing local repository", "path", cli.location.Path)
	} else {
		log.Info("cloning repository",
			"url", cli.location.Path,
			"protocol", cli.location.Kind,
			"dir", cli.localDir)
	}

	var (
		lsErr error
		repo  *ls.RepoContent
	)

	list := ls.NewList().
		IgnoreDotFiles().
		WithPaths(cli.includePaths...).
		ExcludePaths(cli.excludePaths...)

	if cli.location.IsLocal() {
		repo, lsErr = list.LocalRepo(cli.location.Path)
	} else {
		cloneDir := cli.localDir
		if cli.tmpClone {
			tmpDir, err := os.MkdirTemp("", "gitcat-*")
			if err != nil {
				log.Error("failed to create tmp directory", "error", err)
				os.Exit(1)
			}
			defer os.RemoveAll(tmpDir)
			cloneDir = tmpDir
		}
		repo, lsErr = list.RemoteRepo(cli.location, cloneDir)
	}

	if lsErr != nil {
		log.Error("failed to list repo files", "error", lsErr)
		os.Exit(1)
	}

	log.Info("successfully listed repo files", "count", len(repo.Files))

	if len(cli.keepExt) > 0 {
		log.Warn("keeping only files with extensions", "extensions", cli.keepExt)
		repo = files.MatchExt(repo, cli.keepExt...)
	}

	if cli.minSize > 0 || cli.maxSize >= 0 {
		log.Info("applying size filters", "minsize", cli.minSize.InBytes(), "maxsize", cli.maxSize.InBytes())
		repo = files.FilterBySize(repo, cli.minSize.InBytes(), cli.maxSize.InBytes())
	}

	log.Info("files after all filters", "count", len(repo.Files))

	var content string
	switch cli.outFmt {
	case output.FormatJSONL:
		log.Info("writing output to FormatGrouped")
		jsonl, err := output.ToJSONL(repo)
		if err != nil {
			log.Error("failed to generate JSONL output", "error", err)
			os.Exit(1)
		}
		content = jsonl
	case output.FormatText:
		log.Info("writing output to text")
		content = output.ToText(repo)
	case output.FormatMarkdown:
		log.Info("writing output to markdown")
		content = output.ToMarkdown(repo)
	}

	if err := writeOutput(content, cli.outFile, cli.outFmt); err != nil {
		log.Error("failed to write output", "error", err)
		os.Exit(1)
	}
}
