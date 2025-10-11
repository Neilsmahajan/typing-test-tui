package models

import "testing"

func TestNormalizeLanguageAlias(t *testing.T) {
	lang, ok := NormalizeLanguage("code-go")
	if !ok {
		t.Fatalf("expected alias to normalize")
	}
	if lang != Go {
		t.Fatalf("expected language %q, got %q", Go, lang)
	}
}

func TestNormalizeLanguageUnsupported(t *testing.T) {
	if _, ok := NormalizeLanguage("unsupported-lang"); ok {
		t.Fatalf("expected unsupported language to return false")
	}
}

func TestSupportedLanguagesReturnsCopy(t *testing.T) {
	langs := SupportedLanguages()
	if len(langs) == 0 {
		t.Fatalf("expected supported languages to be non-empty")
	}
	originalFirst := langs[0]
	langs[0] = Language("modified")
	fresh := SupportedLanguages()
	if fresh[0] != originalFirst {
		t.Fatalf("expected SupportedLanguages to return copy; got %q instead of %q", fresh[0], originalFirst)
	}
}
