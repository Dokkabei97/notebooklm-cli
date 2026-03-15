package cmd

import (
	"fmt"

	"github.com/jmk/notebooklm-cli/internal/api"
	"github.com/jmk/notebooklm-cli/internal/model"
	"github.com/jmk/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var shareNotebook string

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Manage sharing",
}

var shareSetCmd = &cobra.Command{
	Use:   "set <none|viewer|editor>",
	Short: "Set sharing permissions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, shareNotebook)
		if err != nil {
			return err
		}

		var perm model.SharePermission
		switch args[0] {
		case "none":
			perm = model.SharePermissionNone
		case "viewer":
			perm = model.SharePermissionViewer
		case "editor":
			perm = model.SharePermissionEditor
		default:
			return fmt.Errorf("invalid permission: %s (choose from none, viewer, editor)", args[0])
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		if err := client.SetSharing(nbID, perm); err != nil {
			return err
		}

		output.PrintSuccess(fmt.Sprintf("Sharing permission set: %s", args[0]))
		return nil
	},
}

var shareStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check current sharing status",
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, shareNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		status, err := client.GetSharing(nbID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.PrintJSON(status)
		}

		shared := "No"
		if status.IsShared {
			shared = "Yes"
		}
		output.PrintKeyValue([][2]string{
			{"Shared", shared},
			{"Permission", status.Permission.String()},
			{"Share URL", status.ShareURL},
		})
		return nil
	},
}

func init() {
	shareCmd.PersistentFlags().StringVarP(&shareNotebook, "notebook", "n", "", "Notebook ID")

	shareCmd.AddCommand(shareSetCmd)
	shareCmd.AddCommand(shareStatusCmd)
	rootCmd.AddCommand(shareCmd)
}
