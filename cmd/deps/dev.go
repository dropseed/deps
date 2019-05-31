package main

import (
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Component development commands",
}

func init() {
	rootCmd.AddCommand(devCmd)
}
