package gui

import (
	"main/logger"
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SelectDeviceWindow struct {
	usbMonitor *monitor.USBMonitor
	window     fyne.Window
	cardList   *fyne.Container
}

func NewSelectDevice(myApp fyne.App) *SelectDeviceWindow {
	instance := &SelectDeviceWindow{
		usbMonitor: monitor.GetInstance(),
		window:     myApp.NewWindow("Select device"),
	}
	instance.buildWindow()
	return instance
}

func (s *SelectDeviceWindow) buildWindow() {
	s.window.Hide()
	s.window.SetContent(widget.NewLabel("Select device"))
	s.window.SetCloseIntercept(func() {
		instance.SelectDeviceWindow.Hide()
	})
	s.cardList = container.NewHBox(NewNoDeviceCard())
	s.window.SetContent(container.New(layout.NewCenterLayout(), s.cardList))

	s.usbMonitor.AddDeviceEvent("SelectDeviceWindow", s.updateCards)

	s.window.CenterOnScreen()
	s.window.Resize(fyne.NewSize(800, 600))
}

func (s *SelectDeviceWindow) updateCards(event string, device monitor.ConnectedDevice) {
	if event == monitor.EventDeviceConnected {
		if len(s.cardList.Objects) == 1 && s.cardList.Objects[0].(*DeviceCard).IsDummy {
			s.cardList.Objects = []fyne.CanvasObject{}
		}

		card := NewDeviceCard(
			device,
			canvas.NewImageFromFile("devKeypad.png"),
			s.onClickDevice,
			nil)

		s.cardList.Add(card)
	} else if event == monitor.EventDeviceDisconnected {
		for _, obj := range s.cardList.Objects {
			card := obj.(*DeviceCard)
			if card.device.Identifier.String() == device.Identifier.String() {
				s.cardList.Remove(card)
				break
			}
		}
	}

	if len(s.cardList.Objects) == 0 {
		s.cardList.Add(NewNoDeviceCard())
	}

	s.cardList.Refresh()
}

func (s *SelectDeviceWindow) onClickDevice(device monitor.ConnectedDevice) {
	logger.Log.Infof("Device clicked %s", device.Identifier.String())
}

func (s *SelectDeviceWindow) Show() {
	s.window.Show()
}

func (s *SelectDeviceWindow) Hide() {
	s.window.Hide()
}

type DeviceCard struct {
	widget.Card
	device       monitor.ConnectedDevice
	onLeftClick  func(monitor.ConnectedDevice)
	onRightClick func(monitor.ConnectedDevice)
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

func NewDeviceCard(device monitor.ConnectedDevice, img *canvas.Image, onLeftClick, onRightClick func(monitor.ConnectedDevice)) *DeviceCard {
	img.FillMode = canvas.ImageFillOriginal

	card := &DeviceCard{
		device:       device,
		onLeftClick:  onLeftClick,
		onRightClick: onRightClick,
	}

	card.ExtendBaseWidget(card)
	card.SetTitle(device.Identifier.Name)
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
