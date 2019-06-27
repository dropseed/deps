package main

import (
	"github.com/dropseed/deps/internal/test"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run component tests",
	Run: func(cmd *cobra.Command, args []string) {
		if err := test.Run(); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	devCmd.AddCommand(testCmd)
	// Set these variables directly in the test module
	testCmd.Flags().BoolVarP(&test.UpdateOutputData, "update-output-data", "u", false, "Update output data")
	testCmd.Flags().BoolVarP(&test.LooseOutputDataComparison, "loose-output-data-comparison", "l", false, "Loose output data comparison")
	testCmd.Flags().BoolVarP(&test.ExitEarly, "exit-early", "x", false, "Exit on first failure or error")
	testCmd.Flags().StringVar(&test.FilterName, "filter-name", "", "Filter test cases by name substring")
}
