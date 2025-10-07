package app

import (
	"fmt"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/modes/quote"
	"github.com/neilsmahajan/typing-test-tui/internal/modes/time"
	"github.com/neilsmahajan/typing-test-tui/internal/modes/words"
)

func Run(cfg models.Config) error {
	switch cfg.Mode {
	case models.QuoteMode:
		return quote.Run(cfg)
	case models.WordsMode:
		return words.Run(cfg)
	case models.TimeMode:
		return time.Run(cfg)
	default:
		return fmt.Errorf("unsupported mode: %s", cfg.Mode)
	}
}
