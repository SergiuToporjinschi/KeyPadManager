package gui

import (
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainContent interface {
	Build(device *usb.Device) *fyne.Container
	Destroy()
}

type GUI struct {
	application fyne.App
	device      *usb.Device
	menuOptions menuOptions
	mainWindow  fyne.Window
	menuBar     *MenuBarWidget
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
	g.menuBar = NewMenuBarWidget(widget.NewLabel("Device: "))
	if g.device != nil {
		g.menuBar.Add(widget.NewLabel(g.device.Info.Product))
	}
	g.mainWindow.SetContent(container.NewStack(container.NewBorder(g.menuBar, nil, nil, nil, container.NewHSplit(menu, mainContent))))
	g.mainWindow.Show()
	g.application.Run()
}

func (g *GUI) deviceSelected() {
	if g.device != nil {
		g.menuBar.RemoveAll()
		g.menuBar.Add(widget.NewLabel("Device: "), widget.NewLabel(g.device.Info.Product))
	}
	g.mainWindow.Content().Refresh()
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
