package main

import (
	"context"
	"strings"
	"testing"

	"github.com/i-zaitsev/pkg/assert"
)

// newTestEnv returns an env that writes to the returned buffers.
func newTestEnv() (e *env, stdout, stderr *strings.Builder) {
	stdout, stderr = &strings.Builder{}, &strings.Builder{}
	return newEnv(context.Background(), stdout, stderr), stdout, stderr
}

// runGitcat dispatches args to the gitcat commands and returns what they
// wrote to stdout and stderr.
func runGitcat(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	reg, err := registry()
	assert.NotErr(t, err)
	e, out, errOut := newTestEnv()
	err = dispatch(e, reg, args)
	return out.String(), errOut.String(), err
}

func TestRegistry(t *testing.T) {
	reg, err := registry()

	assert.NotErr(t, err)
	assert.HasKey(t, reg, defaultCmd)
}

func TestDispatch_Help(t *testing.T) {
	for _, arg := range []string{"help", "-h", "-help", "--help"} {
		t.Run(arg, func(t *testing.T) {
			stdout, stderr, err := runGitcat(t, arg)

			assert.NotErr(t, err)
			assert.Eq(t, stderr, "")
			assert.ContainsAll(t, stdout,
				"Usage:\n  gitcat [flags] <repo>",
				"Commands:\n",
				"Concatenate the files of a repository",
				"-fmt value",
				"Examples:",
			)
		})
	}
}

func TestDispatch_CommandHelp(t *testing.T) {
	stdout, _, err := runGitcat(t, "cat", "-h")

	assert.NotErr(t, err)
	assert.ContainsAll(t, stdout, "Usage: gitcat cat [flags] <repo>", "-maxsize value")
}

func TestDispatch_NoArgs_PrintsUsageAndFails(t *testing.T) {
	stdout, stderr, err := runGitcat(t)

	assert.ErrContains(t, err, "cat: repository is required")
	assert.Eq(t, stdout, "")
	assert.Contains(t, stderr, "Usage: gitcat cat [flags] <repo>")
}

func TestDispatch_DefaultCommand_IsCat(t *testing.T) {
	dir := sampleRepo(t)

	explicit, _, errExplicit := runGitcat(t, "cat", "-fmt", "md", dir)
	implicit, _, errImplicit := runGitcat(t, "-fmt", "md", dir)

	assert.NotErr(t, errExplicit)
	assert.NotErr(t, errImplicit)
	assert.Contains(t, implicit, "## main.go")
	assert.Eq(t, implicit, explicit)
}
