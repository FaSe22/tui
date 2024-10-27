# TUI Text Viewer

Ein Terminal-basierter Textbetrachter mit Such- und Navigationsfunktionen.

## Features
- Dateiansicht mit Zeilennummern
- Suchfunktion mit Highlighting
- Konfigurierbare Themes
- Scrollbar
- Statusleiste

## Installation
```bash
go install github.com/fase22/tui/cmd/reader@latest
```

## Verwendung
```bash
reader [filename]
```

## Tastenkombinationen
- `q` oder `Ctrl+C`: Beenden
- `↑` oder `k`: Eine Zeile nach oben
- `↓` oder `j`: Eine Zeile nach unten
- `PgUp`: Seitenweise nach oben
- `PgDn`: Seitenweise nach unten
- `/`: Suchmoduls aktivieren
- `n`: Zum nächsten Suchergebnis
- `N`: Zum vorherigen Suchergebnis
- `ESC`: Suchmodus verlassen

## Konfiguration
Die Konfiguration erfolgt über eine `config.json` Datei, die entweder im aktuellen Verzeichnis oder unter `~/.config/tui/config.json` liegt.

Beispiel-Konfiguration:
```json
{
    "theme": {
        "name": "dracula",
        "background": "#282a36",
        "foreground": "#f8f8f2",
        "selection": "#44475a",
        "accent": "#bd93f9",
        "lineNumbers": "#6272a4"
    },
    "editor": {
        "showLineNumbers": true,
        "tabWidth": 4,
        "wordWrap": true,
        "autoIndent": true
    },
    "ui": {
        "showScrollbar": true,
        "showStatus": true,
        "scrollStyle": "bar"
    }
}
```
