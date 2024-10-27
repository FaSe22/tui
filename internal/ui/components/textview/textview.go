package textview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	ShowLineNumbers bool
	TabWidth        int
	WordWrap        bool
	Style           Style
}

type TextView struct {
	viewport    *viewport.Model
	content     string
	width       int
	height      int
	currentLine int
	config      Config
	style       Style
	searchTerm  string
}

func New(width, height int, cfg Config) TextView {
	vp := viewport.New(width, height)
	return TextView{
		viewport: &vp,
		width:    width,
		height:   height,
		config:   cfg,
		style:    cfg.Style,
	}
}

func (tv *TextView) SetContent(content string) {
	tv.content = content

	// Wenn es einen Suchbegriff gibt, wende Highlighting an
	if tv.searchTerm != "" {
		lines := strings.Split(content, "\n")
		var highlightedLines []string

		for _, line := range lines {
			highlightedLines = append(highlightedLines, tv.highlightSearchTerm(line))
		}

		tv.viewport.SetContent(strings.Join(highlightedLines, "\n"))
	} else {
		tv.viewport.SetContent(content)
	}
}

func (tv *TextView) calculateLineNumberWidth() int {
	// Berechne die Anzahl der Stellen der höchsten Zeilennummer
	totalLines := tv.GetTotalLines()
	digits := len(fmt.Sprintf("%d", totalLines))

	// Mindestens 2 Stellen + 1 für Padding
	if digits < 2 {
		digits = 2
	}

	return digits + 1 // +1 für zusätzliches Padding
}

// SetSearchTerm setzt den Suchbegriff und aktualisiert die Anzeige
func (tv *TextView) SetSearchTerm(term string) {
	tv.searchTerm = term
	// Aktualisiere den Inhalt mit dem neuen Suchbegriff
	if tv.content != "" {
		tv.SetContent(tv.content)
	}
}

func (tv *TextView) highlightSearchTerm(line string) string {
	if tv.searchTerm == "" {
		return line
	}

	lowLine := strings.ToLower(line)
	lowTerm := strings.ToLower(tv.searchTerm)

	var result string
	lastIdx := 0

	// Finde alle Vorkommen und highlighte sie
	for {
		idx := strings.Index(lowLine[lastIdx:], lowTerm)
		if idx == -1 {
			result += line[lastIdx:]
			break
		}

		idx += lastIdx
		match := line[idx : idx+len(tv.searchTerm)]
		result += line[lastIdx:idx] + tv.style.SearchMatch.Render(match)
		lastIdx = idx + len(tv.searchTerm)
	}

	return result
}

func (tv *TextView) Render() string {
	if tv.content == "" {
		return tv.style.EmptyText.Render("Keine Datei geladen")
	}

	lines := strings.Split(tv.viewport.View(), "\n")
	startLine := tv.viewport.YOffset + 1

	var contentBuilder strings.Builder
	maxWidth := tv.width

	lineNumWidth := tv.calculateLineNumberWidth()

	if tv.config.ShowLineNumbers {
		tv.style.LineNumber = tv.style.LineNumber.Width(lineNumWidth)
		maxWidth -= lineNumWidth + 1
	}

	for i, line := range lines {
		lineNum := startLine + i

		var linePrefix string
		if tv.config.ShowLineNumbers {
			format := fmt.Sprintf("%%%dd", lineNumWidth-1)
			linePrefix = tv.style.LineNumber.Render(fmt.Sprintf(format, lineNum))
		}

		// Textzeile mit Highlighting
		lineContent := line
		if len(lineContent) > maxWidth {
			lineContent = lineContent[:maxWidth-3] + "..."
		}

		// Suchbegriff highlighten
		lineContent = tv.highlightSearchTerm(lineContent)

		// Aktuelle Zeile hervorheben
		if i == (tv.viewport.YPosition) {
			lineContent = tv.style.CurrentLine.Render(lineContent)
		}

		var fullLine string
		if tv.config.ShowLineNumbers {
			fullLine = lipgloss.JoinHorizontal(
				lipgloss.Left,
				linePrefix,
				" ",
				lineContent,
			)
		} else {
			fullLine = lineContent
		}

		contentBuilder.WriteString(fullLine + "\n")
	}

	return tv.style.Container.Render(contentBuilder.String())
}

// Hilfsfunktion für Wortumbruch
func (tv *TextView) wrapText(text string) string {
	if !tv.config.WordWrap {
		return text
	}

	width := tv.width
	if tv.config.ShowLineNumbers {
		width -= tv.style.LineNumber.GetWidth()
	}

	var wrapped strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if len(line) <= width {
			wrapped.WriteString(line)
		} else {
			// Einfacher Wortumbruch
			words := strings.Fields(line)
			lineLen := 0
			for _, word := range words {
				wordLen := len(word) + 1 // +1 für Leerzeichen
				if lineLen+wordLen > width {
					wrapped.WriteString("\n")
					lineLen = 0
				}
				if lineLen > 0 {
					wrapped.WriteString(" ")
				}
				wrapped.WriteString(word)
				lineLen += wordLen
			}
		}
		if i < len(lines)-1 {
			wrapped.WriteString("\n")
		}
	}

	return wrapped.String()
}

// Standard Getter/Setter Methoden bleiben gleich
func (tv *TextView) ScrollUp(lines int) {
	tv.viewport.LineUp(lines)
	tv.currentLine = tv.viewport.YOffset
}

func (tv *TextView) ScrollDown(lines int) {
	tv.viewport.LineDown(lines)
	tv.currentLine = tv.viewport.YOffset
}

func (tv *TextView) GetViewport() *viewport.Model {
	return tv.viewport
}

func (tv *TextView) GetCurrentLine() int {
	return tv.currentLine + 1
}

func (tv *TextView) GetTotalLines() int {
	return strings.Count(tv.content, "\n") + 1
}

func (tv *TextView) ToggleLineNumbers() {
	tv.config.ShowLineNumbers = !tv.config.ShowLineNumbers
}

func (tv *TextView) Resize(width, height int) {
	tv.width = width
	tv.height = height
	tv.viewport.Width = width
	tv.viewport.Height = height

	// Bei Größenänderung Wortumbruch neu anwenden
	if tv.config.WordWrap {
		tv.SetContent(tv.content)
	}
}

func (tv *TextView) ToggleWordWrap() {
	tv.config.WordWrap = !tv.config.WordWrap
	if tv.content != "" {
		tv.SetContent(tv.content) // Neu rendern mit/ohne Wrap
	}
}
func (tv *TextView) GetContent() string {
	return tv.content
}

func (tv *TextView) ScrollToLine(line int) {
	// Berücksichtige, dass Zeilennummern bei 1 beginnen
	targetLine := line - 1
	if targetLine < 0 {
		targetLine = 0
	}

	// Berechne die optimale Scrollposition
	viewportHeight := tv.viewport.Height
	halfHeight := viewportHeight / 2

	// Zentriere die Zeile im Viewport wenn möglich
	newPosition := targetLine - halfHeight
	if newPosition < 0 {
		newPosition = 0
	}

	tv.viewport.YOffset = newPosition
	tv.currentLine = targetLine
}
