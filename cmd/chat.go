package cmd

import (
	"fmt"
	"strings"

	"github.com/jmk/notebooklm-cli/internal/api"
	"github.com/jmk/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	chatNotebook  string
	chatSourceIDs []string
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "AI chat",
}

var chatAskCmd = &cobra.Command{
	Use:   "ask <question>",
	Short: "Ask AI a question",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(args[0]) == "" {
			return fmt.Errorf("please enter a question")
		}

		nbID, err := requireNotebook(nil, chatNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		result, err := client.Ask(nbID, args[0], chatSourceIDs)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(result)
		}

		// Print answer
		fmt.Println(result.Answer)

		// Print citations
		if len(result.Sources) > 0 {
			fmt.Println()
			fmt.Println("--- Sources ---")
			for i, src := range result.Sources {
				fmt.Printf("[%d] %s: %s\n", i+1, src.SourceName, output.Truncate(src.Text, 80))
			}
		}

		// Print follow-ups
		if len(result.FollowUps) > 0 {
			fmt.Println()
			fmt.Println("--- Related Questions ---")
			for _, fu := range result.FollowUps {
				fmt.Printf("  • %s\n", fu)
			}
		}

		return nil
	},
}

var chatHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View chat history",
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, chatNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		entries, err := client.GetChatHistory(nbID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(entries)
		}

		for _, entry := range entries {
			role := strings.ToUpper(entry.Role)
			fmt.Printf("[%s] %s\n\n", role, entry.Content)
		}

		return nil
	},
}

func init() {
	chatCmd.PersistentFlags().StringVarP(&chatNotebook, "notebook", "n", "", "Notebook ID")
	chatAskCmd.Flags().StringSliceVarP(&chatSourceIDs, "sources", "s", nil, "Source IDs to reference")

	chatCmd.AddCommand(chatAskCmd)
	chatCmd.AddCommand(chatHistoryCmd)
	rootCmd.AddCommand(chatCmd)
}
