package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Farbpalette definieren
var (
    colors = struct {
        background string
        foreground string
        accent     string
        highlight  string
        subtle     string
        error      string
    }{
        background: "#282c34",
        foreground: "#abb2bf",
        accent:     "#61afef",
        highlight:  "#353b45",
        subtle:     "#4b5263",
        error:      "#e06c75",
    }

    // UI Styles
    statusStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(colors.foreground)).
        Background(lipgloss.Color(colors.subtle)).
        Bold(true).
        Padding(0, 1).
        Width(100).
        MarginBottom(1)

    textStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(colors.foreground)).
        Background(lipgloss.Color(colors.background)).
        Padding(0, 1)

    lineNumberStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(colors.subtle)).
        Background(lipgloss.Color(colors.background)).
        Padding(0, 1).
        Width(4).
        Align(lipgloss.Right)

    currentLineStyle = lipgloss.NewStyle().
        Background(lipgloss.Color(colors.highlight)).
        Foreground(lipgloss.Color(colors.accent)).
        Bold(true)

    scrollbarStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(colors.subtle)).
        Background(lipgloss.Color(colors.background))

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(colors.error)).
        Bold(true)

    // Container-Style für den gesamten Bereich
    mainStyle = lipgloss.NewStyle().
        Background(lipgloss.Color(colors.background)).
        Margin(1).
        Padding(0)
)
