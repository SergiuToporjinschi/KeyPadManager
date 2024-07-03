package widgets

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// BorderedControlImage is a custom widget that displays an image with an optional border.
type BorderedControlImage struct {
	widget.BaseWidget
	image         *canvas.Image
	border        bool
	borderSize    float32
	padding       float32
	binderFocused Binder
}

// NewImageWithBorder creates a new ImageWithBorder widget.
func NewBorderedControlImage(imageResource fyne.Resource, data binding.Bool) *BorderedControlImage {
	img := canvas.NewImageFromResource(imageResource)
	img.SetMinSize(fyne.NewSize(64, 64))

	img.FillMode = canvas.ImageFillContain
	w := &BorderedControlImage{
		image:      img,
		border:     false,
		borderSize: 2, // Default border size
		padding:    5, // Default padding
	}
	w.ExtendBaseWidget(w)
	w.binderFocused.Bind(data, w.updateFromData)
	return w
}

func (w *BorderedControlImage) updateFromData(data binding.DataItem) {
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
	w.SetBorder(val)
}

// CreateRenderer implements the widget.WidgetRenderer interface.
func (w *BorderedControlImage) CreateRenderer() fyne.WidgetRenderer {
	wr := &borderedImageRenderer{
		image:  w.image,
		border: canvas.NewRectangle(color.Transparent),
		widget: w,
	}

	wr.updateBorder()
	return wr
}

// SetBorder toggles the border on or off.
func (w *BorderedControlImage) SetBorder(show bool) {
	w.border = show
	w.Refresh()
}

type borderedImageRenderer struct {
	image  *canvas.Image
	border *canvas.Rectangle
	widget *BorderedControlImage
}

func (r *borderedImageRenderer) Layout(size fyne.Size) {
	padding := r.widget.padding
	borderSize := r.widget.borderSize
	r.image.Resize(size)

	r.border.Resize(fyne.NewSize(r.image.Size().Height+padding*2+borderSize*2, size.Height+padding*2+borderSize*2))
	r.border.Move(fyne.NewPos(size.Width/2-r.image.Size().Height/2-padding-borderSize, r.border.Position().Y-padding+borderSize-1))
}

func (r *borderedImageRenderer) MinSize() fyne.Size {
	innerMinSize := r.image.MinSize()
	padding := r.widget.padding
	borderSize := r.widget.borderSize

	return fyne.NewSize(
		innerMinSize.Width+2*padding+2*borderSize,
		innerMinSize.Height+2*padding+2*borderSize,
	)
}

func (r *borderedImageRenderer) Refresh() {
	r.updateBorder()
	canvas.Refresh(r.widget)
}

func (r *borderedImageRenderer) updateBorder() {
	if r.widget.border {
		r.border.StrokeColor = color.White
		r.border.StrokeWidth = r.widget.borderSize
	} else {
		r.border.StrokeColor = color.Transparent
		r.border.StrokeWidth = 0
	}
	r.border.Refresh()
}

func (r *borderedImageRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *borderedImageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.border, r.image}
}

func (r *borderedImageRenderer) Destroy() {}
