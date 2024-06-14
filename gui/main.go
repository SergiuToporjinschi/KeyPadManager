package gui

import (
	"main/logger"
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	usbMonitor     *monitor.USBMonitor
	window         fyne.Window
	device         *monitor.ConnectedDevice
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

// ------------ Content Manager ------------

type ContentManager struct {
	container.Split
}

func NewContentManager() *ContentManager {
	s := &ContentManager{
		Split: container.Split{
			Offset:     0.5, // Sensible default, can be overridden with SetOffset
			Horizontal: true,
			Leading:    NewMenuOptions(),
			Trailing:   container.NewStack(widget.NewLabel("Menu selection")),
		},
	}
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

func (c *ContentManager) SetDevice(device *monitor.ConnectedDevice) {
	logger.Log.Debugf("Set device: %v", device)
	// c.Trailing.(*MenuOptions).SetDevice("info", NewDeviceInfo())
}

// ------------ Menu Options ------------
type MenuOptions struct {
	widget.Accordion
}

func NewMenuOptions() *MenuOptions {
	menu := &MenuOptions{}
	menu.ExtendBaseWidget(menu)
	inf := widget.NewButton("Info", func() { menu.onSelectionChanged("info", NewDeviceInfo()) })
	menu.Append(widget.NewAccordionItem("Device", container.NewVBox(inf)))
	return menu
}

func (m *MenuOptions) onSelectionChanged(title string, content MainContent) {
	logger.Log.Debugf("Menu selection changed: %v, %v", title, content)
}

func (m *MenuOptions) SetDevice(title string, content MainContent) {
	logger.Log.Debugf("Set device: %v, %v", title, content)
}

type MainContent interface {
	Build() *fyne.Container
	SetDevice(*monitor.ConnectedDevice)
	Destroy()
}
