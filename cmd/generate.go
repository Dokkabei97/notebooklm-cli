package cmd

import (
	"fmt"
	"time"

	"github.com/Dokkabei97/notebooklm-cli/internal/api"
	"github.com/Dokkabei97/notebooklm-cli/internal/output"
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
	"github.com/spf13/cobra"
)

var (
	generateNotebook     string
	generateWait         bool
	generateInstructions string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate content (audio, study guide, FAQ, etc.)",
}

var generateAudioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Generate audio overview",
	RunE: func(cmd *cobra.Command, args []string) error {
		nbID, err := requireNotebook(nil, generateNotebook)
		if err != nil {
			return err
		}

		client, err := api.Authenticate()
		if err != nil {
			return err
		}

		art, err := client.GenerateAudio(nbID, generateInstructions)
		if err != nil {
			return err
		}

		if generateWait && art != nil {
			output.PrintInfo("Generating audio... waiting for completion.")
			art, err = client.WaitForArtifact(nbID, art.ID, 10*time.Minute)
			if err != nil {
				return err
			}
		}

		if jsonOutput {
			return output.PrintJSON(art)
		}

		output.PrintSuccess(fmt.Sprintf("Audio generation started: %s [%s]", art.ID, art.Status))
		if art.AudioURL != "" {
			fmt.Printf("URL: %s\n", art.AudioURL)
		}
		return nil
	},
}

func makeGenerateSubcmd(name, desc string, typeCode rpc.ArtifactTypeCode) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Generate " + desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			nbID, err := requireNotebook(nil, generateNotebook)
			if err != nil {
				return err
			}

			client, err := api.Authenticate()
			if err != nil {
				return err
			}

			art, err := client.CreateArtifact(nbID, typeCode)
			if err != nil {
				return err
			}

			if generateWait && art != nil {
				output.PrintInfo("Generating " + desc + "...")
				art, err = client.WaitForArtifact(nbID, art.ID, 5*time.Minute)
				if err != nil {
					return err
				}
			}

			if jsonOutput {
				return output.PrintJSON(art)
			}

			output.PrintSuccess(fmt.Sprintf("%s generated: %s [%s]", desc, art.ID, art.Status))
			return nil
		},
	}
}

func init() {
	generateCmd.PersistentFlags().StringVarP(&generateNotebook, "notebook", "n", "", "Notebook ID")
	generateCmd.PersistentFlags().BoolVarP(&generateWait, "wait", "w", false, "Wait for generation to complete")
	generateAudioCmd.Flags().StringVarP(&generateInstructions, "instructions", "i", "", "Audio generation instructions")

	generateCmd.AddCommand(generateAudioCmd)
	generateCmd.AddCommand(makeGenerateSubcmd("report", "report", rpc.ArtifactCodeReport))
	generateCmd.AddCommand(makeGenerateSubcmd("quiz", "quiz", rpc.ArtifactCodeQuiz))
	generateCmd.AddCommand(makeGenerateSubcmd("mind-map", "mind map", rpc.ArtifactCodeMindMap))
	generateCmd.AddCommand(makeGenerateSubcmd("video", "video", rpc.ArtifactCodeVideo))
	generateCmd.AddCommand(makeGenerateSubcmd("infographic", "infographic", rpc.ArtifactCodeInfographic))
	generateCmd.AddCommand(makeGenerateSubcmd("slide-deck", "slide deck", rpc.ArtifactCodeSlideDeck))
	rootCmd.AddCommand(generateCmd)
}
