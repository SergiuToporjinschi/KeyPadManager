package widgets

import (
	"log/slog"
	resources "main/assets"
	"main/devicelayout"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

//grayout 8c8c8c
//light gray e5e5e5

// KeySwitchControl is a custom widget that displays an image with an optional border.
type KeySwitchControl struct {
	widget.BaseWidget
	keySelImageRes  fyne.Resource
	keyImageRes     fyne.Resource
	knobImageRes    fyne.Resource
	knobSelImageRes fyne.Resource
	binderFocused   Binder
	keysImages      map[int]*canvas.Image
	knobImage       *canvas.Image
	layout          *devicelayout.DeviceLayoutConfig
}

// NewKeySwitchControl creates a new ImageWithBorder widget.
func NewKeySwitchControl(data binding.ExternalInt, layout *devicelayout.DeviceLayoutConfig) *KeySwitchControl {
	// rotated, err := utility.RotateImageResource(resources.ResKnobPng, 45)
	// if err != nil {
	// 	slog.Error("Error rotating image resource", "error", err)
	// 	return nil
	// }

	w := &KeySwitchControl{
		keySelImageRes:  resources.ResButtonPng,
		keyImageRes:     resources.ResButtonGrayPng,
		knobImageRes:    resources.ResKnobGrayPng,
		knobSelImageRes: resources.ResKnobPng,
		keysImages:      make(map[int]*canvas.Image),
		layout:          layout,
	}

	w.ExtendBaseWidget(w)
	w.binderFocused.Bind(data, w.updateFromData)
	w.createImages()
	return w
}

func (w *KeySwitchControl) createImages() {
	for _, comp := range w.layout.Components {
		if comp.Type == "button" && comp.Value != 1024 {
			w.keysImages[comp.Value] = canvas.NewImageFromResource(w.keyImageRes)
			w.keysImages[comp.Value].FillMode = canvas.ImageFillContain
		}

		if comp.Type == "encoder" {
			w.knobImage = canvas.NewImageFromResource(w.knobImageRes)
			w.knobImage.FillMode = canvas.ImageFillContain
		}
	}
}

func (w *KeySwitchControl) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	boolSource, ok := data.(binding.Int)

	if !ok {
		return
	}

	val, err := boolSource.Get()
	if err != nil {
		slog.Error("Error getting current data value", "error", err)
		return
	}
	slog.Info("KeySwitchControl", "val", val)
	w.setSelected(val)
}

func (r *KeySwitchControl) CreateRenderer() fyne.WidgetRenderer {
	return newDevRenderer(r.keysImages, r.knobImage, r)
}

func (w *KeySwitchControl) setSelected(value int) {
	for key := range w.keysImages {
		if key&value != 0 {
			w.keysImages[key].Resource = w.keySelImageRes
		} else {
			w.keysImages[key].Resource = w.keyImageRes
		}
	}
	w.Refresh()
}

type devRanderer struct {
	iconSSize  float32
	margin     float32
	hGap       float32
	vGap       float32
	iconSize   fyne.Size
	keysImages map[int]*canvas.Image
	knobImage  *canvas.Image
	widget     *KeySwitchControl
}

func newDevRenderer(keysImages map[int]*canvas.Image, knobImage *canvas.Image, widget *KeySwitchControl) fyne.WidgetRenderer {
	inst := &devRanderer{
		iconSSize:  64,
		margin:     10,
		hGap:       5,
		vGap:       5,
		keysImages: keysImages,
		knobImage:  knobImage,
		widget:     widget,
	}

	inst.iconSize = fyne.NewSize(inst.iconSSize, inst.iconSSize)
	for i := range inst.keysImages {
		inst.keysImages[i].SetMinSize(inst.iconSize)
		inst.keysImages[i].Resize(inst.iconSize)
	}

	inst.knobImage.SetMinSize(inst.iconSize)
	inst.knobImage.Resize(inst.iconSize)

	inst.Layout(inst.MinSize())
	return inst

}
func (inst *devRanderer) Layout(_ fyne.Size) {
	var vShift float32 = inst.iconSSize * 0.3

	var xPos float32 = inst.hGap + inst.iconSSize
	var yPos float32 = inst.vGap + inst.iconSSize

	for i := range inst.keysImages {
		inst.keysImages[i].SetMinSize(inst.iconSize)
		inst.keysImages[i].Resize(inst.iconSize)
	}

	inst.keysImages[1].Move(fyne.NewPos(inst.margin, inst.margin+vShift*2))
	inst.keysImages[2].Move(fyne.NewPos(xPos+inst.margin, inst.margin+vShift))
	inst.keysImages[4].Move(fyne.NewPos(inst.margin+xPos*2, inst.margin))
	inst.keysImages[8].Move(fyne.NewPos(inst.margin+xPos*3, inst.margin+vShift))
	inst.keysImages[16].Move(fyne.NewPos(inst.margin+xPos*4, inst.margin+vShift*2))
	inst.keysImages[32].Move(fyne.NewPos(inst.margin, yPos+inst.margin+vShift*2))
	inst.keysImages[64].Move(fyne.NewPos(xPos+inst.margin, yPos+inst.margin+vShift))
	inst.keysImages[128].Move(fyne.NewPos(inst.margin+xPos*2, yPos+inst.margin))
	inst.keysImages[256].Move(fyne.NewPos(inst.margin+xPos*3, yPos+inst.margin+vShift))
	inst.keysImages[512].Move(fyne.NewPos(inst.margin+xPos*4, yPos+inst.margin+vShift*2))

	inst.knobImage.SetMinSize(inst.iconSize)
	inst.knobImage.Resize(inst.iconSize)
	inst.knobImage.Move(fyne.NewPos(inst.margin+xPos*4+vShift, yPos+inst.margin+vShift*2+inst.iconSSize+inst.vGap))
}

func (inst *devRanderer) MinSize() fyne.Size {
	return fyne.NewSize(inst.iconSSize*5+inst.margin*2, inst.iconSSize*3+inst.margin*2)
}

func (inst *devRanderer) Refresh() {
	inst.Layout(inst.MinSize())
	for i := range inst.keysImages {
		inst.keysImages[i].Refresh()
	}
	inst.knobImage.Refresh()

	canvas.Refresh(inst.widget)
}

func (inst *devRanderer) Objects() []fyne.CanvasObject {
	result := []fyne.CanvasObject{}
	for key, _ := range inst.keysImages {
		result = append(result, inst.keysImages[key])
	}
	result = append(result, inst.knobImage)
	return result
}

func (inst *devRanderer) Destroy() {}
