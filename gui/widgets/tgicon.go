package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type onDragFunc func(*fyne.DragEvent)
type onDragEndFunc func()
type onTappedFunc func(*fyne.PointEvent)

type TGIcon struct {
	widget.Icon
	OnDragged         onDragFunc
	OnDragEnd         onDragEndFunc
	OnTapped          onTappedFunc
	OnTappedSecondary onTappedFunc
}

func NewTGIcon(res fyne.Resource) *TGIcon {
	icon := &TGIcon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}

func (t *TGIcon) Tapped(event *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped(event)
	}
}

func (t *TGIcon) TappedSecondary(event *fyne.PointEvent) {
	if t.OnTappedSecondary != nil {
		t.OnTappedSecondary(event)
	}
}

func (t *TGIcon) DragEnd() {
	if t.OnDragEnd != nil {
		t.OnDragEnd()
	}
}

func (t *TGIcon) Dragged(dr *fyne.DragEvent) {
	if t.OnDragged != nil {
		t.OnDragged(dr)
	}
}
