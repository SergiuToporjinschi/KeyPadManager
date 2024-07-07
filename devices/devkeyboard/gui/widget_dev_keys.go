package devkeyboardGui

import (
	"log/slog"
	resources "main/assets"
	"main/devicelayout"
	"main/devices/devkeyboard"
	"main/gui/widgets"
	"main/utility"

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
	binderFocused   widgets.Binder
	keysImages      map[int]*canvas.Image
	knobImage       *canvas.Image
	devDescriptor   *devicelayout.DeviceDescriptor
}

func New(data binding.Bytes, devDesc *devicelayout.DeviceDescriptor) fyne.CanvasObject {
	ins := &KeySwitchControl{
		keySelImageRes:  resources.ResButtonPng,
		keyImageRes:     resources.ResButtonGrayPng,
		knobImageRes:    resources.ResKnobGrayPng,
		knobSelImageRes: resources.ResKnobPng,
		keysImages:      make(map[int]*canvas.Image),
	}
	ins.devDescriptor = devkeyboard.ConvertHardwareDescriptor(devDesc)
	ins.ExtendBaseWidget(ins)
	ins.binderFocused.Bind(data, ins.updateFromData)
	ins.createImages()
	return ins
}

func (w *KeySwitchControl) createImages() {
	hrdDesc := w.devDescriptor.HardwareDescriptor.(devkeyboard.DevKeyboardComponent)
	for _, comp := range hrdDesc.Keys {
		if comp.Type == "button" {
			w.keysImages[comp.Value] = canvas.NewImageFromResource(w.keyImageRes)
			w.keysImages[comp.Value].FillMode = canvas.ImageFillContain
		}
	}
	w.knobImage = canvas.NewImageFromResource(w.knobImageRes)
	w.knobImage.FillMode = canvas.ImageFillContain
}

func (w *KeySwitchControl) updateFromData(dataBinding binding.DataItem) {
	if dataBinding == nil {
		return
	}
	source, ok := dataBinding.(binding.Bytes)

	if !ok {
		slog.Error("Received data is not a binding.Bytes")
		return
	}

	binaryData, err := source.Get()
	if err != nil {
		slog.Error("Error getting current data value", "error", err)
		return
	}

	slog.Debug("Binary data", "binary", utility.AsBinaryString(binaryData))
	slog.Debug("Hex data", "hex", utility.AsHexString(binaryData))
	w.setSelected(binaryData)
}

func (r *KeySwitchControl) CreateRenderer() fyne.WidgetRenderer {
	return newDevRenderer(r.keysImages, r.knobImage, r)
}

func (w *KeySwitchControl) setSelected(data []byte) {
	hrdDesc := w.devDescriptor.HardwareDescriptor.(devkeyboard.DevKeyboardComponent)
	values, encBtnVal, rotValue := devkeyboard.DecodeBinaryValue(data, hrdDesc)

	slog.Debug("Decoded data", "values", values.Keys(), "rotary", rotValue, "encoder", encBtnVal)

	for key := range w.keysImages {
		if values.Contains(key) {
			w.keysImages[key].Resource = w.keySelImageRes
		} else {
			w.keysImages[key].Resource = w.keyImageRes
		}
	}

	if rotValue > 0 {
		res, err := utility.RotateImageResource(w.knobSelImageRes.(*fyne.StaticResource), 70)
		if err != nil {
			slog.Error("Error rotating image", "error", err)
			w.knobImage.Resource = w.knobSelImageRes
		} else {
			w.knobImage.Resource = res
		}
	} else if rotValue < 0 {
		res, err := utility.RotateImageResource(w.knobSelImageRes.(*fyne.StaticResource), -70)
		if err != nil {
			slog.Error("Error rotating image", "error", err)
			w.knobImage.Resource = w.knobSelImageRes
		} else {
			w.knobImage.Resource = res
		}
	} else {
		w.knobImage.Resource = w.knobImageRes
		if encBtnVal != 0 {
			w.knobImage.Resource = w.knobSelImageRes
		} else {
			w.knobImage.Resource = w.knobImageRes
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

	//key display
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

	//knob display
	inst.knobImage.SetMinSize(inst.iconSize)
	inst.knobImage.Resize(inst.iconSize)
	inst.knobImage.Move(fyne.NewPos(inst.margin+xPos*4+vShift, yPos+inst.margin+vShift*2+inst.iconSSize+inst.vGap))
}

func (inst *devRanderer) MinSize() fyne.Size {
	return fyne.NewSize(inst.iconSSize*5+inst.margin*2, (inst.iconSSize*3)+(inst.margin*2)+(inst.iconSSize*0.3*2))
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
