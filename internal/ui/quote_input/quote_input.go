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
	"github.com/neilsmahajan/typing-test-tui/internal/ui/theme"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/typing"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Model struct {
	// Target text
	Target string
	// what user has currentText so far
	currentText    textarea.Model
	languageQuotes models.LanguageQuotes
	rng            *rand.Rand
	viewportWidth  int
	styles         theme.Styles
	session        typing.Session
}

func InitialModel(languageQuotes models.LanguageQuotes) Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	quote := randomQuote(languageQuotes, rng)
	styles := theme.DefaultStyles()
	session := typing.NewSession()

	ti := textarea.New()
	ti.Placeholder = quote.Text
	ti.SetWidth(typing.DefaultBoxWidth)
	ti.Focus()

	return Model{
		Target:         quote.Text,
		currentText:    ti,
		languageQuotes: languageQuotes,
		rng:            rng,
		styles:         styles,
		session:        session,
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
		if m.session.Finished() {
			m.session.Reset()
			m.currentText.SetValue("")
			quote := randomQuote(m.languageQuotes, m.rng)
			m.Target = quote.Text
			m.currentText.Placeholder = m.Target
			m.currentText.SetWidth(m.contentWidth())
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

	if !m.currentText.Focused() {
		m.currentText.Focus()
	}

	m.currentText, cmd = m.currentText.Update(msg)

	if !m.session.Started() && m.currentText.Value() != "" {
		m.session.Start(time.Now())
	}

	// check if completed (capture finish time & wpm only once)
	if !m.session.Finished() && m.currentText.Value() == m.Target {
		m.session.Finish(time.Now(), m.Target)
	}

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

	if m.session.Finished() {
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

	complete := m.styles.Typed.Render(correctSegment) + m.styles.Incorrect.Render(typing.MakeSpacesVisible(incorrectSegment))

	if typedLen > targetLen {
		extra := typed[targetLen:]
		if extra != "" {
			complete += m.styles.Incorrect.Render(typing.MakeSpacesVisible(extra))
		}
	}

	remainingAfterCursor := ""
	if typedLen < targetLen {
		remainingAfterCursor = target[typedLen:]
	}

	if !m.session.Finished() {
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
	if wpm := m.session.CurrentWPM(time.Now(), typed); wpm > 0 {
		wpmValue = fmt.Sprintf("%.1f", wpm)
	}

	elapsedValue := typing.FormatDuration(m.session.Elapsed(time.Now()))

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
	duration := m.session.Elapsed(time.Now())
	summary := fmt.Sprintf("✅ Completed in %s · WPM %.2f", typing.FormatDuration(duration), m.session.WPM())
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
	words := typing.WordCount(m.Target)
	chars := utf8.RuneCountInString(m.Target)
	info := fmt.Sprintf("Language: %s · %d words · %d chars", languageName, words, chars)
	return m.styles.Subtitle.MaxWidth(width).Render(info)
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
		width := m.viewportWidth - typing.BoxHorizontalMargin
		if width < minOuter {
			width = minOuter
		}
		return width
	}

	targetWidth := lipgloss.Width(m.Target)
	inner := typing.DefaultBoxWidth
	if targetWidth > 0 && targetWidth < typing.DefaultBoxWidth {
		inner = targetWidth
	}
	outer := inner + frame
	if outer < minOuter {
		outer = minOuter
	}
	return outer
}

func displayLanguage(lang models.Language) string {
	value := strings.ReplaceAll(strings.ReplaceAll(string(lang), "_", " "), "-", " ")
	value = strings.TrimSpace(strings.TrimPrefix(value, "code "))
	if value == "" {
		return "Unknown"
	}
	return cases.Title(language.English).String(value)
}
