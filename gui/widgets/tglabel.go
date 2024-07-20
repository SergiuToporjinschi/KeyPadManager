package widgets

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TGLabel struct {
	widget.Label
	OnDragged         onDragFunc
	OnDragEnd         onDragEndFunc
	OnTapped          onTappedFunc
	OnTappedSecondary onTappedFunc
}

func NewTGLabel(text string) *TGLabel {
	lable := &TGLabel{
		// Label: widget.Label{
		// 	Text: text,
		// },
	}
	lable.ExtendBaseWidget(lable)
	lable.SetText(text)
	return lable
}

func (t *TGLabel) Tapped(event *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped(event)
	}
	slog.Info("Tapped")
}

func (t *TGLabel) TappedSecondary(event *fyne.PointEvent) {
	if t.OnTappedSecondary != nil {
		t.OnTappedSecondary(event)
	}
}

func (t *TGLabel) DragEnd() {
	if t.OnDragEnd != nil {
		t.OnDragEnd()
	}
}

func (t *TGLabel) Dragged(dr *fyne.DragEvent) {
	if t.OnDragged != nil {
		t.OnDragged(dr)
	}
}

func (t *TGLabel) FocusGained() {
	slog.Info("Focus gained")
}

func (t *TGLabel) FocusLost() {
	slog.Info("Focus lost")
}
