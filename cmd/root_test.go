package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/spf13/cobra"
)

func TestValidateFlagsQuoteMode(t *testing.T) {
	if err := validateFlags(models.QuoteMode, defaultDuration, defaultWordCount, false, false); err != nil {
		t.Fatalf("expected no error for quote mode defaults, got %v", err)
	}
}

func TestValidateFlagsWordsModeInvalidCount(t *testing.T) {
	err := validateFlags(models.WordsMode, defaultDuration, 5, false, false)
	if err == nil || !strings.Contains(err.Error(), "word count") {
		t.Fatalf("expected word count error, got %v", err)
	}
}

func TestValidateFlagsTimeModeInvalidDuration(t *testing.T) {
	err := validateFlags(models.TimeMode, 10, defaultWordCount, false, false)
	if err == nil || !strings.Contains(err.Error(), "duration") {
		t.Fatalf("expected duration error, got %v", err)
	}
}

func TestNormalizeLanguageAlias(t *testing.T) {
	lang, err := normalizeLanguage("go")
	if err != nil {
		t.Fatalf("expected alias normalization to succeed, got %v", err)
	}
	if lang != models.Go {
		t.Fatalf("expected normalized language to be %q, got %q", models.Go, lang)
	}
}

func TestNormalizeLanguageUnsupported(t *testing.T) {
	if _, err := normalizeLanguage("klingon"); err == nil {
		t.Fatalf("expected error for unsupported language")
	}
}

func TestJoinInts(t *testing.T) {
	result := joinInts([]int{1, 2, 3})
	if result != "1, 2, 3" {
		t.Fatalf("unexpected joinInts result: %q", result)
	}
}

func TestListLanguagesOutput(t *testing.T) {
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	listLanguages(cmd, nil)

	output := buf.String()
	if !strings.Contains(output, "Supported Languages:") {
		t.Fatalf("expected header in output, got %q", output)
	}
	if !strings.Contains(output, "english") {
		t.Fatalf("expected english in output, got %q", output)
	}
}

func TestListModesOutput(t *testing.T) {
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	listModes(cmd, nil)

	output := buf.String()
	if !strings.Contains(output, "Supported Modes:") {
		t.Fatalf("expected header in output, got %q", output)
	}
	if !strings.Contains(output, "quote") || !strings.Contains(output, "words") || !strings.Contains(output, "time") {
		t.Fatalf("expected all modes in output, got %q", output)
	}
}
