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

	registerResident *usecase.RegisterResident
	sendMessage      *usecase.SendMessageFromBuildingManagerToResident
}

func NewUI(
	registerResident *usecase.RegisterResident,
	sendMessage *usecase.SendMessageFromBuildingManagerToResident,
) *UI {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(newHakoniwaTheme())
	fyneWindow := fyneApp.NewWindow("hakoniwa")
	fyneWindow.Resize(fyne.NewSize(900, 600))

	return &UI{
		fyneApp:          fyneApp,
		fyneWindow:       fyneWindow,
		registerResident: registerResident,
		sendMessage:      sendMessage,
	}
}

func (u *UI) Run(ctx context.Context) error {
	layout := layoutContainer(
		u.topBar(),
		nil,
		u.residentBoard(),
		u.activityBoard(),
	)

	u.fyneWindow.SetContent(layout)
	u.fyneWindow.ShowAndRun()
	return nil
}
