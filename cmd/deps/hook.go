package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/dependencies-io/deps/internal/hooks"
	"github.com/google/subcommands"
)

// HookCmd starts a new git branch
type HookCmd struct {
}

// Name of the command
func (*HookCmd) Name() string { return "hook" }

// Synopsis of the command
func (*HookCmd) Synopsis() string { return "Run a hook by name" }

// Usage details of the command
func (*HookCmd) Usage() string {
	return `hook <hook_name>
	Run a hook by name
`
}

// SetFlags for the command
func (p *HookCmd) SetFlags(f *flag.FlagSet) {
}

// Execute the command
func (p *HookCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		fmt.Println("hook name not given")
		return subcommands.ExitUsageError
	}

	hookName := f.Arg(0)

	if err := hooks.Run(hookName); err != nil {
		return printErrAndExitFailure(err)
	}

	return subcommands.ExitSuccess
}
