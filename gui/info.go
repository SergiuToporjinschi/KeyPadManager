package gui

import (
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Info struct {
	container *fyne.Container
	gui       *GUI
}

func NewInfo(gui *GUI) *Info {
	return &Info{
		gui: gui,
	}
}

func (i *Info) Init() {

}

func (i *Info) Build(device *usb.Device) *fyne.Container {
	items := []fyne.CanvasObject{}

	items = append(items, widget.NewLabel("PID:"))
	if device != nil {
		pid := device.Info.PID.String()
		items = append(items, widget.NewLabelWithData(binding.BindString(&pid)))
	} else {
		items = append(items, widget.NewLabel(""))
	}

	items = append(items, widget.NewLabel("VID:"))
	if device != nil {
		vid := device.Info.VID.String()
		items = append(items, widget.NewLabelWithData(binding.BindString(&vid)))
	} else {
		items = append(items, widget.NewLabel(""))
	}

	items = append(items, widget.NewLabel("Product:"))
	if device != nil {
		items = append(items, widget.NewLabelWithData(binding.BindString(&device.Info.Product)))
	} else {
		items = append(items, widget.NewLabel(""))
	}

	items = append(items, widget.NewLabel("Manufacturer:"))
	if device != nil {
		items = append(items, widget.NewLabelWithData(binding.BindString(&device.Info.Manufacturer)))
	} else {
		items = append(items, widget.NewLabel(""))
	}

	items = append(items, widget.NewLabel("Serial:"))
	if device != nil {
		items = append(items, widget.NewLabelWithData(binding.BindString(&device.Info.Manufacturer)))
	} else {
		items = append(items, widget.NewLabel(""))
	}
	i.container = container.New(layout.NewFormLayout(), items...)
	return i.container
}

func (i *Info) Destroy() {
	if i.container != nil {
		i.container.Hide()
	}
}
