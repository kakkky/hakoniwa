package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (u *UI) menuContainer(menuItems ...*fyne.MenuItem) *widget.Button {
	var btn *widget.Button
	btn = u.menuButton(func() {
		popup := u.menuPopupContent(menuItems...)
		btnPos := u.fyneApp.Driver().AbsolutePositionForObject(btn)
		popup.ShowAtPosition(fyne.NewPos(btnPos.X, btnPos.Y+btn.Size().Height))
	})
	return btn
}

func (u *UI) menuButton(onClick func()) *widget.Button {
	btn := widget.NewButtonWithIcon("", theme.MenuIcon(), onClick)
	btn.Importance = widget.LowImportance
	return btn
}

func (u *UI) menuPopupContent(menuItems ...*fyne.MenuItem) *widget.PopUpMenu {
	menu := widget.NewPopUpMenu(
		fyne.NewMenu("", menuItems...),
		u.fyneMainWindow.Canvas(),
	)
	return menu
}
