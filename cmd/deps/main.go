package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

func printErrAndExitFailure(err error) subcommands.ExitStatus {
	fmt.Println(err)
	return subcommands.ExitFailure
}

func main() {
	// bulitins
	subcommands.Register(subcommands.HelpCommand(), "")
	// subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	// our commands
	subcommands.Register(&BranchCmd{}, "")
	subcommands.Register(&CommitCmd{}, "")
	subcommands.Register(&PullrequestCmd{}, "")
	subcommands.Register(&CollectCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
