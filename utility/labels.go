package utility

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func NewTitleLabel(text string) *widget.Label {
	lbl := widget.NewLabel(text)
	lbl.TextStyle.Bold = true
	return lbl
}

func NewSizeableText(text string, size float32) *canvas.Text {
	txt := NewSizeableColorText(text, size, nil)
	txt.TextSize = size
	return txt
}

func NewSizeableColorText(text string, size float32, clr color.Color) *canvas.Text {
	color := clr

	if clr == nil && fyne.CurrentApp() != nil { // nil app possible if app not started
		variant := fyne.CurrentApp().Settings().ThemeVariant()
		color = fyne.CurrentApp().Settings().Theme().Color("text", variant) // manually name the size to avoid import loop
	}
	txt := &canvas.Text{
		Color:    color,
		Text:     text,
		TextSize: size,
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
	}
	return txt
}
