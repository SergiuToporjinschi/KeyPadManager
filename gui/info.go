package gui

//go:generate go run ./ bundle -package data -o bundled.go assets

import (
	"main/devicelayout"
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DeviceInfo struct {
	*fyne.Container
	deviceBindingData *devicelayout.DeviceLayoutConfig
}

func NewDeviceInfo() *DeviceInfo {
	instance := &DeviceInfo{
		Container: &fyne.Container{
			Hidden: true,
			Layout: layout.NewVBoxLayout(),
		},
		deviceBindingData: &devicelayout.DeviceLayoutConfig{
			Identifier: devicelayout.DevIdentifier{
				Manufacturer: "",
				Product:      "",
				SerialNumber: "",
				PID:          0,
				VID:          0,
			},
		},
	}
	instance.build()
	return instance
}

func (i *DeviceInfo) build() {
	items := []fyne.CanvasObject{}

	items = append(items, widget.NewLabel("PID:"))
	pid := i.deviceBindingData.Identifier.PID.String()
	items = append(items, widget.NewLabelWithData(binding.BindString(&pid)))

	items = append(items, widget.NewLabel("VID:"))
	vid := i.deviceBindingData.Identifier.VID.String()
	items = append(items, widget.NewLabelWithData(binding.BindString(&vid)))

	items = append(items, widget.NewLabel("Product:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.deviceBindingData.Identifier.Product)))

	items = append(items, widget.NewLabel("Manufacturer:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.deviceBindingData.Identifier.Manufacturer)))

	items = append(items, widget.NewLabel("Serial:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.deviceBindingData.Identifier.SerialNumber)))

	i.Objects = append(i.Objects, container.New(layout.NewFormLayout(), items...))
}

func (i *DeviceInfo) GetContent(dev *monitor.ConnectedDevice) *fyne.Container {
	i.deviceBindingData.Identifier.Manufacturer = dev.Identifier.Manufacturer
	i.deviceBindingData.Identifier.Product = dev.Identifier.Product
	i.deviceBindingData.Identifier.SerialNumber = dev.Identifier.SerialNumber
	i.deviceBindingData.Identifier.PID = dev.Identifier.PID
	i.deviceBindingData.Identifier.VID = dev.Identifier.VID
	i.Container.Show()
	return i.Container
}

func (i *DeviceInfo) Destroy() {
}
