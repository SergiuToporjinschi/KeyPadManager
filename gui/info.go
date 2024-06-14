package gui

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
	fyne.Container
	deviceBindingData *devicelayout.DeviceLayoutConfig
}

func NewDeviceInfo() *DeviceInfo {
	instance := &DeviceInfo{
		Container: fyne.Container{
			Hidden: true,
			Layout: layout.NewVBoxLayout(),
		},
		deviceBindingData: &devicelayout.DeviceLayoutConfig{
			Identifier: devicelayout.DevIdentifier{
				Manufacturer: "",
				Name:         "",
				PID:          0,
				VID:          0,
			},
		},
	}
	instance.Objects = append(instance.Objects, instance.Build())
	return instance
}

func (i *DeviceInfo) Build() *fyne.Container {
	items := []fyne.CanvasObject{}

	items = append(items, widget.NewLabel("PID:"))
	pid := i.deviceBindingData.Identifier.PID.String()
	items = append(items, widget.NewLabelWithData(binding.BindString(&pid)))

	items = append(items, widget.NewLabel("VID:"))
	vid := i.deviceBindingData.Identifier.VID.String()
	items = append(items, widget.NewLabelWithData(binding.BindString(&vid)))

	items = append(items, widget.NewLabel("Product:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.deviceBindingData.Identifier.Name)))

	items = append(items, widget.NewLabel("Manufacturer:"))
	items = append(items, widget.NewLabelWithData(binding.BindString(&i.deviceBindingData.Identifier.Manufacturer)))

	// items = append(items, widget.NewLabel("Serial:"))
	// items = append(items, widget.NewLabelWithData(binding.BindString(&i.deviceBindingData.SerialNumber)))

	return container.New(layout.NewFormLayout(), items...)
}

func (i *DeviceInfo) SetDevice(dev *monitor.ConnectedDevice) {
	// i.deviceBindingData.Identifier.Manufacturer = dev.DeviceLayoutConfig.Identifier.Manufacturer
	// i.deviceBindingData.Identifier.Product = dev.DeviceLayoutConfig.Identifier.Name
	// i.deviceBindingData.Identifier.PID = uint16(dev.DeviceLayoutConfig.Identifier.PID)
	// i.deviceBindingData.Identifier.VID = dev.DeviceLayoutConfig.Identifier.VID
	// i.deviceBindingData.Identifier.SerialNumber = dev.DeviceLayoutConfig.Identifier.SerialNumber
}

func (i *DeviceInfo) Destroy() {
}
