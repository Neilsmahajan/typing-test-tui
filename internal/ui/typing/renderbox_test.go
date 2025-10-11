package typing

import (
	"strings"
	"testing"

	"github.com/neilsmahajan/typing-test-tui/internal/ui/theme"
)

func TestRenderInlinePreservesContent(t *testing.T) {
	styles := theme.DefaultStyles()
	input := "package main\n"
	styled := renderInline(styles.Typed, input)
	if cleaned := sanitizeANSI(styled); cleaned != input {
		t.Fatalf("expected sanitized inline render to equal %q, got %q", input, cleaned)
	}
}

func TestRenderBoxCursorAlignment(t *testing.T) {
	styles := theme.DefaultStyles()
	target := "package main\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}"
	typed := "package main\n"
	metrics := ComputeBoxMetrics(target, styles, 0)
	session := NewSession()

	output := RenderBox(BoxConfig{
		Target:        target,
		Typed:         typed,
		Styles:        styles,
		Session:       &session,
		Metrics:       metrics,
		ViewportWidth: 0,
	})

	cleaned := sanitizeANSI(output)
	lines := strings.Split(cleaned, "\n")
	var importLine string
	for _, line := range lines {
		if strings.Contains(line, "import \"fmt\"") {
			importLine = line
			break
		}
	}

	if importLine == "" {
		t.Fatalf("rendered output missing import line:\n%s", cleaned)
	}

	const expectedPrefix = "│  import"
	if !strings.HasPrefix(importLine, expectedPrefix) {
		t.Fatalf("expected import line to start with %q, got %q", expectedPrefix, importLine)
	}
}

func sanitizeANSI(input string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range input {
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		if r == 0x1b {
			inEscape = true
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
