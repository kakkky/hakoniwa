package ui

import (
	"context"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func (u *UI) topBar(ctx context.Context) *fyne.Container {
	return withColor(
		color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		container.NewBorder(
			nil,
			nil,
			u.buildingManagerInfomationContainer(ctx),
			u.menuContainer(
				fyne.NewMenuItem("住人登録", func() { u.registerResidentContainer(ctx).Show() }),
			),
		),
	)
}

func (u *UI) refreshTopBar(ctx context.Context) {
	u.topBarContent.Objects = []fyne.CanvasObject{u.topBar(ctx)}
	u.topBarContent.Refresh()
}
