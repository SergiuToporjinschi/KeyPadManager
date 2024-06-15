package gui

import (
	resources "main/assets"
	"main/monitor"
	"main/txt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type SelectDeviceWindow struct {
	usbMonitor        *monitor.USBMonitor
	window            fyne.Window
	cardList          *fyne.Container
	selectDevEvent    chan string
	selectDevListners map[string]func(*monitor.ConnectedDevice)
}

func NewSelectDevice(myApp fyne.App) *SelectDeviceWindow {
	instance := &SelectDeviceWindow{
		usbMonitor: monitor.GetInstance(),
		window:     myApp.NewWindow(txt.GetLabel("win.selDevTitle")),
	}

	instance.selectDevEvent = make(chan string)
	instance.selectDevListners = make(map[string]func(*monitor.ConnectedDevice))

	instance.buildWindow()
	return instance
}

func (s *SelectDeviceWindow) buildWindow() {
	s.window.Hide()
	s.window.SetIcon(resources.ResLogoPng)
	s.window.SetCloseIntercept(s.Close)
	s.cardList = container.NewHBox(NewNoDeviceCard())
	s.window.SetContent(container.New(layout.NewCenterLayout(), s.cardList))

	s.usbMonitor.AddDeviceEvent("SelectDeviceWindow", s.updateCards)

	s.window.CenterOnScreen()
	s.window.Resize(fyne.NewSize(800, 600))
}

func (s *SelectDeviceWindow) Close() {
	s.window.Hide()
}

func (s *SelectDeviceWindow) updateCards(event string, device *monitor.ConnectedDevice) {
	if event == monitor.EventDeviceConnected {
		if len(s.cardList.Objects) == 1 && s.cardList.Objects[0].(*DeviceCard).IsDummy {
			s.cardList.Objects = []fyne.CanvasObject{}
		}

		card := NewDeviceCard(
			device,
			canvas.NewImageFromResource(resources.ResDevKeypadPng), //TODO make it dynamic (maybe save the image in the device object)
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

func (s *SelectDeviceWindow) AddSelectDeviceListener(name string, callback func(device *monitor.ConnectedDevice)) {
	s.selectDevListners[name] = callback

}

func (s *SelectDeviceWindow) onClickDevice(device *monitor.ConnectedDevice) {
	for _, listener := range s.selectDevListners {
		listener(device)
	}
}

func (s *SelectDeviceWindow) Show() {
	s.window.Show()
}

func (s *SelectDeviceWindow) Hide() {
	s.window.Hide()
}
