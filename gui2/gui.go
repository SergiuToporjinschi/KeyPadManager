package gui

import (
	"main/logger"
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainContent interface {
	Build(*usb.Device) *fyne.Container
	DeviceSelectionChanged(*usb.Device)
	Destroy()
}

type GUI struct {
	application fyne.App
	device      *usb.Device
	menuOptions menuOptions
	mainWindow  fyne.Window
}

func NewGUI() *GUI {
	return &GUI{
		application: app.New(),
	}
}

func (g *GUI) OpenMain() {
	g.menuOptions = g.getMenuConfig()
	g.mainWindow = g.application.NewWindow("Device Manager")
	g.mainWindow.Resize(fyne.NewSize(800, 600))
	g.mainWindow.CenterOnScreen()
	g.mainWindow.SetPadded(true)

	mainContent := container.NewStack()
	menu := g.buidMenu(mainContent)

	menuBar := container.NewBorder(
		nil,                                 //top
		nil,                                 //bottom
		nil,                                 //left
		NewDeviceList(g.onSelectionChanged), //right
		nil,                                 //center
	)

	main := container.NewBorder(
		menuBar,                                //top
		nil,                                    //bootom
		nil,                                    //left
		nil,                                    //right
		container.NewHSplit(menu, mainContent), //center
	)

	g.mainWindow.SetContent(container.NewStack(main))
	g.mainWindow.Show()
	g.application.Run()
}

func (g *GUI) onSelectionChanged(dev *usb.Device) {
	logger.Log.Infof("GUI: iterate menu trigger onSelectionChanged, %v", dev)
	g.menuOptions["Device"]["Info"].DeviceSelectionChanged(dev)
	g.menuOptions["Device"]["RawValues"].DeviceSelectionChanged(dev)
}

type MenuBarWidget struct {
	widget.BaseWidget
	container *fyne.Container
}

func NewMenuBarWidget(objects ...fyne.CanvasObject) *MenuBarWidget {
	item := &MenuBarWidget{
		container: container.NewHBox(objects...),
	}
	item.ExtendBaseWidget(item)
	return item
}

func (item *MenuBarWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, nil, item.container, nil))
}

func (item *MenuBarWidget) Add(object ...fyne.CanvasObject) {
	for _, obj := range object {
		item.container.Add(obj)
	}
	item.Refresh()
}

func (item *MenuBarWidget) RemoveAll() {
	item.container.Objects = nil
	item.Refresh()
}
