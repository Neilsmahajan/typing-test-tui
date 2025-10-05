package quotes

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

//go:embed *.json
var quoteFiles embed.FS

func Load(language models.Language) (models.LanguageQuotes, error) {
	filename := fmt.Sprintf("%s.json", language)

	data, err := quoteFiles.ReadFile(filename)
	if err != nil {
		return models.LanguageQuotes{}, fmt.Errorf("quotes: load %s: %w", language, err)
	}

	var languageQuotes models.LanguageQuotes
	if err := json.Unmarshal(data, &languageQuotes); err != nil {
		return models.LanguageQuotes{}, fmt.Errorf("quotes: decode %s: %w", language, err)
	}

	if len(languageQuotes.Quotes) == 0 {
		return models.LanguageQuotes{}, fmt.Errorf("quotes: no entries for %s", language)
	}

	languageQuotes.Language = language

	return languageQuotes, nil
}
