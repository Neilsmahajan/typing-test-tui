package loaders

import (
	"fmt"

	quotesdata "github.com/neilsmahajan/typing-test-tui/internal/data/quotes"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func LoadQuotes(language models.Language) (models.LanguageQuotes, error) {
	languageQuotes, err := quotesdata.Load(language)
	if err != nil {
		return models.LanguageQuotes{}, fmt.Errorf("load quotes: %w", err)
	}

	return languageQuotes, nil
}
