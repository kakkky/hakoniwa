package components

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func RegisterBuildingManagerButton(onClick func()) *widget.Button {
	return widget.NewButtonWithIcon(
		"新規登録",
		theme.ContentAddIcon(),
		onClick,
	)
}
