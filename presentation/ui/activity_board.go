package ui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (u *UI) activityBoard() *fyne.Container {
	activityCards := activityMock()

	return withColor(
		color.NRGBA{R: 222, G: 222, B: 222, A: 255},
		container.NewVBox(
			activityCards...,
		),
	)
}

func activityMock() []fyne.CanvasObject {
	var cards []fyne.CanvasObject
	activities := []struct {
		with    []string
		content string
	}{
		{
			with:    []string{"田中", "佐藤"},
			content: "会話中",
		},
		{
			with:    []string{"岡田", "田町"},
			content: "行動中",
		},
	}
	for _, a := range activities {
		cards = append(
			cards,
			container.NewPadded(
				withColor(
					color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					widget.NewCard("", "", container.NewHBox(widget.NewLabel(strings.Join(a.with, "と")), widget.NewLabel(string(a.content)))),
				),
			),
		)
	}

	return cards
}
