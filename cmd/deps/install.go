package main

import (
	"github.com/dropseed/deps/internal/install"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies in this directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := install.Install(); err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
