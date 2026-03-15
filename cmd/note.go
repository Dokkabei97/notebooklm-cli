package cmd

import (
	"fmt"

	"github.com/Dokkabei97/notebooklm-cli/internal/api"
	"github.com/Dokkabei97/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var noteNotebook string

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
}

var noteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, noteNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		notes, err := client.ListNotes(nbID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(notes)
		}

		var rows [][]string
		for _, n := range notes {
			rows = append(rows, []string{
				output.FormatID(n.ID),
				output.Truncate(n.Title, 40),
				output.Truncate(n.Content, 50),
				output.FormatTime(n.UpdatedAt),
			})
		}
		output.PrintTable([]string{"ID", "Title", "Content", "Updated"}, rows)
		return nil
	},
}

var noteCreateCmd = &cobra.Command{
	Use:   "create <title> <content>",
	Short: "Create a new note",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, noteNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		note, err := client.CreateNote(nbID, args[0], args[1])
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(note)
		}

		output.PrintSuccess(fmt.Sprintf("Note created: %s (%s)", note.Title, note.ID))
		return nil
	},
}

var noteUpdateCmd = &cobra.Command{
	Use:   "update <note-id> <title> <content>",
	Short: "Update a note",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, noteNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		note, err := client.UpdateNote(nbID, args[0], args[1], args[2])
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(note)
		}

		output.PrintSuccess(fmt.Sprintf("Note updated: %s", note.Title))
		return nil
	},
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete <note-id>",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, noteNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.DeleteNote(nbID, args[0]); err != nil {
			return err
		}

		output.PrintSuccess("Note deleted.")
		return nil
	},
}

func init() {
	noteCmd.PersistentFlags().StringVarP(&noteNotebook, "notebook", "n", "", "Notebook ID")

	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(noteUpdateCmd)
	noteCmd.AddCommand(noteDeleteCmd)
	rootCmd.AddCommand(noteCmd)
}
