package typing

import (
	"fmt"
	"strings"
	"time"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MakeSpacesVisible(text string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' {
			return '_'
		}
		return r
	}, text)
}

func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "00:00"
	}
	seconds := int(d.Seconds())
	minutes := seconds / 60
	remaining := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, remaining)
}

func WordCount(text string) int {
	return len(strings.Fields(text))
}

func DisplayLanguage(lang models.Language) string {
	value := strings.ReplaceAll(strings.ReplaceAll(string(lang), "_", " "), "-", " ")
	value = strings.TrimSpace(strings.TrimPrefix(value, "code "))
	if value == "" {
		return "Unknown"
	}
	return cases.Title(language.English).String(value)
}
