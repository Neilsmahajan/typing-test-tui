package typing

import (
	"testing"
	"time"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestMakeSpacesVisible(t *testing.T) {
	if got := MakeSpacesVisible("a b"); got != "a_b" {
		t.Fatalf("expected spaces to be visible, got %q", got)
	}
}

func TestFormatDuration(t *testing.T) {
	if got := FormatDuration(75 * time.Second); got != "01:15" {
		t.Fatalf("expected formatted duration 01:15, got %q", got)
	}
}

func TestWordCount(t *testing.T) {
	if got := WordCount("one two  three\n"); got != 3 {
		t.Fatalf("expected word count 3, got %d", got)
	}
}

func TestDisplayLanguage(t *testing.T) {
	if got := DisplayLanguage(models.Go); got != "Go" {
		t.Fatalf("expected display language Go, got %q", got)
	}
	if got := DisplayLanguage(models.Language("code_typescript")); got != "Typescript" {
		t.Fatalf("expected display language Typescript, got %q", got)
	}
}
