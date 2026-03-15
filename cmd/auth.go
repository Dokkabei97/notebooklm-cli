package cmd

import (
	"fmt"
	"net/http"

	"github.com/jmk/notebooklm-cli/internal/auth"
	"github.com/jmk/notebooklm-cli/internal/output"
	"github.com/spf13/cobra"
)

var reuseChrome bool

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with Google account via browser",
	Long: `Set up Google account authentication.

Default:  Open a new Chrome window for Google login
--reuse:  Extract cookies directly from Chrome without opening it (preserves login state)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cookies []*http.Cookie
		var err error

		if reuseChrome {
			output.PrintInfo("Extracting cookies directly from Chrome cookie DB (no Chrome launch needed)...")
			cookies, err = auth.ExtractChromeCookies()
		} else {
			cookies, err = auth.BrowserLogin()
		}
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if err := auth.SaveStorageState(cookies); err != nil {
			return fmt.Errorf("failed to save cookies: %w", err)
		}

		tokens, err := auth.ExtractTokens(cookies)
		if err != nil {
			return fmt.Errorf("failed to extract tokens: %w", err)
		}

		if tokens.IsValid() {
			output.PrintSuccess(fmt.Sprintf("Authentication successful! (%d cookies saved)", len(cookies)))
		} else {
			output.PrintError("Cookies saved but token extraction failed. Please try again.")
		}
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cookies, err := auth.LoadStorageState()
		if err != nil {
			output.PrintError("No saved credentials found. Log in with 'nlm auth login'.")
			return nil
		}

		tokens, err := auth.ExtractTokens(cookies)
		if err != nil {
			output.PrintError(fmt.Sprintf("Authentication expired: %v", err))
			output.PrintInfo("Re-authenticate with 'nlm auth login --reuse'.")
			return nil
		}

		if tokens.IsValid() {
			output.PrintSuccess("Authentication valid")
			output.PrintKeyValue([][2]string{
				{"Cookies", fmt.Sprintf("%d", len(tokens.Cookies))},
				{"CSRF Token", tokens.CSRFToken[:min(20, len(tokens.CSRFToken))] + "..."},
				{"Session ID", tokens.SessionID},
			})
		} else {
			output.PrintError("Authentication credentials are invalid.")
		}
		return nil
	},
}

var authClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete saved credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.ClearStorageState(); err != nil {
			return fmt.Errorf("failed to delete credentials: %w", err)
		}
		output.PrintSuccess("Credentials deleted.")
		return nil
	},
}

func init() {
	authLoginCmd.Flags().BoolVar(&reuseChrome, "reuse", false, "Extract directly from Chrome cookie DB (no Chrome shutdown needed)")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authClearCmd)
	rootCmd.AddCommand(authCmd)
}
