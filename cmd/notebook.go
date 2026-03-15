package cmd

import (
	"fmt"

	"github.com/Dokkabei97/notebooklm-cli/internal/api"
	"github.com/Dokkabei97/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var notebookCmd = &cobra.Command{
	Use:     "notebook",
	Aliases: []string{"nb"},
	Short:   "Manage notebooks",
}

var notebookListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List notebooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		notebooks, err := client.ListNotebooks()
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(notebooks)
		}

		var rows [][]string
		for _, nb := range notebooks {
			rows = append(rows, []string{
				output.FormatID(nb.ID),
				output.Truncate(nb.Title, 40),
				fmt.Sprintf("%d", nb.SourceCount),
				output.FormatTime(nb.UpdatedAt),
			})
		}
		output.PrintTable([]string{"ID", "Title", "Sources", "Updated"}, rows)
		return nil
	},
}

var notebookCreateCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new notebook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		nb, err := client.CreateNotebook(args[0])
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(nb)
		}

		output.PrintSuccess(fmt.Sprintf("Notebook created: %s (%s)", nb.Title, nb.ID))
		return nil
	},
}

var notebookDeleteCmd = &cobra.Command{
	Use:   "delete <notebook-id>",
	Short: "Delete a notebook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.DeleteNotebook(args[0]); err != nil {
			return err
		}

		output.PrintSuccess("Notebook deleted.")
		return nil
	},
}

var notebookRenameCmd = &cobra.Command{
	Use:   "rename <notebook-id> <new-title>",
	Short: "Rename a notebook",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.RenameNotebook(args[0], args[1]); err != nil {
			return err
		}

		output.PrintSuccess(fmt.Sprintf("Notebook renamed: %s", args[1]))
		return nil
	},
}

func init() {
	notebookCmd.AddCommand(notebookListCmd)
	notebookCmd.AddCommand(notebookCreateCmd)
	notebookCmd.AddCommand(notebookDeleteCmd)
	notebookCmd.AddCommand(notebookRenameCmd)
	rootCmd.AddCommand(notebookCmd)
}
