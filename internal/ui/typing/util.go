package typing

import (
	"fmt"
	"strings"
	"time"
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
