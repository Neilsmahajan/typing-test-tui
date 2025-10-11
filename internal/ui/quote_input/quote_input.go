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
	"github.com/neilsmahajan/typing-test-tui/internal/ui/theme"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/typing"
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
	prevValue := m.currentText.Value()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
		m.currentText.SetWidth(metrics.ContentWidth)
		return m, nil
	case tea.KeyMsg:
		if m.session.Finished() {
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEnter:
				m.session.Reset()
				m.currentText.SetValue("")
				quote := randomQuote(m.languageQuotes, m.rng)
				m.Target = quote.Text
				m.currentText.Placeholder = m.Target
				metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
				m.currentText.SetWidth(metrics.ContentWidth)
			}
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

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter && !m.session.Finished() {
		prevNormalized := normalizeTypedValue(prevValue, m.Target)
		currentNormalized := normalizeTypedValue(m.currentText.Value(), m.Target)

		if len([]rune(currentNormalized)) > len([]rune(prevNormalized)) {
			targetRunes := []rune(m.Target)
			currentRunes := []rune(currentNormalized)
			if prefixMatches(targetRunes, currentRunes) {
				if indent := indentAfter(targetRunes, len(currentRunes)); indent != "" {
					m.currentText.InsertString(indent)
				}
			}
		}

	}

	typedNormalized := normalizeTypedValue(m.currentText.Value(), m.Target)

	if !m.session.Started() && typedNormalized != "" {
		m.session.Start(time.Now())
	}

	// check if completed (capture finish time & wpm only once)
	if !m.session.Finished() && typedNormalized == m.Target {
		m.session.Finish(time.Now(), m.Target)
	}

	return m, cmd
}

// View defines UI rendering
func (m Model) View() string {
	typed := normalizeTypedValue(m.currentText.Value(), m.Target)
	metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
	now := time.Now()

	sections := []string{
		m.renderHeader(metrics.OuterWidth),
		m.renderSubtitle(metrics.OuterWidth),
		typing.RenderBox(typing.BoxConfig{
			Target:        m.Target,
			Typed:         typed,
			Styles:        m.styles,
			Session:       &m.session,
			Metrics:       metrics,
			ViewportWidth: m.viewportWidth,
		}),
		typing.RenderStats(typing.StatsConfig{
			Target:  m.Target,
			Typed:   typed,
			Width:   metrics.OuterWidth,
			Styles:  m.styles,
			Session: &m.session,
			Now:     now,
		}),
	}

	if m.session.Finished() {
		sections = append(sections, typing.RenderCompletion(typing.CompletionConfig{
			Width:   metrics.OuterWidth,
			Styles:  m.styles,
			Session: &m.session,
			Now:     now,
			Prompt:  "Press Enter for another quote or Ctrl+C to exit.",
		}))
	} else {
		sections = append(sections, typing.RenderInstructions(typing.InstructionsConfig{
			Width:  metrics.OuterWidth,
			Styles: m.styles,
		}))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return "\n" + m.styles.Container.Width(metrics.OuterWidth).Render(body)
}

func (m Model) renderHeader(width int) string {
	return m.styles.Header.MaxWidth(width).Render("Quote Mode")
}

func (m Model) renderSubtitle(width int) string {
	languageName := typing.DisplayLanguage(m.languageQuotes.Language)
	words := typing.WordCount(m.Target)
	chars := utf8.RuneCountInString(m.Target)
	info := fmt.Sprintf("Language: %s · %d words · %d chars", languageName, words, chars)
	return m.styles.Subtitle.MaxWidth(width).Render(info)
}

func normalizeTypedValue(typed, target string) string {
	typedRunes := []rune(typed)
	targetRunes := []rune(target)

	result := make([]rune, 0, len(typedRunes))
	ti, to := 0, 0

	for ti < len(typedRunes) && to < len(targetRunes) {
		currentTyped := typedRunes[ti]
		currentTarget := targetRunes[to]

		if currentTarget == '\t' {
			if currentTyped == '\t' {
				result = append(result, '\t')
				ti++
			} else {
				spaces := 0
				for j := ti; j < len(typedRunes) && typedRunes[j] == ' '; j++ {
					spaces++
					if spaces == 4 {
						break
					}
				}
				if spaces == 4 {
					result = append(result, '\t')
					ti += spaces
				} else {
					result = append(result, currentTyped)
					ti++
				}
			}
			to++
			continue
		}

		result = append(result, currentTyped)
		ti++
		to++
	}

	for ; ti < len(typedRunes); ti++ {
		result = append(result, typedRunes[ti])
	}

	return string(result)
}

func prefixMatches(target, typed []rune) bool {
	if len(typed) > len(target) {
		return false
	}
	for i := range typed {
		if target[i] != typed[i] {
			return false
		}
	}
	return true
}

func indentAfter(target []rune, position int) string {
	if position < 0 || position >= len(target) {
		return ""
	}

	var builder strings.Builder
	for i := position; i < len(target); i++ {
		r := target[i]
		if r == ' ' || r == '\t' {
			builder.WriteRune(r)
			continue
		}
		break
	}

	return builder.String()
}
