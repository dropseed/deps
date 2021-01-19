package main

import (
	"github.com/dropseed/deps/internal/runner"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Locally upgrade deps in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.Local(); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
