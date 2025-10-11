package theme

import "testing"

func TestDefaultStyles(t *testing.T) {
	styles := DefaultStyles()
	if styles.StatSeparator == "" {
		t.Fatalf("expected stat separator to be initialized")
	}
	rendered := styles.Header.Render("Title")
	if rendered == "" {
		t.Fatalf("expected header style to render non-empty string")
	}
	if styles.Cursor.Render(" ") == "" {
		t.Fatalf("expected cursor style to render non-empty string")
	}
}
