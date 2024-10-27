package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ThemeName definiert die verfügbaren Themes
type ThemeName string

const (
	ThemeDark    ThemeName = "dark"
	ThemeLight   ThemeName = "light"
	ThemeDracula ThemeName = "dracula"
)

// Theme enthält die Farbkonfiguration
type Theme struct {
	Name        ThemeName `json:"name"`
	Background  string    `json:"background"`
	Foreground  string    `json:"foreground"`
	Selection   string    `json:"selection"`
	Accent      string    `json:"accent"`
	LineNumbers string    `json:"lineNumbers"`
}

// Vordefinierte Themes
var (
	DarkTheme = Theme{
		Name:        ThemeDark,
		Background:  "#282c34",
		Foreground:  "#abb2bf",
		Selection:   "#3e4451",
		Accent:      "#61afef",
		LineNumbers: "#4b5263",
	}

	LightTheme = Theme{
		Name:        ThemeLight,
		Background:  "#ffffff",
		Foreground:  "#383a42",
		Selection:   "#e5e5e6",
		Accent:      "#4078f2",
		LineNumbers: "#9d9d9f",
	}

	DraculaTheme = Theme{
		Name:        ThemeDracula,
		Background:  "#282a36",
		Foreground:  "#f8f8f2",
		Selection:   "#44475a",
		Accent:      "#bd93f9",
		LineNumbers: "#6272a4",
	}
)

type Config struct {
	// Theme-Einstellungen
	Theme Theme `json:"theme"`

	// Editor-Einstellungen
	Editor struct {
		ShowLineNumbers bool `json:"showLineNumbers"`
		TabWidth        int  `json:"tabWidth"`
		WordWrap        bool `json:"wordWrap"`
		AutoIndent      bool `json:"autoIndent"`
	} `json:"editor"`

	// UI-Einstellungen
	UI struct {
		ShowScrollbar bool   `json:"showScrollbar"`
		ShowStatus    bool   `json:"showStatus"`
		ScrollStyle   string `json:"scrollStyle"` // "bar" oder "block"
	} `json:"ui"`

	// Tastatur-Shortcuts
	Keybindings struct {
		QuitKey        string `json:"quitKey"`
		SaveKey        string `json:"saveKey"`
		ToggleWrapKey  string `json:"toggleWrapKey"`
		ToggleLinesKey string `json:"toggleLinesKey"`
	} `json:"keybindings"`
}

// DefaultConfig erstellt eine Standardkonfiguration
func DefaultConfig() Config {
	cfg := Config{}

	// Standard-Theme (Dark)
	cfg.Theme = DarkTheme

	// Standard Editor-Einstellungen
	cfg.Editor.ShowLineNumbers = true
	cfg.Editor.TabWidth = 4
	cfg.Editor.WordWrap = false
	cfg.Editor.AutoIndent = true

	// Standard UI-Einstellungen
	cfg.UI.ShowScrollbar = true
	cfg.UI.ShowStatus = true
	cfg.UI.ScrollStyle = "bar"

	// Standard Keybindings
	cfg.Keybindings.QuitKey = "q"
	cfg.Keybindings.SaveKey = "ctrl+s"
	cfg.Keybindings.ToggleWrapKey = "ctrl+w"
	cfg.Keybindings.ToggleLinesKey = "ctrl+l"

	return cfg
}

// GetTheme gibt das konfigurierte Theme zurück
func (c Config) GetTheme() Theme {
	switch c.Theme.Name {
	case ThemeLight:
		return LightTheme
	case ThemeDracula:
		return DraculaTheme
	default:
		return DarkTheme
	}
}

// LoadConfig lädt die Konfiguration aus einer Datei
func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		// Zuerst lokale config.json im Projektverzeichnis suchen
		if _, err := os.Stat("config.json"); err == nil {
			path = "config.json"
		} else {
			// Falls nicht gefunden, nutze den Home-Directory-Pfad
			home, err := os.UserHomeDir()
			if err != nil {
				return cfg, err
			}
			path = filepath.Join(home, ".config", "tui", "config.json")
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, SaveConfig(path, cfg)
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// SaveConfig speichert die Konfiguration in einer Datei
func SaveConfig(path string, cfg Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
