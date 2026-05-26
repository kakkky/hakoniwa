package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kakkky/hakoniwa/domain"
)

func (u *UI) residentBoard() *fyne.Container {
	residentCards := residentCardsMock()

	return withColor(
		color.NRGBA{R: 232, G: 232, B: 232, A: 255},
		container.NewBorder(
			nil, nil, nil, nil,
			container.NewGridWrap(
				fyne.NewSize(260, 70),
				residentCards...,
			),
		))
}

func residentCardsMock() []fyne.CanvasObject {
	var cards []fyne.CanvasObject
	residents := []domain.Resident{
		{Name: "田中", Mood: "😅"},
		{Name: "佐藤", Mood: "😀"},
		{Name: "岡田", Mood: "😡"},
		{Name: "木下", Mood: "🤔"},
	}
	for _, r := range residents {
		cards = append(
			cards,
			container.NewPadded(
				withColor(
					color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					widget.NewCard("", "", container.NewHBox(widget.NewLabel(string(r.Name)), widget.NewLabel(string(r.Mood)))),
				),
			),
		)
	}

	return cards
}
