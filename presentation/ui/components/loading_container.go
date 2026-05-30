package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoadingContainer(loading *widget.Activity) *fyne.Container {
	loading = widget.NewActivity()
	loading.Start()
	return container.NewCenter(
		container.NewHBox(
			loading,
			widget.NewLabel("処理中です..."),
		),
	)
}
