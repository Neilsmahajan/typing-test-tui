package app

import (
	"fmt"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/modes/quote"
)

func Run(cfg models.Config) error {
	switch cfg.Mode {
	case models.QuoteMode:
		return quote.Run(cfg)
	case models.WordsMode:
		// Run words mode
		return fmt.Errorf("words mode not implemented yet")
	case models.TimeMode:
		// Run time mode
		return fmt.Errorf("time mode not implemented yet")
	default:
		return fmt.Errorf("unsupported mode: %s", cfg.Mode)
	}
}
