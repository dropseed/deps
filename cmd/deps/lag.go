package main

import (
	"os"

	"github.com/dropseed/deps/internal/lag"
	"github.com/spf13/cobra"
)

var lagCmd = &cobra.Command{
	Use:   "lag",
	Short: "See if installed dependencies are lagging",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if err := lag.Run(cwd); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lagCmd)
}
