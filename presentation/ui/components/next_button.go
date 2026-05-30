package components

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewNavigateNextButton(text string, navigate func()) *widget.Button {
	btn := widget.NewButtonWithIcon(text, theme.NavigateNextIcon(), navigate)
	btn.IconPlacement = widget.ButtonIconTrailingText
	return btn
}
