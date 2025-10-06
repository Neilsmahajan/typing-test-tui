package quote_input

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container     lipgloss.Style
	Header        lipgloss.Style
	Subtitle      lipgloss.Style
	QuoteBox      lipgloss.Style
	QuoteContent  lipgloss.Style
	Instruction   lipgloss.Style
	StatsRow      lipgloss.Style
	StatBlock     lipgloss.Style
	StatLabel     lipgloss.Style
	StatValue     lipgloss.Style
	StatSeparator string
	Success       lipgloss.Style
	Typed         lipgloss.Style
	Incorrect     lipgloss.Style
	Remaining     lipgloss.Style
	Cursor        lipgloss.Style
}

func defaultStyles() Styles {
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("60")).Padding(0, 1).Render("│")

	return Styles{
		Container:     lipgloss.NewStyle().Align(lipgloss.Left),
		Header:        lipgloss.NewStyle().Foreground(lipgloss.Color("218")).Bold(true),
		Subtitle:      lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		QuoteBox:      lipgloss.NewStyle().MarginTop(1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(1, 2),
		QuoteContent:  lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Instruction:   lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Faint(true).Italic(true).MarginTop(1),
		StatsRow:      lipgloss.NewStyle().MarginTop(1),
		StatBlock:     lipgloss.NewStyle().Padding(0, 2, 0, 0),
		StatLabel:     lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Bold(true),
		StatValue:     lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		StatSeparator: separator,
		Success:       lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true).MarginTop(1),
		Typed:         lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
		Incorrect:     lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Underline(true),
		Remaining:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Cursor:        lipgloss.NewStyle().Background(lipgloss.Color("218")).Foreground(lipgloss.Color("0")),
	}
}
