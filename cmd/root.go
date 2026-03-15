package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "nlm",
	Short: "NotebookLM CLI - Unofficial Google NotebookLM client",
	Long: `nlm is an unofficial CLI client for Google NotebookLM.
Manage notebooks, add sources, chat with AI, generate artifacts, and more from the command line.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Config file path")
}
