package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/kakkky/hakoniwa/usecase"
)

type UI struct {
	fyneApp    fyne.App
	fyneWindow fyne.Window

	topBarContent        *fyne.Container
	residentBoardContent *fyne.Container
	activityBoardContent *fyne.Container

	usecases *usecase.Usecases
}

func NewUI(usecases *usecase.Usecases) *UI {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(newHakoniwaTheme())
	fyneWindow := fyneApp.NewWindow("hakoniwa")
	fyneWindow.Resize(fyne.NewSize(900, 600))

	return &UI{
		fyneApp:    fyneApp,
		fyneWindow: fyneWindow,
		usecases:   usecases,
	}
}

func (u *UI) Run(ctx context.Context) error {
	u.topBarContent = u.topBar(ctx)
	u.residentBoardContent = u.residentBoard()
	u.activityBoardContent = u.activityBoard()
	layout := layoutContainer(
		u.topBarContent,
		nil,
		u.residentBoardContent,
		u.activityBoardContent,
	)

	u.fyneWindow.SetContent(layout)
	u.fyneWindow.ShowAndRun()
	return nil
}
