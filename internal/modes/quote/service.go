package quote

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/quote_input"
)

func Run(cfg models.Config) error {
	p := tea.NewProgram(quote_input.InitialModel("The quick brown fox jumps over the lazy dog."))

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
