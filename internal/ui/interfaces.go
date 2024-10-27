package ui

import tea "github.com/charmbracelet/bubbletea"

type Renderer interface {
	Render() string
}

type Component interface {
	Renderer
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
}
