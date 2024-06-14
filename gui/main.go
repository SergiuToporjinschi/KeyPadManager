package gui

import (
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContent interface {
	GetContent(*monitor.ConnectedDevice) *fyne.Container
	Destroy()
}

type MainWindow struct {
	usbMonitor *monitor.USBMonitor
	window     fyne.Window
	// device         *monitor.ConnectedDevice
	contentManager *Navigation
}

func NewMainWindow(myApp fyne.App) *MainWindow {
	instance := &MainWindow{
		usbMonitor: monitor.GetInstance(),
		window:     myApp.NewWindow("Manager"), //TOOD change title
	}
	instance.buildWindow()
	return instance
}

func (m *MainWindow) buildWindow() {
	m.window.Hide()
	m.window.SetCloseIntercept(m.Close)
	m.window.SetPadded(true)
	m.window.CenterOnScreen()
	m.window.Resize(fyne.NewSize(800, 600))
	m.contentManager = NewContentManager()
	m.window.SetContent(container.NewStack(
		container.NewBorder(
			nil, //bootom
			nil, //left
			nil, //right
			m.contentManager,
		)))
}

func (m *MainWindow) Close() {
	m.window.Hide()
}

func (m *MainWindow) Show(device *monitor.ConnectedDevice) {
	m.contentManager.SetDevice(device)
	// m.window.Content().(*widget.Label).SetText(device.Identifier.String())
	m.window.Show()
}
