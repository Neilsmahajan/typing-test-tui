package typing

import (
	"fmt"
	"math"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/theme"
)

type BoxMetrics struct {
	OuterWidth   int
	ContentWidth int
}

type BoxConfig struct {
	Target        string
	Typed         string
	Styles        theme.Styles
	Session       *Session
	Metrics       BoxMetrics
	ViewportWidth int
}

type StatsConfig struct {
	Target        string
	Typed         string
	Width         int
	Styles        theme.Styles
	Session       *Session
	Now           time.Time
	ProgressLabel string
	WPMLabel      string
	TimeLabel     string
	ProgressValue string
	WPMValue      string
	TimeValue     string
}

type CompletionConfig struct {
	Width   int
	Styles  theme.Styles
	Session *Session
	Now     time.Time
	Prompt  string
}

type InstructionsConfig struct {
	Width   int
	Styles  theme.Styles
	Message string
}

const DefaultInstructionsMessage = "Esc: blur focus • Ctrl+C: exit"

func ComputeBoxMetrics(target string, styles theme.Styles, viewportWidth int) BoxMetrics {
	frame := styles.QuoteBox.GetHorizontalFrameSize()
	outer := computeOuterWidth(target, styles, viewportWidth, frame)
	inner := outer - frame
	if inner < 1 {
		inner = 1
	}
	return BoxMetrics{
		OuterWidth:   outer,
		ContentWidth: inner,
	}
}

func RenderBox(cfg BoxConfig) string {
	metrics := cfg.Metrics
	if metrics.OuterWidth == 0 || metrics.ContentWidth == 0 {
		metrics = ComputeBoxMetrics(cfg.Target, cfg.Styles, cfg.ViewportWidth)
	}

	target := cfg.Target
	typed := cfg.Typed
	typedLen := len(typed)
	targetLen := len(target)
	limit := typedLen
	if limit > targetLen {
		limit = targetLen
	}

	incorrectIndex := limit
	for i := 0; i < limit; i++ {
		if typed[i] != target[i] {
			incorrectIndex = i
			break
		}
	}

	correctSegment := target[:incorrectIndex]
	incorrectSegment := ""
	if incorrectIndex < limit {
		incorrectSegment = target[incorrectIndex:limit]
	}

	complete := cfg.Styles.Typed.Render(correctSegment) + cfg.Styles.Incorrect.Render(MakeSpacesVisible(incorrectSegment))

	if typedLen > targetLen {
		extra := typed[targetLen:]
		if extra != "" {
			complete += cfg.Styles.Incorrect.Render(MakeSpacesVisible(extra))
		}
	}

	remainingAfterCursor := ""
	if typedLen < targetLen {
		remainingAfterCursor = target[typedLen:]
	}

	if cfg.Session != nil && !cfg.Session.Finished() {
		cursorGlyph := " "
		if len(remainingAfterCursor) > 0 {
			r, size := utf8.DecodeRuneInString(remainingAfterCursor)
			if size > 0 {
				if r == utf8.RuneError {
					cursorGlyph = remainingAfterCursor[:size]
				} else {
					cursorGlyph = string(r)
				}
				remainingAfterCursor = remainingAfterCursor[size:]
			}
		}
		complete += cfg.Styles.Cursor.Render(cursorGlyph)
	}

	complete += cfg.Styles.Remaining.Render(remainingAfterCursor)
	innerWidth := metrics.ContentWidth
	wrapped := cfg.Styles.QuoteContent.Width(innerWidth).Render(complete)
	return cfg.Styles.QuoteBox.Width(metrics.OuterWidth).Render(wrapped)
}

func RenderStats(cfg StatsConfig) string {
	if cfg.Session == nil {
		return ""
	}

	now := cfg.Now
	if now.IsZero() {
		now = time.Now()
	}

	totalRunes := utf8.RuneCountInString(cfg.Target)
	typedRunes := utf8.RuneCountInString(cfg.Typed)
	progress := fmt.Sprintf("%d/%d chars", typedRunes, totalRunes)
	if totalRunes > 0 {
		percent := int(math.Round(float64(typedRunes) / float64(totalRunes) * 100))
		if percent > 100 {
			percent = 100
		}
		progress = fmt.Sprintf("%s (%d%%)", progress, percent)
	}

	if cfg.ProgressValue != "" {
		progress = cfg.ProgressValue
	}

	progressLabel := cfg.ProgressLabel
	if progressLabel == "" {
		progressLabel = "Progress"
	}
	wpmLabel := cfg.WPMLabel
	if wpmLabel == "" {
		wpmLabel = "WPM"
	}
	timeLabel := cfg.TimeLabel
	if timeLabel == "" {
		timeLabel = "Time"
	}

	wpmValue := "--"
	if cfg.WPMValue != "" {
		wpmValue = cfg.WPMValue
	} else if wpm := cfg.Session.CurrentWPM(now, cfg.Typed); wpm > 0 {
		wpmValue = fmt.Sprintf("%.1f", wpm)
	}

	elapsedValue := FormatDuration(cfg.Session.Elapsed(now))
	if cfg.TimeValue != "" {
		elapsedValue = cfg.TimeValue
	}

	statEntries := []string{
		renderStatBlock(cfg.Styles, progressLabel, progress),
		renderStatBlock(cfg.Styles, wpmLabel, wpmValue),
		renderStatBlock(cfg.Styles, timeLabel, elapsedValue),
	}

	row := statEntries[0]
	for i := 1; i < len(statEntries); i++ {
		row = lipgloss.JoinHorizontal(lipgloss.Left, row, cfg.Styles.StatSeparator, statEntries[i])
	}

	return cfg.Styles.StatsRow.MaxWidth(cfg.Width).Render(row)
}

func RenderCompletion(cfg CompletionConfig) string {
	if cfg.Session == nil {
		return ""
	}

	now := cfg.Now
	if now.IsZero() {
		now = time.Now()
	}

	duration := cfg.Session.Elapsed(now)
	summary := fmt.Sprintf("✅ Completed in %s · WPM %.2f", FormatDuration(duration), cfg.Session.WPM())
	prompt := cfg.Prompt
	if prompt == "" {
		prompt = "Press any key to continue or Ctrl+C to exit."
	}

	lines := []string{
		cfg.Styles.Success.MaxWidth(cfg.Width).Render(summary),
		cfg.Styles.Instruction.MaxWidth(cfg.Width).Render(prompt),
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func RenderInstructions(cfg InstructionsConfig) string {
	message := cfg.Message
	if message == "" {
		message = DefaultInstructionsMessage
	}
	return cfg.Styles.Instruction.MaxWidth(cfg.Width).Render(message)
}

func computeOuterWidth(target string, styles theme.Styles, viewportWidth int, frame int) int {
	minOuter := frame + 1

	if viewportWidth > 0 {
		width := viewportWidth - BoxHorizontalMargin
		if width < minOuter {
			width = minOuter
		}
		return width
	}

	targetWidth := lipgloss.Width(target)
	inner := DefaultBoxWidth
	if targetWidth > 0 && targetWidth < DefaultBoxWidth {
		inner = targetWidth
	}
	outer := inner + frame
	if outer < minOuter {
		outer = minOuter
	}
	return outer
}

func renderStatBlock(styles theme.Styles, label, value string) string {
	block := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.StatLabel.Render(label),
		styles.StatValue.Render(value),
	)
	return styles.StatBlock.Render(block)
}
