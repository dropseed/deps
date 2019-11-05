package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dropseed/deps/internal/output"

	"github.com/dropseed/deps/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a deps config in the current",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.InferredConfigFromDir(".")
		if err != nil {
			printErrAndExitFailure(err)
		}
		inferred, err := cfg.DumpYAML()
		if err != nil {
			printErrAndExitFailure(err)
		}

		filename := config.DefaultFilenames[0]

		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			printErrAndExitFailure(fmt.Errorf("%s already exists!", filename))
		}

		fmt.Printf("Generating config...\n\n%s\n", inferred)

		err = ioutil.WriteFile(filename, []byte(inferred), 0644)
		if err != nil {
			printErrAndExitFailure(err)
		}
		output.Success("✔ saved as %s", filename)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
