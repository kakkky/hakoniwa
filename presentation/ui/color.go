package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func withColor(color color.Color, content fyne.CanvasObject) *fyne.Container {
	rectangle := canvas.NewRectangle(color)
	return container.NewStack(rectangle, content)
}
