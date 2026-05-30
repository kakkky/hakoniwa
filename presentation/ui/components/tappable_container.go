package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TappableContainer struct {
	widget.BaseWidget
	content  fyne.CanvasObject
	OnTapped func()
}

func NewTappableContainer(content fyne.CanvasObject, tapped func()) *TappableContainer {
	t := &TappableContainer{content: content, OnTapped: tapped}
	t.ExtendBaseWidget(t)
	return t
}

func (t *TappableContainer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.content)
}

// fyne.Tappable
func (t *TappableContainer) Tapped(_ *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped()
	}
}
