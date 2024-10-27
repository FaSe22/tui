package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fase22/tui/internal/config"
	"github.com/fase22/tui/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Bitte geben Sie einen Dateinamen an")
		os.Exit(1)
	}

	// Lade Konfiguration
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Printf("Warnung: Konnte Konfiguration nicht laden: %v\n", err)
		cfg = config.DefaultConfig()
	}

	filename := os.Args[1]
	p := tea.NewProgram(ui.NewModel(filename, &cfg))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Ahhh, es gab einen Fehler: %v", err)
		os.Exit(1)
	}
}
