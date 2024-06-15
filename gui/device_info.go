package gui

import (
	"main/monitor"
	"main/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DeviceInfo struct {
	title       string
	navTitle    string
	button      *widget.Button
	body        *fyne.Container
	bindingData infoBindingData
}
type infoBindingData struct {
	Manufacturer binding.ExternalString
	Product      binding.ExternalString
	SerialNumber binding.ExternalString
	PID          binding.ExternalString
	VID          binding.ExternalString
}

func NewDeviceInfo(onNavClick func(NavigationItem)) *DeviceInfo {
	inst := &DeviceInfo{
		title:    "Info",
		navTitle: "Info",
		bindingData: infoBindingData{
			Manufacturer: binding.BindString(nil),
			Product:      binding.BindString(nil),
			SerialNumber: binding.BindString(nil),
			PID:          binding.BindString(nil),
			VID:          binding.BindString(nil),
		},
	}
	inst.buildButton(onNavClick)
	inst.buildBody()
	return inst
}

func (i *DeviceInfo) buildButton(onNavClick func(NavigationItem)) {
	i.button = widget.NewButton(i.title, func() {
		onNavClick(i)
	})
}

func (i *DeviceInfo) buildBody() {
	i.body = container.New(layout.NewFormLayout()) //TODO if not works then try to use i.body.Add()

	i.body.Add(utility.NewTitleLabel("PID:"))
	i.body.Add(widget.NewLabelWithData(i.bindingData.PID))

	i.body.Add(utility.NewTitleLabel("VID:"))
	i.body.Add(widget.NewLabelWithData(i.bindingData.VID))

	i.body.Add(utility.NewTitleLabel("Product:"))
	i.body.Add(widget.NewLabelWithData(i.bindingData.Product))

	i.body.Add(utility.NewTitleLabel("Manufacturer:"))
	i.body.Add(widget.NewLabelWithData(i.bindingData.Manufacturer))

	i.body.Add(utility.NewTitleLabel("Serial:"))
	i.body.Add(widget.NewLabelWithData(i.bindingData.SerialNumber))
}

func (i *DeviceInfo) setData(dev *monitor.ConnectedDevice) {
	i.bindingData.Manufacturer.Set(dev.Identifier.Manufacturer)
	i.bindingData.Product.Set(dev.Identifier.Product)
	i.bindingData.SerialNumber.Set(dev.Identifier.SerialNumber)
	i.bindingData.PID.Set(dev.Identifier.PID.String())
	i.bindingData.VID.Set(dev.Identifier.VID.String())
}

func (i *DeviceInfo) GetContent(dev *monitor.ConnectedDevice) *fyne.Container {
	i.setData(dev)
	return i.body
}

func (i *DeviceInfo) GetButton() *widget.Button {
	return i.button
}

func (i *DeviceInfo) GetTitle() string {
	return i.title
}

func (i *DeviceInfo) GetNavTitle() string {
	return i.navTitle
}

func (i *DeviceInfo) Destroy() {
	//TODO
}
