package quote_input

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	typedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	remainingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Model struct {
	// Target text
	Target string
	// what user has currentText so far
	currentText textarea.Model
	// timing
	started  bool
	start    time.Time
	finished bool
	end      time.Time
	wpm      float64
	err      error
}

func InitialModel(target string) Model {
	ti := textarea.New()
	ti.Placeholder = target
	ti.Focus()

	return Model{
		Target:      target,
		currentText: ti,
		err:         nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages (key presses, etc.)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.finished {
			return m, tea.Quit
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
	cmds = append(cmds, cmd)
	cmd = tea.Batch(cmds...)

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

	b.WriteString(renderBox(typed, remaining, len(m.Target)) + "\n\n")

	if m.finished {
		b.WriteString(fmt.Sprintf("✅ Done! WPM: %.2f\n", m.wpm))
		b.WriteString("Press any key to exit.\n")
	}

	return b.String()
}

func renderBox(typed string, remaining string, width int) string {
	top := "+" + strings.Repeat("-", width) + "+\n"
	complete := typedStyle.Render(typed) + remainingStyle.Render(remaining)
	if len(complete) > width {
		complete = complete[:width]
	} else if len(typed) < width {
		complete += strings.Repeat(" ", width-len(complete))
	}
	middle := "|" + complete + "|\n"
	bottom := "+" + strings.Repeat("-", width) + "+\n"
	return top + middle + bottom
}
