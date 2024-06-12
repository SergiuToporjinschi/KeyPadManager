package gui

import (
	"main/logger"
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Info struct {
	container *fyne.Container
	device    *usb.DevInfo
}

func NewInfo() *Info {
	return &Info{
		device: &usb.DevInfo{
			Manufacturer: "",
			Product:      "",
			SerialNumber: "",
			PID:          0,
			VID:          0,
		},
	}
}

func (i *Info) DeviceSelectionChanged(dev *usb.Device) {
	logger.Log.Infof("INFO: Device selection changed %v", dev)
	i.device.Manufacturer = dev.Info.Manufacturer
	i.device.Product = dev.Info.Product
	i.device.SerialNumber = dev.Info.SerialNumber
	i.device.PID = dev.Info.PID
	i.device.VID = dev.Info.VID
}

func (i *Info) Build(device *usb.Device) *fyne.Container {
	items := []fyne.CanvasObject{}

	items = append(items, widget.NewLabel("PID:"))
	pid := i.device.PID.String()
	items = append(items, widget.NewLabelWithData(binding.BindString(&pid)))

	items = append(items, widget.NewLabel("VID:"))
	vid := i.device.VID.String()
	items = append(items, widget.NewLabelWithData(binding.BindString(&vid)))

	items = append(items, widget.NewLabel("Product:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.device.Product)))

	items = append(items, widget.NewLabel("Manufacturer:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.device.Manufacturer)))

	items = append(items, widget.NewLabel("Serial:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.device.SerialNumber)))

	i.container = container.New(layout.NewFormLayout(), items...)
	return i.container
}

func (i *Info) Destroy() {
	if i.container != nil {
		i.container.Hide()
	}
}
