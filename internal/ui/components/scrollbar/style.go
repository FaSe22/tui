package scrollbar

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/fase22/tui/internal/config"
)

type Style struct {
    Track   lipgloss.Style
    Thumb   lipgloss.Style
    Symbols ScrollbarSymbols
}

type ScrollbarSymbols struct {
    Single string
    Top    string
    Bottom string
    Body   string
    Track  string
}

func NewStyleFromConfig(cfg *config.Config) Style {
    theme := cfg.Theme

    return Style{
        Track: lipgloss.NewStyle().
            Foreground(lipgloss.Color(theme.LineNumbers)).
            Background(lipgloss.Color(theme.Background)),

        Thumb: lipgloss.NewStyle().
            Foreground(lipgloss.Color(theme.Selection)).
            Background(lipgloss.Color(theme.Background)),

        Symbols: ScrollbarSymbols{
            Single: "█",
            Top:    "▀",
            Bottom: "▄",
            Body:   "█",
            Track:  "│",
        },
    }
}
