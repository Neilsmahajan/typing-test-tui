package words

import (
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestRunInvalidLanguage(t *testing.T) {
	cfg := models.Config{Mode: models.WordsMode, Language: models.Language("not-real"), WordCount: models.WordCount(10)}
	if err := Run(cfg); err == nil {
		t.Fatalf("expected error when language data is missing")
	}
}
