package loaders

import (
	wordsdata "github.com/neilsmahajan/typing-test-tui/internal/data/words"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func LoadWords(language models.Language) (models.LanguageWords, error) {
	languageWords, err := wordsdata.Load(language)
	if err != nil {
		return models.LanguageWords{}, err
	}
	return languageWords, nil
}
