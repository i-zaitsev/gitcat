package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"slices"

	"github.com/i-zaitsev/gitcat/internal/version"
	"github.com/i-zaitsev/gitcat/pkg/log"
	"github.com/i-zaitsev/pkg/cmd"
)

// prog is the program name shown in usage and log messages.
const prog = "gitcat"

// defaultCmd is the command that runs when the first argument is not the
// name of a command.
const defaultCmd = "cat"

// usageExamples ends the usage text of the program.
const usageExamples = `
Examples:
  gitcat git@github.com:i-zaitsev/gitcat.git
  gitcat -tmp -fmt md https://github.com/i-zaitsev/gitcat
  gitcat -dryrun -dir ./gitcat-src git@github.com:i-zaitsev/gitcat.git
  gitcat -tmp -path pkg/files,cmd -exclude testdata https://github.com/i-zaitsev/gitcat
  gitcat -tmp -maxsize 500 -keep .go https://github.com/i-zaitsev/gitcat
  gitcat -tmp -head 50 -out gitcat https://github.com/i-zaitsev/gitcat
`

// helpArgs are the first arguments that print the usage of the program.
var helpArgs = map[string]bool{
	"help":   true,
	"-h":     true,
	"-help":  true,
	"--help": true,
}

// command is a gitcat subcommand.
type command = cmd.Command[*env]

// flagPrinter is implemented by commands that can describe their flags.
type flagPrinter interface {
	PrintFlags(w io.Writer)
}

// env is the environment that commands run in: the output streams and a
// context that is canceled on interrupt.
type env struct {
	ctx    context.Context
	stdout io.Writer
	stderr io.Writer
}

// newEnv returns an env bound to ctx that writes to stdout and stderr.
func newEnv(ctx context.Context, stdout, stderr io.Writer) *env {
	return &env{ctx: ctx, stdout: stdout, stderr: stderr}
}

// Error writes a formatted message to the standard error of the env.
func (e *env) Error(format string, args ...any) {
	_, _ = fmt.Fprintf(e.stderr, format, args...)
}

// dispatch runs the command named by the first argument with the remaining
// arguments. When the first argument is not a command name, all arguments
// go to the default command, so that
//
//	gitcat -fmt md <repo>
//
// is the same as
//
//	gitcat cat -fmt md <repo>
//
// A help argument prints the usage of the program, and -h after a command
// prints the usage of the command. When a required argument is missing, the
// usage of the command is printed to stderr along with the error.
func dispatch(e *env, reg map[string]command, args []string) error {
	if len(args) > 0 && helpArgs[args[0]] {
		printUsage(e.stdout, reg)
		return nil
	}
	c := reg[defaultCmd]
	if len(args) > 0 {
		if named, ok := reg[args[0]]; ok {
			c, args = named, args[1:]
		}
	}
	err := c.Parse(args)
	switch {
	case errors.Is(err, flag.ErrHelp):
		printCmdUsage(e.stdout, c)
		return nil
	case errors.As(err, new(*cmd.ArgMissingError)):
		printCmdUsage(e.stderr, c)
		return err
	case err != nil:
		return err
	}
	return c.Run(e)
}

// printUsage writes the usage of the program to w: the synopsis, the
// commands, the flags of the default command, and examples.
func printUsage(w io.Writer, reg map[string]command) {
	_, _ = fmt.Fprintf(w, "%s %s - concatenates a git repository into a single file\n\n", prog, version.Get())
	_, _ = fmt.Fprintf(w, "Usage:\n  %s [flags] <repo>\n  %s <command> [flags] [args]\n\n", prog, prog)
	_, _ = fmt.Fprintf(w, "Commands:\n%s\n", cmd.Usage(slices.Collect(maps.Values(reg))...))
	_, _ = fmt.Fprintf(w, "Flags of %s, the default command:\n", defaultCmd)
	printFlags(w, reg[defaultCmd])
	_, _ = fmt.Fprint(w, usageExamples)
}

// printCmdUsage writes the usage of the command c and its flags to w.
func printCmdUsage(w io.Writer, c command) {
	_, _ = fmt.Fprintf(w, "Usage: %s %s [flags] %s\n\n%s\n\nFlags:\n", prog, c.Name(), c.Args(), c.Desc())
	printFlags(w, c)
}

// printFlags writes the flags of the command c to w, if c can describe them.
func printFlags(w io.Writer, c command) {
	if p, ok := c.(flagPrinter); ok {
		p.PrintFlags(w)
	}
}

// setupLog sends the logs of gitcat to w, colored by level. Debug messages
// are logged only when debug is set.
func setupLog(w io.Writer, debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(log.NewColorWriter(w), &slog.HandlerOptions{Level: level})
	log.SetLogger(slog.New(handler).With("app", prog))
}
