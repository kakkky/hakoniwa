package ui

import (
	"context"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kakkky/hakoniwa/presentation/ui/components"
)

func (u *UI) topBar(ctx context.Context) *fyne.Container {
	bdm, err := u.usecases.GetBuildingManager.Exec(ctx)
	if err != nil {

	}

	if bdm != nil {
		return withColor(
			color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			container.NewBorder(nil, nil, widget.NewLabel(fmt.Sprintf("管理人：%s", bdm.Name)), nil),
		)
	}

	return withColor(
		color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		container.NewBorder(
			nil,
			nil,
			components.RegisterBuildingManagerButton(
				func() {
					components.RegisterBuildingManagerFormDialog(
						ctx,
						u.usecases.RegisterBuildingManager,
						u.fyneWindow,
						func() { u.refreshTopBar(ctx) },
					).Show()
				},
			),
			nil,
		),
	)
}

// 再構築
func (u *UI) refreshTopBar(ctx context.Context) {
	u.topBarContent.Objects = []fyne.CanvasObject{u.topBar(ctx)}
	u.topBarContent.Refresh()
}
