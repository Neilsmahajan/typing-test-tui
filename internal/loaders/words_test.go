package loaders

import (
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestLoadWordsSuccess(t *testing.T) {
	data, err := LoadWords(models.English)
	if err != nil {
		t.Fatalf("expected words to load, got %v", err)
	}
	if data.Language != models.English {
		t.Fatalf("expected language to be %q, got %q", models.English, data.Language)
	}
	if len(data.Words) == 0 {
		t.Fatalf("expected words to be non-empty")
	}
}

func TestLoadWordsInvalidLanguage(t *testing.T) {
	if _, err := LoadWords(models.Language("not-real")); err == nil {
		t.Fatalf("expected error for missing language")
	}
}
