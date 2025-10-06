package quote_input

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

var (
	typedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	remainingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	boxStyle       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	cursorStyle    = lipgloss.NewStyle().Reverse(true)
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
	err            error
}

func InitialModel(languageQuotes models.LanguageQuotes) Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	quote := randomQuote(languageQuotes, rng)

	ti := textarea.New()
	ti.Placeholder = quote.Text
	ti.Focus()

	return Model{
		Target:         quote.Text,
		currentText:    ti,
		languageQuotes: languageQuotes,
		rng:            rng,
		err:            nil,
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
	case tea.KeyMsg:
		if m.finished {
			m.finished = false
			m.started = false
			m.currentText.SetValue("")
			m.Target = randomQuote(m.languageQuotes, m.rng).Text
			m.currentText.Placeholder = m.Target
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
		m.err = msg
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
	var b strings.Builder

	b.WriteString("\nType the following:\n\n")
	b.WriteString(m.Target + "\n\n")

	typed := m.currentText.Value()
	remaining := ""
	if len(typed) < len(m.Target) {
		remaining = m.Target[len(typed):]
	}

	b.WriteString(m.renderBox(typed, remaining) + "\n\n")

	if m.finished {
		b.WriteString(fmt.Sprintf("✅ Done! WPM: %.2f\n", m.wpm))
		b.WriteString("Press any key to continue. Press Ctrl+C to exit.\n")
	}

	return b.String()
}

func (m Model) renderBox(typed string, remaining string) string {
	incorrectIndex := len(typed)
	for i := 0; i < len(typed); i++ {
		if typed[i] != m.Target[i] {
			incorrectIndex = i
			break
		}
	}
	correctSegment := typed[:incorrectIndex]
	incorrectSegment := m.Target[incorrectIndex:len(typed)]
	complete := typedStyle.Render(correctSegment) + incorrectStyle.Render(makeSpacesVisible(incorrectSegment))
	remainingAfterCursor := remaining
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
		complete += cursorStyle.Render(cursorGlyph)
	}
	complete += remainingStyle.Render(remainingAfterCursor)
	targetWidth := lipgloss.Width(m.Target)
	return boxStyle.Width(targetWidth).Render(complete)
}

func makeSpacesVisible(text string) string {
	if text == "" {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))

	for _, r := range text {
		if r == ' ' {
			b.WriteRune('_')
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
