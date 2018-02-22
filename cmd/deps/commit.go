package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/dependencies-io/pullrequest/internal/git"
	"github.com/dependencies-io/pullrequest/internal/hooks"
	"github.com/google/subcommands"
)

// CommitCmd adds and commits paths to git
type CommitCmd struct {
	message string
}

// Name of the command
func (*CommitCmd) Name() string { return "commit" }

// Synopsis of the command
func (*CommitCmd) Synopsis() string { return "Add and commit paths" }

// Usage details of the command
func (*CommitCmd) Usage() string {
	return `addcommit [-m <message>] [<paths>...]
	Add and commit paths
`
}

// SetFlags for the command
func (c *CommitCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.message, "m", "", "message")
}

// Execute the command
func (c *CommitCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.message == "" {
		fmt.Println("commit message is required")
		return subcommands.ExitUsageError
	}

	if len(f.Args()) == 0 {
		fmt.Println("paths not given")
		return subcommands.ExitUsageError
	}

	// if we had dependencies JSON...
	// - automatically generate commit title
	// - before_lockfile_commit
	// - before_manifest_commit

	if err := hooks.Run("before_commit"); err != nil {
		return printErrAndExitFailure(err)
	}

	if err := git.AddCommit(c.message, f.Args()); err != nil {
		return printErrAndExitFailure(err)
	}

	if err := hooks.Run("after_commit"); err != nil {
		return printErrAndExitFailure(err)
	}

	return subcommands.ExitSuccess
}
