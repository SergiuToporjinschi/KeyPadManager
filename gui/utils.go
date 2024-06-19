package gui

import (
	"image/color"
	"main/txt"
	"main/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type style struct {
	// size  float32
	align fyne.TextAlign
	wrap  fyne.TextWrap
	trun  fyne.TextTruncation
	imp   widget.Importance

	Bold      bool // Should text be bold
	Italic    bool // Should text be italic
	Monospace bool // Use the system monospace font instead of regular
	// Since: 2.2
	Symbol bool // Use the system symbol font.
	// Since: 2.1
	TabWidth int // Width of tabs in spaces
}

func NewCustomLabel(text string, sty *style) *widget.Label {
	lbl := widget.NewLabel(text)
	if sty != nil {
		if sty.imp > 0 {
			lbl.Importance = sty.imp
		}
		if sty.trun > 0 {
			lbl.Truncation = sty.trun
		}
		if sty.wrap > 0 {
			lbl.Wrapping = sty.wrap
		}
		if sty.align > 0 {
			lbl.Alignment = sty.align
		}
		lbl.TextStyle = fyne.TextStyle{
			Bold:      sty.Bold,
			Italic:    sty.Italic,
			Monospace: sty.Monospace,
			Symbol:    sty.Symbol,
			TabWidth:  sty.TabWidth,
		}
	}
	return lbl
}

// NewCustomLocaleLabel creates a new label widget with the set text content with a specific style
func NewCustomLocaleLabel(id string, sty *style) *widget.Label {
	return NewCustomLabel(txt.GetLabel(id), sty)
}

func NewLocaleLabel(id string) *widget.Label {
	return widget.NewLabel(txt.GetLabel(id))
}
func NewTitleLocaleText(id string) *canvas.Text {
	return NewTitleText(txt.GetLabel(id))
}

func NewTitleText(text string) *canvas.Text {
	return utility.NewSizeableColorText(text, 20, color.NRGBA{R: 0xFE, G: 0x58, B: 0x62, A: 0xFF})
}
