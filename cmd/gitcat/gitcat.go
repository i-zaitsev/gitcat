package main

import (
	"fmt"
	"os"

	"github.com/i-zaitsev/gitcat/pkg/ls"
)

func main() {
	cli := NewCLI()

	if err := cli.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	lg := cli.lg

	if cli.dryRun {
		lg.Info("dry run mode enabled - no actions will be executed")
		if cli.location.IsLocal() {
			lg.Info("would list local repository", "path", cli.location.Path)
		} else {
			lg.Info("would clone repository",
				"url", cli.location.Path,
				"protocol", cli.location.Kind,
				"dir", cli.localDir)
			if cli.tmpClone {
				lg.Warn("cloning to a tmp directory - deleted after execution")
			}
		}
		return
	}

	if cli.location.IsLocal() {
		lg.Info("listing local repository", "path", cli.location.Path)
	} else {
		lg.Info("cloning repository",
			"url", cli.location.Path,
			"protocol", cli.location.Kind,
			"dir", cli.localDir)
	}

	var (
		lsErr error
		repo  *ls.RepoContent
	)

	if cli.location.IsLocal() {
		repo, lsErr = ls.LocalRepo(cli.location.Path)
	} else {
		cloneDir := cli.localDir
		if cli.tmpClone {
			tmpDir, err := os.MkdirTemp("", "gitcat-*")
			if err != nil {
				lg.Error("failed to create tmp directory", "error", err)
				os.Exit(1)
			}
			defer os.RemoveAll(tmpDir)
			cloneDir = tmpDir
		}
		repo, lsErr = ls.RemoteRepo(cli.location, cloneDir)
	}

	if lsErr != nil {
		lg.Error("failed to list repo files", "error", lsErr)
		os.Exit(1)
	}

	lg.Info("successfully listed repo files", "count", len(repo.Content))

	for _, file := range repo.Content {
		fmt.Println(file)
	}
}
