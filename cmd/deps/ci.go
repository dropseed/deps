package main

import (
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/runner"
	"github.com/spf13/cobra"
)

var ciManual bool
var ciTypes []string
var ciQuiet bool

var ciCMD = &cobra.Command{
	Use:   "ci",
	Short: "Update all dependencies of the current branch, as pull requests",
	Run: func(cmd *cobra.Command, args []string) {
		// CI will run verbose by default
		if !ciQuiet {
			output.Verbosity = 1
		}

		auto := !ciManual
		if err := runner.CI(auto, ciTypes); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	ciCMD.Flags().BoolVarP(&ciManual, "manual", "m", false, "do not automatically configure repo")
	ciCMD.Flags().BoolVarP(&ciQuiet, "quiet", "q", false, "disable verbose output")
	ciCMD.Flags().StringArrayVarP(&ciTypes, "type", "t", []string{}, "only run on specified dependency types")
	rootCmd.AddCommand(ciCMD)
}
