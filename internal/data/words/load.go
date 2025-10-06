package words

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

//go:embed *.json
var wordFiles embed.FS

func Load(language models.Language) (models.LanguageWords, error) {
	filename := fmt.Sprintf("%s.json", language)

	data, err := wordFiles.ReadFile(filename)
	if err != nil {
		return models.LanguageWords{}, fmt.Errorf("words: load %s: %w", language, err)
	}

	var languageWords models.LanguageWords
	if err := json.Unmarshal(data, &languageWords); err != nil {
		return models.LanguageWords{}, fmt.Errorf("words: decode %s: %w", language, err)
	}

	if len(languageWords.Words) == 0 {
		return models.LanguageWords{}, fmt.Errorf("words: no entries for %s", language)
	}

	languageWords.Language = language

	return languageWords, nil
}
