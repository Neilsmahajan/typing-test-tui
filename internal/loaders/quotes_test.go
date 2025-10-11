package loaders

import (
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestLoadQuotesSuccess(t *testing.T) {
	data, err := LoadQuotes(models.English)
	if err != nil {
		t.Fatalf("expected quotes to load, got %v", err)
	}
	if data.Language != models.English {
		t.Fatalf("expected language to be %q, got %q", models.English, data.Language)
	}
	if len(data.Quotes) == 0 {
		t.Fatalf("expected quotes to be non-empty")
	}
}

func TestLoadQuotesInvalidLanguage(t *testing.T) {
	if _, err := LoadQuotes(models.Language("not-real")); err == nil {
		t.Fatalf("expected error for missing language")
	}
}
