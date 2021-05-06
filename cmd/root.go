package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.1"
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:     "stew-builder",
	Short:   "Package builder for Stew.",
	Version: "v0.0.1",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
