package words_input

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/theme"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/typing"
)

type Model struct {
	// Target text
	Target string
	// what user has currentText so far
	currentText   textarea.Model
	rng           *rand.Rand
	viewportWidth int
	styles        theme.Styles
	session       typing.Session
}

func InitialModel(languageWords models.LanguageWords) Model {
	tempWords := languageWords.Words
	if len(tempWords) > 50 {
		tempWords = tempWords[:50]
	}

	tempTextToType := ""
	for i, word := range tempWords {
		if i != 0 {
			tempTextToType += " "
		}
		tempTextToType += word
	}

	ti := textarea.New()
	ti.Placeholder = tempTextToType
	ti.SetWidth(typing.DefaultBoxWidth)
	ti.Focus()

	return Model{
		Target:      tempTextToType,
		currentText: ti,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		styles:      theme.DefaultStyles(),
		session:     typing.NewSession(),
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

// View defines UI rendering
func (m Model) View() string {
	return m.currentText.View()
}
