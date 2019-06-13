package main

import (
	"github.com/dropseed/deps/internal/runner"
	"github.com/spf13/cobra"
)

var ciUpdateLimit int

var ciCMD = &cobra.Command{
	Use:   "ci",
	Short: "Run deps on the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.CI(ciUpdateLimit); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	ciCMD.Flags().IntVarP(&ciUpdateLimit, "limit", "l", -1, "limit the number of updates performed")
	rootCmd.AddCommand(ciCMD)
}
