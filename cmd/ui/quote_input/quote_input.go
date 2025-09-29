package quote_input

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// Target text
	Target string
	// what user has typed so far
	typed textinput.Model
	// timing
	started  bool
	start    time.Time
	finished bool
	end      time.Time
	wpm      float64
	err      error
}

func InitialModel(target string) Model {
	ti := textinput.New()
	ti.Placeholder = target
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 156

	return Model{
		Target: target,
		typed:  ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages (key presses, etc.)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.finished {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		default:
			// start timer on first key
			if !m.started && len(msg.String()) == 1 {
				m.started = true
				m.start = time.Now()
			}
			m.typed, cmd = m.typed.Update(msg)
		}

	case error:
		m.err = msg
		return m, nil
	}

	// check if completed (capture finish time & wpm only once)
	if !m.finished && m.typed.Value() == m.Target {
		m.finished = true
		m.end = time.Now()
		elapsedMinutes := m.end.Sub(m.start).Minutes()
		if elapsedMinutes > 0 {
			m.wpm = float64(len(strings.Fields(m.Target))) / elapsedMinutes
		}
	}

	return m, cmd
}

// View defines UI rendering
func (m Model) View() string {
	var b strings.Builder

	b.WriteString("\nType the following:\n\n")
	b.WriteString(m.Target + "\n\n")

	// highlight typed portion
	b.WriteString("You typed: " + m.typed.View() + "\n\n")

	if m.finished {
		b.WriteString(fmt.Sprintf("✅ Done! WPM: %.2f\n", m.wpm))
		b.WriteString("Press any key to exit.\n")
	}

	return b.String()
}
