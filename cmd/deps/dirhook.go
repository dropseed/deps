package main

import (
	"github.com/dropseed/deps/internal/lag"
	"github.com/spf13/cobra"
)

var dirhookCmd = &cobra.Command{
	Use:    "dirhook",
	Hidden: true,
	Short:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if err := lag.Run(); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dirhookCmd)
}
