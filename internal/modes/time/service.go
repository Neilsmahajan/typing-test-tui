package time

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/internal/loaders"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/time_input"
)

func Run(cfg models.Config) error {
	languageWords, err := loaders.LoadWords(cfg.Language)
	if err != nil {
		return fmt.Errorf("error loading words: %w", err)
	}

	p := tea.NewProgram(time_input.InitialModel(languageWords, cfg.Duration))

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
