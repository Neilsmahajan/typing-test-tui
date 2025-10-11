package quote

import (
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestRunInvalidLanguage(t *testing.T) {
	cfg := models.Config{Mode: models.QuoteMode, Language: models.Language("not-real")}
	if err := Run(cfg); err == nil {
		t.Fatalf("expected error when language data is missing")
	}
}
