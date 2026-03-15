package cmd

import (
	"fmt"
	"os"

	"github.com/jmk/notebooklm-cli/internal/api"
	"github.com/jmk/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var artifactNotebook string

var artifactCmd = &cobra.Command{
	Use:     "artifact",
	Aliases: []string{"art"},
	Short:   "Manage artifacts",
}

var artifactListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, artifactNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		artifacts, err := client.ListArtifacts(nbID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(artifacts)
		}

		var rows [][]string
		for _, art := range artifacts {
			rows = append(rows, []string{
				output.FormatID(art.ID),
				output.Truncate(art.Title, 40),
				art.Type.String(),
				art.Status,
				output.FormatTime(art.CreatedAt),
			})
		}
		output.PrintTable([]string{"ID", "Title", "Type", "Status", "Created"}, rows)
		return nil
	},
}

var artifactDeleteCmd = &cobra.Command{
	Use:   "delete <artifact-id>",
	Short: "Delete an artifact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, artifactNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.DeleteArtifact(nbID, args[0]); err != nil {
			return err
		}

		output.PrintSuccess("Artifact deleted.")
		return nil
	},
}

var artifactExportCmd = &cobra.Command{
	Use:   "export <artifact-id> [output-file]",
	Short: "Export artifact content",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, artifactNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		art, err := client.GetArtifact(nbID, args[0])
		if err != nil {
			return err
		}

		content := art.Content
		if content == "" {
			return fmt.Errorf("artifact has no content to export")
		}

		if len(args) > 1 {
			if err := os.WriteFile(args[1], []byte(content), 0644); err != nil {
				return err
			}
			output.PrintSuccess(fmt.Sprintf("Saved: %s", args[1]))
		} else {
			fmt.Println(content)
		}

		return nil
	},
}

func init() {
	artifactCmd.PersistentFlags().StringVarP(&artifactNotebook, "notebook", "n", "", "Notebook ID")

	artifactCmd.AddCommand(artifactListCmd)
	artifactCmd.AddCommand(artifactDeleteCmd)
	artifactCmd.AddCommand(artifactExportCmd)
	rootCmd.AddCommand(artifactCmd)
}
