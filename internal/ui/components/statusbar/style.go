package statusbar

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/fase22/tui/internal/config"
)

type Style struct {
	Base          lipgloss.Style
	LeftSection   lipgloss.Style
	MiddleSection lipgloss.Style
	RightSection  lipgloss.Style
	SearchMode    lipgloss.Style
}

func NewStyleFromConfig(cfg *config.Config) Style {
	theme := cfg.Theme

	return Style{
		Base: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Selection)).
			Padding(0, 1),

		LeftSection: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground)).
			Bold(true),

		MiddleSection: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.LineNumbers)),

		RightSection: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground)).
			Bold(true),
		SearchMode: lipgloss.NewStyle(). // Neu
							Foreground(lipgloss.Color(theme.Accent)).
							Background(lipgloss.Color(theme.Selection)).
							Bold(true),
	}
}
