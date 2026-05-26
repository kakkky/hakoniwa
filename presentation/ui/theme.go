package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type hakoniwaTheme struct {
	fyne.Theme
}

func newHakoniwaTheme() *hakoniwaTheme {
	return &hakoniwaTheme{Theme: theme.LightTheme()}
}

func (h *hakoniwaTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 42, G: 42, B: 42, A: 255}
	}
	return h.Theme.Color(name, variant)
}
