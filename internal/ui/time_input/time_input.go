package time_input

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
	currentText   textarea.Model
	duration      models.Duration
	languageWords models.LanguageWords
	rng           *rand.Rand
	viewportWidth int
	styles        theme.Styles
	session       typing.Session
}

func InitialModel(languageWords models.LanguageWords, duration models.Duration) Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := generateTargetWords(rng, languageWords, duration)

	ti := textarea.New()
	ti.Placeholder = target
	ti.SetWidth(typing.DefaultBoxWidth)
	ti.Focus()

	return Model{
		Target:        target,
		currentText:   ti,
		duration:      duration,
		languageWords: languageWords,
		rng:           rng,
		styles:        theme.DefaultStyles(),
		session:       typing.NewSession(),
	}
}

func generateTargetWords(rng *rand.Rand, languageWords models.LanguageWords, duration models.Duration) string {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	available := languageWords.Words
	count := int(duration / 2) // assuming average 2 words per second
	if len(available) == 0 || count <= 0 {
		return ""
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		idx := rng.Intn(len(available))
		result[i] = available[idx]
	}
	return strings.Join(result, " ")
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
		metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
		m.currentText.SetWidth(metrics.ContentWidth)
		return m, nil
	case tea.KeyMsg:
		if m.session.Finished() {
			m.session.Reset()
			m.currentText.SetValue("")
			target := generateTargetWords(m.rng, m.languageWords, m.duration)
			m.Target = target
			m.currentText.Placeholder = m.Target
			metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
			m.currentText.SetWidth(metrics.ContentWidth)
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
			Prompt:  "Press any key for another word set or Ctrl+C to exit.",
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
	return m.styles.Header.MaxWidth(width).Render("Time Mode")
}

func (m Model) renderSubtitle(width int) string {
	languageName := typing.DisplayLanguage(m.languageWords.Language)
	chars := utf8.RuneCountInString(m.Target)
	info := fmt.Sprintf("Language: %s · duration: %d · %d chars", languageName, m.duration, chars)
	return m.styles.Subtitle.MaxWidth(width).Render(info)
}
