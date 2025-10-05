package cmd

import (
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/spf13/cobra"
)

var languagesCmd = &cobra.Command{
	Use:     "languages",
	Aliases: []string{"list-languages", "langs"},
	Short:   "List supported languages",
	Long:    `List all the languages supported by the typing test application.`,
	Example: "typing-test-tui languages",
	Args:    cobra.NoArgs,
	Run:     listLanguages,
}

func listLanguages(cmd *cobra.Command, args []string) {
	cmd.Println("Supported Languages:")
	for _, lang := range models.SupportedLanguages() {
		cmd.Println(" -", lang)
	}
	cmd.Println("\nYou can specify a language using the --language or -l flag when starting a typing test.")
}

func init() {
	rootCmd.AddCommand(languagesCmd)
}
