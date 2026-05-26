package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (u *UI) topBar() *fyne.Container {
	buildingManagerInfomation := widget.NewLabel("管理人: hogeさん")

	return withColor(
		color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		container.NewBorder(nil, nil, buildingManagerInfomation, nil),
	)
}
