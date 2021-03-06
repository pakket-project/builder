package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.2"
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:     "pakket-builder",
	Short:   "Package builder for Pakket.",
	Version: Version,
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
