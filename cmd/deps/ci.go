package main

import (
	"github.com/dropseed/deps/internal/runner"
	"github.com/spf13/cobra"
)

var ciAuto bool
var ciTypes []string

var ciCMD = &cobra.Command{
	Use:   "ci",
	Short: "Run deps on the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.CI(ciAuto, ciTypes); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	ciCMD.Flags().BoolVarP(&ciAuto, "autoconfigure", "a", false, "automatically configure repo for deps usage (push access)")
	ciCMD.Flags().StringArrayVarP(&ciTypes, "type", "t", []string{}, "only run on specified dependency types")
	rootCmd.AddCommand(ciCMD)
}
