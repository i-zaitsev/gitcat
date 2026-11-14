package main

import (
	"fmt"

	"github.com/i-zaitsev/gitcat/internal/version"
	"github.com/i-zaitsev/pkg/cmd"
)

// versionCmd prints the version of gitcat.
type versionCmd struct{}

// Name returns the name of the command.
func (c *versionCmd) Name() string { return "version" }

// Desc returns a one-line description of the command.
func (c *versionCmd) Desc() string { return "Print the version of gitcat" }

// Args returns the positional arguments of the command, which are none.
func (c *versionCmd) Args() string { return "" }

// Parse fails when given any arguments, since the command takes none.
func (c *versionCmd) Parse(args []string) error {
	if len(args) > 0 {
		return cmd.Errorf(c, "%w", &cmd.ArgInvalidError{Name: "argument", Value: args[0]})
	}
	return nil
}

// Run writes the version to the standard output of e.
func (c *versionCmd) Run(e *env) error {
	_, err := fmt.Fprintln(e.stdout, version.Get())
	return err
}
