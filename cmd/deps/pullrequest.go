package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dependencies-io/deps/internal/git"
	"github.com/dependencies-io/deps/internal/hooks"
	"github.com/dependencies-io/deps/internal/pullrequest/adapter"
	"github.com/google/subcommands"
)

// PullrequestCmd creates a pull request
type PullrequestCmd struct {
	branch               string
	title                string
	body                 string
	dependenciesJSON     string
	relatedPRTitleSearch string
}

// Name of the command
func (*PullrequestCmd) Name() string { return "pullrequest" }

// Synopsis of the command
func (*PullrequestCmd) Synopsis() string { return "Create a pull request" }

// Usage details of the command
func (*PullrequestCmd) Usage() string {
	return `pullrequest <dependencies_json_path>
	Create a pull request using a dependencies JSON file
`
}

// SetFlags for the command
func (p *PullrequestCmd) SetFlags(f *flag.FlagSet) {
}

// Execute the command
func (p *PullrequestCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	dependenciesJSONPath := f.Arg(0)
	if dependenciesJSONPath == "" {
		fmt.Println("path is required as arg")
		return subcommands.ExitUsageError
	}

	// get this first to make sure there are no problems
	pr, err := adapter.PullrequestAdapterFromDependenciesJSONPathAndHost(dependenciesJSONPath, os.Getenv("GIT_HOST"))
	if err != nil {
		return printErrAndExitFailure(err)
	}

	// assume this is the first thing to happen after updates are made
	if err = hooks.Run("after_update"); err != nil {
		return printErrAndExitFailure(err)
	}
	if err = hooks.Run("before_pullrequest"); err != nil {
		return printErrAndExitFailure(err)
	}

	if err = git.PushJobBranch(); err != nil {
		return printErrAndExitFailure(err)
	}

	if err = pr.Create(); err != nil {
		return printErrAndExitFailure(err)
	}
	if err = pr.DoRelated(); err != nil {
		return printErrAndExitFailure(err)
	}
	if err = pr.OutputActions(); err != nil {
		return printErrAndExitFailure(err)
	}

	if err = hooks.Run("after_pullrequest"); err != nil {
		return printErrAndExitFailure(err)
	}

	return subcommands.ExitSuccess
}
