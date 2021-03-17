package main

import (
	"os"

	"github.com/dropseed/deps/internal/install"
	"github.com/dropseed/deps/internal/lag"
	"github.com/dropseed/deps/internal/output"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies in this directory",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		lockfiles := install.FindLockfiles(cwd)

		if len(lockfiles) < 1 {
			return
		}

		output.Event("Running install for:")
		for _, lockfile := range lockfiles {
			output.Event("- %s", lockfile.RelPath())
		}

		lagManager, err := lag.NewLagManager()
		if err != nil {
			printErrAndExitFailure(err)
		}

		for _, lockfile := range lockfiles {
			if err := lockfile.Install(); err == nil {
				id := lag.IdentifierForFile(lockfile.Path)
				lagManager.SaveLockfileIdentifier(lockfile.Path, id)
			} else {
				printErrAndExitFailure(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
