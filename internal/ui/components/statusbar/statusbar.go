package statusbar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	fileName      string
	viewportWidth int
	currentLine   int
	totalLines    int
	fileSize      int
	style         Style
	searchMode    bool
	searchQuery   string
	searchResults string
}

func New(filename string, viewportWidth int, style Style) StatusBar {
	return StatusBar{
		fileName:      filename,
		viewportWidth: viewportWidth,
		style:         style,
	}
}

func (s *StatusBar) Update(currentLine, totalLines, fileSize int) {
	s.currentLine = currentLine
	s.totalLines = totalLines
	s.fileSize = fileSize
}

func (s *StatusBar) SetSearchInfo(active bool, query string, current, total int) {
	s.searchMode = active
	s.searchQuery = query
	if total > 0 {
		s.searchResults = fmt.Sprintf("%d/%d", current+1, total)
	} else {
		s.searchResults = "No matches"
	}
}

func (s StatusBar) Render() string {
	if s.searchMode {
		searchPrompt := fmt.Sprintf("/%s", s.searchQuery)
		if s.searchResults != "" {
			searchPrompt += fmt.Sprintf(" (%s)", s.searchResults)
		}
		return s.style.SearchMode.Render(searchPrompt)
	}

	// Berechne Prozentposition
	percentage := 0
	if s.totalLines > 0 {
		percentage = int((float64(s.currentLine) / float64(s.totalLines)) * 100)
	}

	// Linke Seite
	leftStatus := fmt.Sprintf(
		"%s - %s",
		s.fileName,
		s.formatFileSize(),
	)

	// Mittlerer Teil (Shortcuts)
	middleStatus := "NORMAL | ^F: Suche | ^S: Speichern | q: Beenden"

	// Rechte Seite
	rightStatus := fmt.Sprintf(
		"Zeile %d/%d [%d%%]",
		s.currentLine,
		s.totalLines,
		percentage,
	)

	// Layout berechnen
	leftWidth := lipgloss.Width(leftStatus)
	rightWidth := lipgloss.Width(rightStatus)
	middleWidth := s.viewportWidth - leftWidth - rightWidth - 2

	// Zentrieren des mittleren Teils
	middleStatus = s.centerText(middleStatus, middleWidth)

	return s.style.Base.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			s.style.LeftSection.Render(leftStatus),
			s.style.MiddleSection.Render(middleStatus),
			s.style.RightSection.Render(rightStatus),
		),
	)
}

func (s StatusBar) formatFileSize() string {
	size := float64(s.fileSize)
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", int(size))
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", size/1024)
	default:
		return fmt.Sprintf("%.1f MB", size/(1024*1024))
	}
}

func (s StatusBar) centerText(text string, width int) string {
	if width <= 0 {
		return ""
	}
	textWidth := lipgloss.Width(text)
	if textWidth >= width {
		return text[:width]
	}
	leftPad := (width - textWidth) / 2
	rightPad := width - textWidth - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}
