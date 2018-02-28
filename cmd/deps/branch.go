package main

import (
	"context"
	"flag"

	"github.com/dependencies-io/deps/internal/git"
	"github.com/dependencies-io/deps/internal/hooks"
	"github.com/google/subcommands"
)

// BranchCmd starts a new git branch
type BranchCmd struct {
}

// Name of the command
func (*BranchCmd) Name() string { return "branch" }

// Synopsis of the command
func (*BranchCmd) Synopsis() string { return "Start a branch for pullrequest" }

// Usage details of the command
func (*BranchCmd) Usage() string {
	return `branch
	Start a branch for pullrequest
`
}

// SetFlags for the command
func (p *BranchCmd) SetFlags(f *flag.FlagSet) {
}

// Execute the command
func (p *BranchCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := hooks.Run("before_branch"); err != nil {
		return printErrAndExitFailure(err)
	}

	if _, err := git.BranchForJob(); err != nil {
		return printErrAndExitFailure(err)
	}

	if err := hooks.Run("after_branch"); err != nil {
		return printErrAndExitFailure(err)
	}

	// we'll assume the update is following the creation of the branch
	if err := hooks.Run("before_update"); err != nil {
		return printErrAndExitFailure(err)
	}

	return subcommands.ExitSuccess
}
