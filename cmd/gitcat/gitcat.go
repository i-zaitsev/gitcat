package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/i-zaitsev/gitcat/pkg/files"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/gitcat/pkg/ls"
)

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

	list := ls.NewList().IgnoreDotFiles()

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

	allFileExt := files.DiscoverExt(repo)
	log.Info("found file extensions", "found", strings.Join(allFileExt, ", "))

	log.Info("taking only .go files")
	onlyGo := files.MatchExt(repo, ".go")

	for _, line := range files.Cat(onlyGo.Files) {
		fmt.Printf("%s\n", line)
	}
}
