package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func layoutContainer(
	topBar *fyne.Container,
	sideBar *fyne.Container,
	residentBord *fyne.Container,
	activityBord *fyne.Container,
) *fyne.Container {
	mainContent := container.NewVSplit(
		residentBord,
		activityBord,
	)
	mainContent.Offset = 0.6

	return container.NewBorder(
		topBar,
		nil,
		nil,
		nil,
		mainContent,
	)
}
