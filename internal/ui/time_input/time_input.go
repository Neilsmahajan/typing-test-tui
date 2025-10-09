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
	currentText        textarea.Model
	duration           models.Duration
	languageWords      models.LanguageWords
	includePunctuation bool
	includeNumbers     bool
	rng                *rand.Rand
	viewportWidth      int
	styles             theme.Styles
	session            typing.Session
	totalDuration      time.Duration
	remaining          time.Duration
	deadline           time.Time
	tickInterval       time.Duration
}

type tickMsg struct {
	now time.Time
}

const (
	wordsPerSecondEstimate = 6
	minInitialWords        = 120
	wordBufferChunk        = 40
	minRemainingRunes      = 200
	defaultTickInterval    = 100 * time.Millisecond
)

func InitialModel(languageWords models.LanguageWords, duration models.Duration, includePunctuation bool, includeNumbers bool) Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := generateTargetWords(rng, languageWords, duration)
	totalDuration := time.Duration(duration) * time.Second
	if totalDuration <= 0 {
		totalDuration = 60 * time.Second
	}
	styles := theme.DefaultStyles()
	session := typing.NewSession()

	ti := textarea.New()
	ti.Placeholder = target
	ti.SetWidth(typing.DefaultBoxWidth)
	ti.Focus()

	return Model{
		Target:             target,
		currentText:        ti,
		duration:           duration,
		languageWords:      languageWords,
		includePunctuation: includePunctuation,
		includeNumbers:     includeNumbers,
		rng:                rng,
		styles:             styles,
		session:            session,
		totalDuration:      totalDuration,
		remaining:          totalDuration,
		tickInterval:       defaultTickInterval,
	}
}

func generateTargetWords(rng *rand.Rand, languageWords models.LanguageWords, duration models.Duration) string {
	wordCount := estimateInitialWordCount(duration)
	return generateWordString(rng, languageWords, wordCount)
}

func generateWordString(rng *rand.Rand, languageWords models.LanguageWords, count int) string {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	available := languageWords.Words
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

func estimateInitialWordCount(duration models.Duration) int {
	count := int(duration) * wordsPerSecondEstimate
	if count < minInitialWords {
		count = minInitialWords
	}
	return count
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages (key presses, etc.)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
		m.currentText.SetWidth(metrics.ContentWidth)
		return m, nil
	case tickMsg:
		if m.session.Started() && !m.session.Finished() {
			remaining := m.deadline.Sub(msg.now)
			if remaining <= 0 {
				remaining = 0
				m.session.Finish(msg.now, m.currentText.Value())
			}
			m.remaining = remaining
		}
		if m.session.Started() && !m.session.Finished() {
			if tick := m.scheduleTick(); tick != nil {
				cmds = append(cmds, tick)
			}
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		if m.session.Finished() {
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEnter:
				m.session.Reset()
				m.currentText.SetValue("")
				target := generateTargetWords(m.rng, m.languageWords, m.duration)
				m.Target = target
				m.currentText.Placeholder = m.Target
				metrics := typing.ComputeBoxMetrics(m.Target, m.styles, m.viewportWidth)
				m.currentText.SetWidth(metrics.ContentWidth)
				m.remaining = m.totalDuration
				m.deadline = time.Time{}
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

	updated, cmd := m.currentText.Update(msg)
	m.currentText = updated

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	now := time.Now()
	if !m.session.Started() && m.currentText.Value() != "" {
		m.session.Start(now)
		m.remaining = m.totalDuration
		m.deadline = now.Add(m.totalDuration)
		if tick := m.scheduleTick(); tick != nil {
			cmds = append(cmds, tick)
		}
	}

	m.ensureTargetBuffer()

	return m, tea.Batch(cmds...)
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
			Target:        m.Target,
			Typed:         typed,
			Width:         metrics.OuterWidth,
			Styles:        m.styles,
			Session:       &m.session,
			Now:           now,
			ProgressLabel: "Words",
			ProgressValue: fmt.Sprintf("%d", typing.WordCount(typed)),
			TimeLabel:     "Time Left",
			TimeValue:     typing.FormatDuration(m.displayRemaining()),
		}),
	}

	if m.session.Finished() {
		sections = append(sections, typing.RenderCompletion(typing.CompletionConfig{
			Width:   metrics.OuterWidth,
			Styles:  m.styles,
			Session: &m.session,
			Now:     now,
			Prompt:  "Press Enter to start another word set or Ctrl+C to exit.",
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

func (m Model) ensureTargetBuffer() {
	if len(m.languageWords.Words) == 0 {
		return
	}

	typed := m.currentText.Value()
	remainingRunes := utf8.RuneCountInString(m.Target) - utf8.RuneCountInString(typed)
	if remainingRunes > minRemainingRunes {
		return
	}

	additional := generateWordString(m.rng, m.languageWords, wordBufferChunk)
	if additional == "" {
		return
	}

	if strings.TrimSpace(m.Target) == "" {
		m.Target = additional
	} else {
		m.Target = strings.TrimSpace(m.Target + " " + additional)
	}
	if !m.session.Started() {
		m.currentText.Placeholder = m.Target
	}
}

func (m Model) scheduleTick() tea.Cmd {
	if m.tickInterval <= 0 {
		return nil
	}
	return tea.Tick(m.tickInterval, func(t time.Time) tea.Msg {
		return tickMsg{now: t}
	})
}

func (m Model) displayRemaining() time.Duration {
	if !m.session.Started() {
		return m.totalDuration
	}
	if m.remaining <= 0 {
		return 0
	}
	return m.remaining
}
