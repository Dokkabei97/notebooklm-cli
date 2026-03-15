package cmd

import (
	"fmt"
	"time"

	"github.com/jmk/notebooklm-cli/internal/api"
	"github.com/jmk/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var researchNotebook string

var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Deep research",
}

var researchStartCmd = &cobra.Command{
	Use:   "start <query>",
	Short: "Start deep research",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, researchNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		result, err := client.StartResearch(nbID, args[0])
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(result)
		}

		output.PrintSuccess(fmt.Sprintf("Research started: %s", result.ID))
		output.PrintInfo("Check progress with 'nlm research poll <id>'.")
		return nil
	},
}

var researchPollCmd = &cobra.Command{
	Use:   "poll <research-id>",
	Short: "Check research progress",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, researchNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		result, err := client.PollResearch(nbID, args[0])
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(result)
		}

		output.PrintKeyValue([][2]string{
			{"ID", result.ID},
			{"Status", result.Status},
			{"Progress", fmt.Sprintf("%d%%", result.Progress)},
		})

		if result.Content != "" {
			fmt.Println()
			fmt.Println(result.Content)
		}
		return nil
	},
}

var researchImportCmd = &cobra.Command{
	Use:   "import <research-id>",
	Short: "Import research results as a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, researchNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		// Poll until complete
		var result *api.ResearchResult
		deadline := time.Now().Add(10 * time.Minute)
		for time.Now().Before(deadline) {
			result, err = client.PollResearch(nbID, args[0])
			if err != nil {
				return err
			}
			if result.Status == "completed" {
				break
			}
			if result.Status == "error" {
				return fmt.Errorf("research failed")
			}
			time.Sleep(3 * time.Second)
		}

		if result == nil || result.Content == "" {
			return fmt.Errorf("no research results available")
		}

		note, err := client.CreateNote(nbID, "Research: "+args[0], result.Content)
		if err != nil {
			return err
		}

		output.PrintSuccess(fmt.Sprintf("Research results saved as note: %s", note.ID))
		return nil
	},
}

func init() {
	researchCmd.PersistentFlags().StringVarP(&researchNotebook, "notebook", "n", "", "Notebook ID")

	researchCmd.AddCommand(researchStartCmd)
	researchCmd.AddCommand(researchPollCmd)
	researchCmd.AddCommand(researchImportCmd)
	rootCmd.AddCommand(researchCmd)
}
