package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fase22/tui/internal/config"
	"github.com/fase22/tui/internal/file"
	"github.com/fase22/tui/internal/ui/components/scrollbar"
	"github.com/fase22/tui/internal/ui/components/statusbar"
	"github.com/fase22/tui/internal/ui/components/textview"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeSearch
)

type Model struct {
	textView    textview.TextView
	statusBar   statusbar.StatusBar
	scrollBar   scrollbar.Scrollbar
	currentFile string
	err         error
	state       string
	config      *config.Config
	mode        Mode
	searchQuery string
	searchIndex int   // Aktueller Treffer-Index
	searchHits  []int // Zeilennummern der Treffer
}

type errMsg struct {
	err error
}

type fileLoadedMsg struct {
	content string
}

type searchMsg struct {
	query string
}

type searchHitMsg struct {
	hits []int
}

func NewModel(filename string, cfg *config.Config) *Model {
	// Styles aus der Konfiguration erstellen
	tvStyle := textview.NewStyleFromConfig(cfg)
	sbStyle := statusbar.NewStyleFromConfig(cfg)
	scrollStyle := scrollbar.NewStyleFromConfig(cfg)

	return &Model{
		textView: textview.New(80, 24, textview.Config{
			ShowLineNumbers: cfg.Editor.ShowLineNumbers,
			TabWidth:        cfg.Editor.TabWidth,
			WordWrap:        cfg.Editor.WordWrap,
			Style:           tvStyle,
		}),
		statusBar:   statusbar.New(filename, 80, sbStyle),
		scrollBar:   scrollbar.New(24, 0, 0, scrollStyle),
		currentFile: filename,
		state:       "initialized",
		config:      cfg,
		mode:        ModeNormal,
	}
}

func (m *Model) Init() tea.Cmd {
	if m.currentFile != "" {
		return m.loadFile
	}
	return nil
}

func (m *Model) loadFile() tea.Msg {
	content, err := file.ReadFile(m.currentFile)
	if err != nil {
		return errMsg{err}
	}
	return fileLoadedMsg{content: content}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Komponenten an neue Fenstergröße anpassen
		m.textView.Resize(msg.Width-2, msg.Height-2) // -2 für Statusleiste

		// Styles aus der Konfiguration erstellen
		sbStyle := statusbar.NewStyleFromConfig(m.config)
		scrollStyle := scrollbar.NewStyleFromConfig(m.config)

		// Komponenten mit Styles aktualisieren
		m.statusBar = statusbar.New(m.currentFile, msg.Width, sbStyle)
		m.scrollBar = scrollbar.New(
			msg.Height-2,
			m.textView.GetTotalLines(),
			m.textView.GetCurrentLine()-1,
			scrollStyle,
		)

	case tea.KeyMsg:
		if m.mode == ModeSearch {
			switch msg.Type {
			case tea.KeyEnter:
				// Suche starten
				m.mode = ModeNormal
				return m, m.search
			case tea.KeyEsc:
				// Suchmodus verlassen
				m.mode = ModeNormal
				m.searchQuery = ""
				m.searchHits = nil
				m.textView.SetSearchTerm("") // Highlighting entfernen
			case tea.KeyBackspace:
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
			default:
				// Ignoriere Steuerungstasten
				if msg.Type != tea.KeyCtrlC && msg.Type != tea.KeyCtrlH {
					// Zeichen zur Suchanfrage hinzufügen
					m.searchQuery += msg.String()
				}
			}
		} else {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				m.textView.ScrollUp(1)
			case "down", "j":
				m.textView.ScrollDown(1)
			case "pgup":
				m.textView.ScrollUp(m.textView.GetViewport().Height)
			case "pgdown":
				m.textView.ScrollDown(m.textView.GetViewport().Height)
			case "/":
				// In den Suchmodus wechseln
				m.mode = ModeSearch
				m.searchQuery = ""
				m.searchHits = nil
				m.searchIndex = 0
			case "n":
				// Zum nächsten Treffer
				if len(m.searchHits) > 0 {
					m.searchIndex = (m.searchIndex + 1) % len(m.searchHits)
					m.jumpToLine(m.searchHits[m.searchIndex])
				}
			case "N":
				// Zum vorherigen Treffer
				if len(m.searchHits) > 0 {
					m.searchIndex--
					if m.searchIndex < 0 {
						m.searchIndex = len(m.searchHits) - 1
					}
					m.jumpToLine(m.searchHits[m.searchIndex])
				}
			}
		}

	case fileLoadedMsg:
		m.textView.SetContent(msg.content)
		m.state = "ready"

	case errMsg:
		m.err = msg.err
		m.state = "error"

	case searchHitMsg:
		m.searchHits = msg.hits
		m.textView.SetSearchTerm(m.searchQuery) // Highlighting aktivieren
		if len(m.searchHits) > 0 {
			m.jumpToLine(m.searchHits[0])
		}
	}

	// Update StatusBar
	m.statusBar.Update(
		m.textView.GetCurrentLine(),
		m.textView.GetTotalLines(),
		len(m.textView.GetViewport().View()),
	)

	// Update search info in status bar
	if m.mode == ModeSearch {
		m.statusBar.SetSearchInfo(
			true,
			m.searchQuery,
			m.searchIndex,
			len(m.searchHits),
		)
	} else {
		m.statusBar.SetSearchInfo(false, "", 0, 0)
	}

	// Update ScrollBar
	m.scrollBar = scrollbar.New(
		m.textView.GetViewport().Height,
		m.textView.GetTotalLines(),
		m.textView.GetCurrentLine()-1,
		scrollbar.NewStyleFromConfig(m.config),
	)

	return m, cmd
}
func (m *Model) View() string {
	if m.err != nil {
		return errorStyle.Render(m.err.Error())
	}

	// Hauptinhalt
	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.textView.Render(),
		m.scrollBar.Render(),
	)

	// Status und Sucheingabe
	var status string
	if m.mode == ModeSearch {
		searchPrompt := fmt.Sprintf("/%s", m.searchQuery)
		if len(m.searchHits) > 0 {
			searchPrompt += fmt.Sprintf(" (%d/%d)", m.searchIndex+1, len(m.searchHits))
		}
		status = searchPrompt
	} else {
		status = m.statusBar.Render()
	}

	// Kombiniere alles
	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		status,
	)
}

func (m *Model) search() tea.Msg {
	if m.searchQuery == "" {
		return searchHitMsg{hits: nil}
	}

	var hits []int
	lines := strings.Split(m.textView.GetContent(), "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(m.searchQuery)) {
			hits = append(hits, i+1)
		}
	}
	return searchHitMsg{hits: hits}
}

func (m *Model) jumpToLine(line int) {
	m.textView.ScrollToLine(line)
}
