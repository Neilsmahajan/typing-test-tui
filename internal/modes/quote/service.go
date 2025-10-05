package quote

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/quote_input"
)

func Run(cfg models.Config) error {
	languageQuotes, err := LoadQuotes(cfg.Language)
	if err != nil {
		return fmt.Errorf("error loading quotes: %w", err)
	}

	p := tea.NewProgram(quote_input.InitialModel(languageQuotes))

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
