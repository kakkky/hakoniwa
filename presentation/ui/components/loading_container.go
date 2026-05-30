package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoadingContainer() *fyne.Container {
	loading := widget.NewActivity()
	loading.Start()
	return container.NewCenter(
		container.NewHBox(
			container.NewGridWrap(fyne.NewSize(32, 32), loading),
			widget.NewLabel("処理中です..."),
		),
	)
}
