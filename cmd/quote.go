package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/cmd/ui/quote_input"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(quoteCmd)
}

var quoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "Get a random quote",
	Long:  `Get a random quote`,
	Run: func(cmd *cobra.Command, args []string) {
		tprogram := tea.NewProgram(quote_input.InitialQuoteModel())
		if _, err := tprogram.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
