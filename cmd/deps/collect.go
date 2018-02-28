package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/dependencies-io/deps/internal/collect"
	"github.com/google/subcommands"
)

// CollectCmd collects a dependencies JSON file at a path
type CollectCmd struct {
}

// Name of the command
func (*CollectCmd) Name() string { return "collect" }

// Synopsis of the command
func (*CollectCmd) Synopsis() string { return "Collect a dependencies JSON file" }

// Usage details of the command
func (*CollectCmd) Usage() string {
	return `collect <dependencies JSON file>
	Collect a dependencies JSON file
`
}

// SetFlags for the command
func (p *CollectCmd) SetFlags(f *flag.FlagSet) {
}

// Execute the command
func (p *CollectCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)

	if path == "" {
		fmt.Println("path must be given")
		return subcommands.ExitUsageError
	}

	if err := collect.DependenciesJSONPath(path); err != nil {
		return printErrAndExitFailure(err)
	}

	return subcommands.ExitSuccess
}
