package main

import (
	"os"

	"github.com/dropseed/deps/internal/output"
)

func printErrAndExitFailure(err error) {
	output.Error(err.Error())
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		printErrAndExitFailure(err)
	}
}
