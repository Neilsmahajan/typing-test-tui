package quote_input

import (
	"testing"
	"time"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/typing"
)

func TestNormalizeTypedValueConvertsSpacesToTabs(t *testing.T) {
	target := "\tfunc main()"
	typed := "    func main()"
	normalized := normalizeTypedValue(typed, target)
	if normalized != target {
		t.Fatalf("expected normalized value %q, got %q", target, normalized)
	}
}

func TestPrefixMatches(t *testing.T) {
	target := []rune("hello")
	if !prefixMatches(target, []rune("hel")) {
		t.Fatalf("expected prefixMatches to return true for matching prefix")
	}
	if prefixMatches(target, []rune("hey")) {
		t.Fatalf("expected prefixMatches to return false for mismatch")
	}
}

func TestIndentAfter(t *testing.T) {
	target := []rune("\t    body")
	indent := indentAfter(target, 1)
	if indent != "    " {
		t.Fatalf("expected indent of four spaces, got %q", indent)
	}
}

func TestInitialModelSetsIndicators(t *testing.T) {
	codeQuotes := models.LanguageQuotes{
		Language: models.Go,
		Quotes:   []models.Quote{{Text: "code sample"}},
	}
	codeModel := InitialModel(codeQuotes)
	if codeModel.newlineIndicator != typing.DefaultNewlineIndicator {
		t.Fatalf("expected newline indicator for code language")
	}

	plainQuotes := models.LanguageQuotes{
		Language: models.English,
		Quotes:   []models.Quote{{Text: "hello world"}},
	}
	plainModel := InitialModel(plainQuotes)
	if plainModel.newlineIndicator != "" {
		t.Fatalf("expected no newline indicator for natural language")
	}
	if plainModel.session.Started() {
		t.Fatalf("expected new model session to be not started")
	}
	if plainModel.currentText.Placeholder != "hello world" {
		t.Fatalf("expected placeholder to match quote text")
	}
}

func TestModelUpdateIgnoresErrors(t *testing.T) {
	codeQuotes := models.LanguageQuotes{
		Language: models.English,
		Quotes:   []models.Quote{{Text: "hello"}},
	}
	model := InitialModel(codeQuotes)
	// Should not panic when receiving an unexpected message type
	model.Update(time.Now())
}
