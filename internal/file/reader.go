package file

import (
	"fmt"
	"io"
	"os"
)

// ReadFile liest eine Textdatei und gibt deren Inhalt zurück
func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Fehler beim Öffnen der Datei: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Fehler beim Lesen der Datei: %w", err)
	}

	return string(content), nil
}
