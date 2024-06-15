package gui

import (
	"main/monitor"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Gui struct {
	App                fyne.App
	MainWindow         *MainWindow
	SelectDeviceWindow *SelectDeviceWindow
	UsbMonitor         *monitor.USBMonitor
}

var once sync.Once
var instance *Gui

func GetInstance() *Gui {
	once.Do(func() {
		app := app.NewWithID("KeyPadManager")

		instance = &Gui{
			App:                app,
			UsbMonitor:         monitor.GetInstance(),
			MainWindow:         NewMainWindow(app),
			SelectDeviceWindow: NewSelectDevice(app),
		}

		instance.App.Settings().SetTheme(NewTheme(instance.App.Preferences()))
		instance.SelectDeviceWindow.AddSelectDeviceListener("Gui", instance.onDeviceSelected)
		NewSysTrayMenu(instance).Start()
	})
	return instance
}

func (g Gui) ShowDeviceWindow() {
	g.SelectDeviceWindow.Show()
}

func (g *Gui) onDeviceSelected(device *monitor.ConnectedDevice) {
	instance.SelectDeviceWindow.Close()
	instance.MainWindow.Show(device)
}
