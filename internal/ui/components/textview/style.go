package textview

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/fase22/tui/internal/config"
)

type Style struct {
	Container   lipgloss.Style
	LineNumber  lipgloss.Style
	NormalLine  lipgloss.Style
	CurrentLine lipgloss.Style
	EmptyText   lipgloss.Style
	SearchMatch lipgloss.Style
}

func NewStyleFromConfig(cfg *config.Config) Style {
	theme := cfg.Theme

	return Style{
		Container: lipgloss.NewStyle().
			Padding(0),

		LineNumber: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.LineNumbers)).
			Background(lipgloss.Color(theme.Background)).
			Padding(0, 1).
			Width(4).
			Align(lipgloss.Right),

		NormalLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground)).
			Background(lipgloss.Color(theme.Background)),

		CurrentLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground)).
			Background(lipgloss.Color(theme.Selection)).
			Bold(true),

		EmptyText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.LineNumbers)).
			Align(lipgloss.Center),
		SearchMatch: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Accent)).
			Foreground(lipgloss.Color(theme.Background)).
			Bold(true),
	}
}
