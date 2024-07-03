package widgets

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// KeySwitchControl is a custom widget that displays an image with an optional border.
type KeySwitchControl struct {
	widget.BaseWidget
	selImageRes   fyne.Resource
	imageRes      fyne.Resource
	image         *canvas.Image
	selected      bool
	binderFocused Binder
}

// NewKeySwitchControl creates a new ImageWithBorder widget.
func NewKeySwitchControl(imageResource fyne.Resource, selImageResource fyne.Resource, data binding.Bool) *KeySwitchControl {
	img := canvas.NewImageFromResource(imageResource)
	img.SetMinSize(fyne.NewSize(64, 64))

	img.FillMode = canvas.ImageFillContain
	w := &KeySwitchControl{
		selImageRes: selImageResource,
		imageRes:    imageResource,
		image:       img,
		selected:    false,
	}

	w.ExtendBaseWidget(w)
	w.binderFocused.Bind(data, w.updateFromData)

	return w
}

func (w *KeySwitchControl) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	boolSource, ok := data.(binding.Bool)

	if !ok {
		return
	}

	val, err := boolSource.Get()
	if err != nil {
		slog.Error("Error getting current data value", "error", err)
		return
	}
	w.setSelected(val)
}

// CreateRenderer implements the widget.WidgetRenderer interface.
func (w *KeySwitchControl) CreateRenderer() fyne.WidgetRenderer {
	wr := &keySwitchControlRenderer{
		image:       w.image,
		selImageRes: w.selImageRes,
		imageRes:    w.imageRes,
		selected:    w.selected,
		widget:      w,
	}

	wr.updateSelection()
	return wr
}
func (w *KeySwitchControl) setSelected(show bool) {
	w.selected = show
	if show {
		w.image.Resource = w.selImageRes
	} else {
		w.image.Resource = w.imageRes
	}
	w.Refresh()
}

type keySwitchControlRenderer struct {
	selImageRes fyne.Resource
	imageRes    fyne.Resource

	image    *canvas.Image
	selected bool
	widget   *KeySwitchControl
}

func (r *keySwitchControlRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
}

func (r *keySwitchControlRenderer) MinSize() fyne.Size {
	return r.image.MinSize()
}

func (r *keySwitchControlRenderer) Refresh() {
	r.updateSelection()
	canvas.Refresh(r.widget)
}

func (r *keySwitchControlRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.image}
}

func (r *keySwitchControlRenderer) Destroy() {}

func (r *keySwitchControlRenderer) updateSelection() {
	if r.widget.selected {
		r.image.Resource = r.selImageRes
	} else {
		r.image.Resource = r.imageRes
	}
	r.image.Refresh()
}

func (r *keySwitchControlRenderer) BackgroundColor() color.Color {
	return color.Transparent
}
