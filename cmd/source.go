package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmk/notebooklm-cli/internal/api"
	"github.com/jmk/notebooklm-cli/internal/model"
	"github.com/jmk/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	sourceNotebook string
	sourceWait     bool
)

var sourceCmd = &cobra.Command{
	Use:     "source",
	Aliases: []string{"src"},
	Short:   "Manage sources",
}

var sourceListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, sourceNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		sources, err := client.ListSources(nbID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(sources)
		}

		var rows [][]string
		for _, src := range sources {
			rows = append(rows, []string{
				output.FormatID(src.ID),
				output.Truncate(src.Title, 40),
				src.Type.String(),
				src.Status,
				output.FormatTime(src.CreatedAt),
			})
		}
		output.PrintTable([]string{"ID", "Title", "Type", "Status", "Created"}, rows)
		return nil
	},
}

var sourceAddCmd = &cobra.Command{
	Use:   "add <url-or-file>",
	Short: "Add a source (URL or file)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, sourceNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		input := strings.TrimSpace(args[0])
		if input == "" {
			return fmt.Errorf("please provide a URL or file path")
		}

		var src *model.Source

		if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
			src, err = client.AddSourceURL(nbID, input)
		} else {
			src, err = client.AddSourceFile(nbID, input)
		}

		if err != nil {
			return err
		}

		if sourceWait && src != nil {
			output.PrintInfo("Processing source... waiting for completion.")
			src, err = client.WaitForSource(nbID, src.ID, 5*time.Minute)
			if err != nil {
				return err
			}
		}

		if jsonOutput {
			return output.PrintJSON(src)
		}

		output.PrintSuccess(fmt.Sprintf("Source added: %s (%s) [%s]", src.Title, src.ID, src.Status))
		return nil
	},
}

var sourceGetCmd = &cobra.Command{
	Use:   "get <source-id>",
	Short: "Get source details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, sourceNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		src, err := client.GetSource(nbID, args[0])
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(src)
		}

		output.PrintKeyValue([][2]string{
			{"ID", src.ID},
			{"Title", src.Title},
			{"Type", src.Type.String()},
			{"Status", src.Status},
			{"URL", src.URL},
			{"Created", output.FormatTime(src.CreatedAt)},
		})
		return nil
	},
}

var sourceDeleteCmd = &cobra.Command{
	Use:   "delete <source-id>",
	Short: "Delete a source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, sourceNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.DeleteSource(nbID, args[0]); err != nil {
			return err
		}

		output.PrintSuccess("Source deleted.")
		return nil
	},
}

var sourceRefreshCmd = &cobra.Command{
	Use:   "refresh <source-id>",
	Short: "Refresh a source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, sourceNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.RefreshSource(nbID, args[0]); err != nil {
			return err
		}

		output.PrintSuccess("Source refresh started.")
		return nil
	},
}

var sourceWaitCmd = &cobra.Command{
	Use:   "wait <source-id>",
	Short: "Wait for source processing to complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, sourceNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		output.PrintInfo("Waiting for source processing...")
		src, err := client.WaitForSource(nbID, args[0], 5*time.Minute)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(src)
		}

		output.PrintSuccess(fmt.Sprintf("Source ready: %s [%s]", src.Title, src.Status))
		return nil
	},
}

func init() {
	sourceCmd.PersistentFlags().StringVarP(&sourceNotebook, "notebook", "n", "", "Notebook ID")
	sourceAddCmd.Flags().BoolVarP(&sourceWait, "wait", "w", false, "Wait for source processing to complete")

	sourceCmd.AddCommand(sourceListCmd)
	sourceCmd.AddCommand(sourceAddCmd)
	sourceCmd.AddCommand(sourceGetCmd)
	sourceCmd.AddCommand(sourceDeleteCmd)
	sourceCmd.AddCommand(sourceRefreshCmd)
	sourceCmd.AddCommand(sourceWaitCmd)
	rootCmd.AddCommand(sourceCmd)
}
