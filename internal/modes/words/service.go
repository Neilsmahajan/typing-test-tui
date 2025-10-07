package words

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/internal/loaders"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/words_input"
)

func Run(cfg models.Config) error {
	languageWords, err := loaders.LoadWords(cfg.Language)
	if err != nil {
		return fmt.Errorf("error loading words: %w", err)
	}

	p := tea.NewProgram(words_input.InitialModel(languageWords, cfg.WordCount))

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
