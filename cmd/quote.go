package cmd

import (
	"fmt"

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
		p := tea.NewProgram(quote_input.Model{
			Target: "The quick brown fox jumps over the lazy dog.",
		})
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
		}
	},
}
