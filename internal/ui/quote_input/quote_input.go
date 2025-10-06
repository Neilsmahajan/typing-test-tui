package quote_input

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	boxHorizontalMargin = 4
	defaultBoxWidth     = 60
)

type Model struct {
	// Target text
	Target string
	// what user has currentText so far
	currentText textarea.Model
	// timing
	started        bool
	start          time.Time
	finished       bool
	end            time.Time
	wpm            float64
	languageQuotes models.LanguageQuotes
	rng            *rand.Rand
	viewportWidth  int
	styles         Styles
}

func InitialModel(languageQuotes models.LanguageQuotes) Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	quote := randomQuote(languageQuotes, rng)
	styles := defaultStyles()

	ti := textarea.New()
	ti.Placeholder = quote.Text
	ti.SetWidth(defaultBoxWidth)
	ti.Focus()

	return Model{
		Target:         quote.Text,
		currentText:    ti,
		languageQuotes: languageQuotes,
		rng:            rng,
		styles:         styles,
	}
}

func randomQuote(languageQuotes models.LanguageQuotes, rng *rand.Rand) models.Quote {
	count := len(languageQuotes.Quotes)
	if count == 0 {
		return models.Quote{}
	}

	return languageQuotes.Quotes[rng.Intn(count)]
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages (key presses, etc.)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.currentText.SetWidth(m.contentWidth())
		return m, nil
	case tea.KeyMsg:
		if m.finished {
			m.finished = false
			m.started = false
			m.currentText.SetValue("")
			m.Target = randomQuote(m.languageQuotes, m.rng).Text
			m.currentText.Placeholder = m.Target
			m.currentText.SetWidth(m.contentWidth())
			m.wpm = 0
			return m, nil
		}

		switch msg.Type {
		case tea.KeyEsc:
			if m.currentText.Focused() {
				m.currentText.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case error:
		return m, nil
	}

	if !m.started && m.currentText.LineCount() == 1 {
		m.started = true
		m.start = time.Now()
	}
	if !m.currentText.Focused() {
		m.currentText.Focus()
	}

	// check if completed (capture finish time & wpm only once)
	if !m.finished && m.currentText.Value() == m.Target {
		m.finished = true
		m.end = time.Now()
		elapsedMinutes := m.end.Sub(m.start).Minutes()
		if elapsedMinutes > 0 {
			m.wpm = float64(len(strings.Fields(m.Target))) / elapsedMinutes
		}
	}

	m.currentText, cmd = m.currentText.Update(msg)

	return m, cmd
}

// View defines UI rendering
func (m Model) View() string {
	typed := m.currentText.Value()
	width := m.boxOuterWidth()

	sections := []string{
		m.renderHeader(width),
		m.renderSubtitle(width),
		m.renderBox(typed),
		m.renderStats(typed, width),
	}

	if m.finished {
		sections = append(sections, m.renderCompletion(width))
	} else {
		sections = append(sections, m.renderInstructions(width))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return "\n" + m.styles.Container.Width(width).Render(body)
}

func (m Model) renderBox(typed string) string {
	target := m.Target
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

	complete := m.styles.Typed.Render(correctSegment) + m.styles.Incorrect.Render(makeSpacesVisible(incorrectSegment))

	if typedLen > targetLen {
		extra := typed[targetLen:]
		if extra != "" {
			complete += m.styles.Incorrect.Render(makeSpacesVisible(extra))
		}
	}

	remainingAfterCursor := ""
	if typedLen < targetLen {
		remainingAfterCursor = target[typedLen:]
	}

	if !m.finished {
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
		complete += m.styles.Cursor.Render(cursorGlyph)
	}

	complete += m.styles.Remaining.Render(remainingAfterCursor)
	innerWidth := m.contentWidth()
	wrapped := m.styles.QuoteContent.Width(innerWidth).Render(complete)
	return m.styles.QuoteBox.Width(m.boxOuterWidth()).Render(wrapped)
}

func (m Model) renderStats(typed string, width int) string {
	totalRunes := utf8.RuneCountInString(m.Target)
	typedRunes := utf8.RuneCountInString(typed)
	progress := fmt.Sprintf("%d/%d chars", typedRunes, totalRunes)
	if totalRunes > 0 {
		percent := int(math.Round(float64(typedRunes) / float64(totalRunes) * 100))
		if percent > 100 {
			percent = 100
		}
		progress = fmt.Sprintf("%s (%d%%)", progress, percent)
	}

	wpmValue := "--"
	if wpm := m.currentWPM(typed); wpm > 0 {
		wpmValue = fmt.Sprintf("%.1f", wpm)
	}

	elapsedValue := formatDuration(m.elapsed())

	statEntries := []string{
		m.renderStat("Progress", progress),
		m.renderStat("WPM", wpmValue),
		m.renderStat("Time", elapsedValue),
	}

	row := statEntries[0]
	for i := 1; i < len(statEntries); i++ {
		row = lipgloss.JoinHorizontal(lipgloss.Left, row, m.styles.StatSeparator, statEntries[i])
	}

	return m.styles.StatsRow.MaxWidth(width).Render(row)
}

func (m Model) renderStat(label, value string) string {
	block := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.StatLabel.Render(label),
		m.styles.StatValue.Render(value),
	)
	return m.styles.StatBlock.Render(block)
}

func (m Model) renderInstructions(width int) string {
	parts := []string{"Esc: blur focus", "Ctrl+C: exit"}
	message := strings.Join(parts, " • ")
	return m.styles.Instruction.MaxWidth(width).Render(message)
}

func (m Model) renderCompletion(width int) string {
	duration := m.end.Sub(m.start)
	summary := fmt.Sprintf("✅ Completed in %s · WPM %.2f", formatDuration(duration), m.wpm)
	lines := []string{
		m.styles.Success.MaxWidth(width).Render(summary),
		m.styles.Instruction.MaxWidth(width).Render("Press any key for another quote or Ctrl+C to exit."),
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) renderHeader(width int) string {
	return m.styles.Header.MaxWidth(width).Render("Quote Mode")
}

func (m Model) renderSubtitle(width int) string {
	languageName := displayLanguage(m.languageQuotes.Language)
	words := wordCount(m.Target)
	chars := utf8.RuneCountInString(m.Target)
	info := fmt.Sprintf("Language: %s · %d words · %d chars", languageName, words, chars)
	return m.styles.Subtitle.MaxWidth(width).Render(info)
}

func makeSpacesVisible(text string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' {
			return '_'
		}
		return r
	}, text)
}

func (m Model) contentWidth() int {
	outer := m.boxOuterWidth()
	frame := m.styles.QuoteBox.GetHorizontalFrameSize()
	inner := outer - frame
	if inner < 1 {
		inner = 1
	}
	return inner
}

func (m Model) boxOuterWidth() int {
	frame := m.styles.QuoteBox.GetHorizontalFrameSize()
	minOuter := frame + 1

	if m.viewportWidth > 0 {
		width := m.viewportWidth - boxHorizontalMargin
		if width < minOuter {
			width = minOuter
		}
		return width
	}

	targetWidth := lipgloss.Width(m.Target)
	inner := defaultBoxWidth
	if targetWidth > 0 && targetWidth < defaultBoxWidth {
		inner = targetWidth
	}
	outer := inner + frame
	if outer < minOuter {
		outer = minOuter
	}
	return outer
}

func (m Model) currentWPM(typed string) float64 {
	if m.finished {
		return m.wpm
	}
	if !m.started {
		return 0
	}
	elapsed := time.Since(m.start).Minutes()
	if elapsed <= 0 {
		return 0
	}
	words := float64(len(strings.Fields(typed)))
	if words == 0 {
		words = float64(utf8.RuneCountInString(typed)) / 5
	}
	if words <= 0 {
		return 0
	}
	return words / elapsed
}

func (m Model) elapsed() time.Duration {
	if !m.started {
		return 0
	}
	if m.finished {
		return m.end.Sub(m.start)
	}
	return time.Since(m.start)
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "00:00"
	}
	seconds := int(d.Seconds())
	minutes := seconds / 60
	remaining := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, remaining)
}

func displayLanguage(lang models.Language) string {
	value := strings.ReplaceAll(strings.ReplaceAll(string(lang), "_", " "), "-", " ")
	value = strings.TrimSpace(strings.TrimPrefix(value, "code "))
	if value == "" {
		return "Unknown"
	}
	return cases.Title(language.English).String(value)
}

func wordCount(text string) int {
	return len(strings.Fields(text))
}
