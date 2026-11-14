// Gitcat concatenates the files of a git repository into a single document
// in the text, JSON Lines, or Markdown format, ready to be loaded into the
// context of a language model.
//
// Usage:
//
//	gitcat [flags] <repo>
//	gitcat <command> [flags] [args]
//
// The repository is a local directory, an http(s) or ssh URL, or an scp-like
// address such as git@github.com:i-zaitsev/gitcat.git. Remote repositories are
// cloned with a depth of one. The files of a git work tree are listed by git,
// so files ignored by .gitignore are left out.
//
// Run "gitcat help" for the list of commands and flags.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/i-zaitsev/pkg/cmd"
)

// main runs gitcat and exits with its status code.
func main() {
	os.Exit(run(os.Args[1:]))
}

// run executes gitcat with args and returns the exit status.
// An interrupt cancels the context of the commands.
func run(args []string) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e := newEnv(ctx, os.Stdout, os.Stderr)

	reg, err := registry()
	if err != nil {
		e.Error("error: %v\n", err)
		return 1
	}

	if err = dispatch(e, reg, args); err != nil {
		e.Error("error: %v\n", err)
		return 1
	}

	return 0
}

// registry returns the commands of gitcat keyed by name.
func registry() (map[string]command, error) {
	reg := map[string]command{}
	if err := cmd.Register(reg, newCat()); err != nil {
		return nil, err
	}
	return reg, nil
}
