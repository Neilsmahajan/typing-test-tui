package app

import (
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

func TestRunUnsupportedMode(t *testing.T) {
	cfg := models.Config{Mode: models.Mode("unsupported")}
	if err := Run(cfg); err == nil {
		t.Fatalf("expected error for unsupported mode")
	}
}
