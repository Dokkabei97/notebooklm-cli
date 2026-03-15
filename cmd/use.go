package cmd

import (
	"fmt"

	"github.com/Dokkabei97/notebooklm-cli/internal/config"
	"github.com/Dokkabei97/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <notebook-id>",
	Short: "Set active notebook",
	Long:  "Set the default notebook to use for subsequent commands.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		notebookID := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg.ActiveNotebook = notebookID
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		output.PrintSuccess(fmt.Sprintf("Active notebook: %s", notebookID))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

// requireNotebook returns the notebook ID from args or active config.
func requireNotebook(args []string, flagNotebook string) (string, error) {
	if flagNotebook != "" {
		return flagNotebook, nil
	}
	if len(args) > 0 {
		return args[0], nil
	}
	if id := config.GetActiveNotebook(); id != "" {
		return id, nil
	}
	return "", fmt.Errorf("please specify a notebook. Use 'nlm use <id>' or the --notebook flag")
}
