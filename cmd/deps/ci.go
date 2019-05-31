package main

import (
	"github.com/dependencies-io/deps/internal/runner"
	"github.com/spf13/cobra"
)

var ciCMD = &cobra.Command{
	Use:   "ci",
	Short: "Run deps on the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.CI(); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ciCMD)
}
