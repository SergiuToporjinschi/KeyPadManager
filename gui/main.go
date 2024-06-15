package gui

import (
	resources "main/assets"
	"main/monitor"

	"fyne.io/fyne/v2"
)

type MainWindow struct {
	usbMonitor     *monitor.USBMonitor
	window         fyne.Window
	contentManager *ContentManager
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
	m.window.SetIcon(resources.ResLogoPng)
	m.window.SetCloseIntercept(m.Close)
	m.window.CenterOnScreen()
	m.window.Resize(fyne.NewSize(800, 600))
	m.contentManager = NewContentManager()
	m.window.SetContent(m.contentManager)
}

func (m *MainWindow) Close() {
	m.window.Hide()
}

func (m *MainWindow) Show(device *monitor.ConnectedDevice) {
	m.contentManager.SetDevice(device)
	// m.window.Content().(*widget.Label).SetText(device.Identifier.String())
	m.window.Show()
}
