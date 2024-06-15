package gui

import (
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type DeviceCard struct {
	widget.Card
	device       *monitor.ConnectedDevice
	onLeftClick  func(*monitor.ConnectedDevice)
	onRightClick func(*monitor.ConnectedDevice)
	IsDummy      bool
}

func NewNoDeviceCard() *DeviceCard {
	card := &DeviceCard{
		IsDummy: true,
		Card: widget.Card{
			Title: "No device connected",
		},
	}
	card.ExtendBaseWidget(card)
	return card
}

func NewDeviceCard(device *monitor.ConnectedDevice, img *canvas.Image, onLeftClick, onRightClick func(*monitor.ConnectedDevice)) *DeviceCard {
	img.FillMode = canvas.ImageFillOriginal

	card := &DeviceCard{
		device:       device,
		onLeftClick:  onLeftClick,
		onRightClick: onRightClick,
	}

	card.ExtendBaseWidget(card)
	card.SetTitle(device.Identifier.Product)
	card.SetSubTitle(device.Identifier.String())
	card.SetImage(img)
	return card

}

func (b *DeviceCard) Tapped(_ *fyne.PointEvent) {
	if b.onLeftClick != nil {
		b.onLeftClick(b.device)
	}
}

func (b *DeviceCard) TappedSecondary(_ *fyne.PointEvent) {
	if b.onRightClick != nil {
		b.onRightClick(b.device)
	}
}
