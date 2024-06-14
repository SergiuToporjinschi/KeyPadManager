package gui

import (
	"main/logger"
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	usbMonitor *monitor.USBMonitor
	window     fyne.Window
	// device         *monitor.ConnectedDevice
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

// ------------------------ Content Manager ------------------------

type ContentManager struct {
	container.Split
	currentDevice *monitor.ConnectedDevice
}

func NewContentManager() *ContentManager {
	s := &ContentManager{
		Split: container.Split{
			Offset:     0.5, // Sensible default, can be overridden with SetOffset
			Horizontal: true,
			Trailing:   container.NewStack(),
		},
	}
	s.Split.Leading = NewMenuOptions(s.onMenuClicked)
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

func (c *ContentManager) SetDevice(device *monitor.ConnectedDevice) {
	c.currentDevice = device
}

func (c *ContentManager) onMenuClicked(content MainContent) {
	logger.Log.Debugf("Menu selection changed: device: %v; content: %v", c.currentDevice, content)
	c.Trailing.(*fyne.Container).Add(content.GetContent(c.currentDevice))
	// c.Split.Resize(fyne.NewSize(100, 100)) //TODO change size
	c.Refresh()
}

// ------------------------ Menu Options ------------------------
type MenuOptions struct {
	widget.Accordion
}

func NewMenuOptions(onClick func(MainContent)) *MenuOptions {
	menu := &MenuOptions{}
	menu.ExtendBaseWidget(menu)
	devInf := NewDeviceInfo()
	inf := widget.NewButton("Info", func() { onClick(devInf) })
	// inf := widget.NewButton("Info", func() { onClick(NewDeviceInfo()) })
	//.......
	// inf := widget.NewButton("Info", func() { onClick(NewDeviceInfo()) })
	menu.Append(widget.NewAccordionItem("Device", container.NewVBox(inf)))
	return menu
}

type MainContent interface {
	GetContent(*monitor.ConnectedDevice) *fyne.Container
	Destroy()
}
