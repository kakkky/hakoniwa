package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewEditableField(label, originalValue string, onSubmitted func(v string)) *fyne.Container {
	entry := widget.NewEntry()
	entry.SetText(originalValue)
	valueLabel := widget.NewLabel(originalValue)
	editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

	setEditing := func(editing bool) {
		if editing {
			valueLabel.Hide()
			entry.Show()
			editBtn.Hide()
		} else {
			entry.Hide()
			valueLabel.Show()
			editBtn.Show()
		}
	}
	setEditing(false)

	editBtn.OnTapped = func() { setEditing(true) }
	entry.OnSubmitted = func(v string) {
		valueLabel.SetText(v)
		onSubmitted(v)
		setEditing(false)
	}

	return container.NewBorder(
		nil, nil,
		widget.NewLabel(fmt.Sprintf("%s :", label)),
		editBtn,
		container.NewStack(valueLabel, entry),
	)
}
