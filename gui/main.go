package gui

import (
	resources "main/assets"
	"main/monitor"
	"main/txt"

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
		window:     myApp.NewWindow(txt.GetLabel("win.mainTitle")),
	}
	instance.buildWindow()
	instance.usbMonitor.AddDeviceEvent("contentManager", instance.OnDeviceEvent)
	return instance
}

func (m *MainWindow) buildWindow() {
	m.window.Hide()
	m.window.SetIcon(resources.ResLogoPng)
	m.window.SetCloseIntercept(m.Close)
	m.window.CenterOnScreen()
	m.window.Resize(fyne.NewSize(1224, 768))
	m.contentManager = NewContentManager()
	m.window.SetContent(m.contentManager)
}

func (m *MainWindow) Close() {

	// fyne.CurrentApp().Preferences().SetFloat("mainWindowX", size.X)
	// fyne.CurrentApp().Settings()
	m.contentManager.OnMainWindowHide()
	m.window.Hide()
}

func (m *MainWindow) Show(device *monitor.ConnectedDevice) {
	m.contentManager.SetDevice(device)
	m.window.Show()
}

func (m *MainWindow) OnDeviceEvent(event string, device *monitor.ConnectedDevice) {
	if event == "disconnected" {
		m.Close()
	}
}
