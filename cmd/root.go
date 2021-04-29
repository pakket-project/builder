package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.1"
)

var rootCmd = &cobra.Command{
	Use:   "stew-builder",
	Short: "Package builder for Stew.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
