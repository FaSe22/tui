package scrollbar

import (
	"strings"
)

type Scrollbar struct {
	height        int
	contentHeight int
	offset        int
	style         Style
}

func New(height int, contentHeight int, offset int, style Style) Scrollbar {
	return Scrollbar{
		height:        height,
		contentHeight: contentHeight,
		offset:        offset,
		style:         style,
	}
}

func (s Scrollbar) Render() string {
	if s.contentHeight <= s.height {
		return strings.Repeat(s.style.Track.Render("│"), s.height)
	}

	// Scrollbar-Komponenten aus dem Style
	symbols := s.style.Symbols

	ratio := float64(s.height) / float64(s.contentHeight)
	thumbSize := int(float64(s.height) * ratio)
	if thumbSize < 1 {
		thumbSize = 1
	}

	thumbPosition := int(float64(s.offset) / float64(s.contentHeight) * float64(s.height))

	var scrollbar strings.Builder

	for i := 0; i < s.height; i++ {
		switch {
		case i == thumbPosition && thumbSize == 1:
			scrollbar.WriteString(s.style.Thumb.Render(symbols.Single))
		case i == thumbPosition:
			scrollbar.WriteString(s.style.Thumb.Render(symbols.Top))
		case i == thumbPosition+thumbSize-1:
			scrollbar.WriteString(s.style.Thumb.Render(symbols.Bottom))
		case i > thumbPosition && i < thumbPosition+thumbSize-1:
			scrollbar.WriteString(s.style.Thumb.Render(symbols.Body))
		default:
			scrollbar.WriteString(s.style.Track.Render(symbols.Track))
		}
		scrollbar.WriteString("\n")
	}

	return scrollbar.String()
}
