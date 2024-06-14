package gui

import (
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	usbMonitor *monitor.USBMonitor
	window     fyne.Window
	monitor.ConnectedDevice
}

func NewMainWindow(myApp fyne.App) *MainWindow {
	instance := &MainWindow{
		usbMonitor: monitor.GetInstance(),
		window:     myApp.NewWindow("Manager"), //TOOD change title
	}
	instance.buildWindow()
	return instance
}

func (s *MainWindow) buildWindow() {
	s.window.Hide()
	s.window.SetContent(widget.NewLabel("Manager device"))
	s.window.SetCloseIntercept(s.Close)

	s.window.CenterOnScreen()
	s.window.Resize(fyne.NewSize(800, 600))
}

func (s *MainWindow) Close() {
	s.window.Hide()
}

func (s *MainWindow) Show(device *monitor.ConnectedDevice) {
	s.window.Content().(*widget.Label).SetText(device.Identifier.String())
	s.window.Show()
}
