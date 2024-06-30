package gui

import (
	"log/slog"
	"main/monitor"
	"main/txt"
	"main/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DeviceInfo struct {
	title       string
	button      *widget.Button
	body        *container.Scroll
	bindingData infoBindingData
}
type infoBindingData struct {
	Manufacturer binding.ExternalString
	Product      binding.ExternalString
	SerialNumber binding.ExternalString
	PID          binding.ExternalString
	VID          binding.ExternalString
}

func NewDeviceInfo() NavigationItem {
	inst := &DeviceInfo{
		title: txt.GetLabel("navi.deviceInfoTitle"),
		bindingData: infoBindingData{
			Manufacturer: binding.BindString(nil),
			Product:      binding.BindString(nil),
			SerialNumber: binding.BindString(nil),
			PID:          binding.BindString(nil),
			VID:          binding.BindString(nil),
		},
	}
	inst.buildBody()
	return inst
}

func (i *DeviceInfo) buildBody() {
	i.body = container.NewVScroll(container.New(layout.NewFormLayout())) //TODO if not works then try to use i.body.Add()

	i.body.Content.(*fyne.Container).Add(utility.NewTitleLabel(txt.GetLabel("cont.pid")))
	i.body.Content.(*fyne.Container).Add(widget.NewLabelWithData(i.bindingData.PID))

	i.body.Content.(*fyne.Container).Add(utility.NewTitleLabel(txt.GetLabel("cont.vid")))
	i.body.Content.(*fyne.Container).Add(widget.NewLabelWithData(i.bindingData.VID))

	i.body.Content.(*fyne.Container).Add(utility.NewTitleLabel(txt.GetLabel("cont.product")))
	i.body.Content.(*fyne.Container).Add(widget.NewLabelWithData(i.bindingData.Product))

	i.body.Content.(*fyne.Container).Add(utility.NewTitleLabel(txt.GetLabel("cont.manufacturer")))
	i.body.Content.(*fyne.Container).Add(widget.NewLabelWithData(i.bindingData.Manufacturer))

	i.body.Content.(*fyne.Container).Add(utility.NewTitleLabel(txt.GetLabel("cont.serial")))
	i.body.Content.(*fyne.Container).Add(widget.NewLabelWithData(i.bindingData.SerialNumber))
}

func (i *DeviceInfo) setData(dev *monitor.ConnectedDevice) {
	i.bindingData.Manufacturer.Set(dev.Identifier.Manufacturer)
	i.bindingData.Product.Set(dev.Identifier.Product)
	i.bindingData.SerialNumber.Set(dev.Identifier.SerialNumber)
	i.bindingData.PID.Set(dev.Identifier.PID.String())
	i.bindingData.VID.Set(dev.Identifier.VID.String())
}

func (i *DeviceInfo) GetContent(dev *monitor.ConnectedDevice) *container.Scroll {
	i.setData(dev)
	return i.body
}

func (i *DeviceInfo) GetButton() *widget.Button {
	return i.button
}

func (i *DeviceInfo) Destroy() {
	slog.Debug("Destroying DeviceInfo")
}
