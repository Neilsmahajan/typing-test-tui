package time

import (
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestRunInvalidLanguage(t *testing.T) {
	cfg := models.Config{Mode: models.TimeMode, Language: models.Language("not-real"), Duration: models.Duration(60)}
	if err := Run(cfg); err == nil {
		t.Fatalf("expected error when language data is missing")
	}
}
