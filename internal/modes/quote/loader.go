package quote

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func LoadQuotes(language models.Language) (*models.LanguageQuotes, error) {
	filePath := fmt.Sprintf("internal/data/quotes/%s.json", language)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var languageQuotes models.LanguageQuotes
	if err := json.Unmarshal(content, &languageQuotes); err != nil {
		return nil, err
	}

	return &languageQuotes, nil
}
