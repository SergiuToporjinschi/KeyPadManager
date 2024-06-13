package gui

import (
	"main/monitor"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Gui struct {
	App                fyne.App
	MainWindow         fyne.Window
	SelectDeviceWindow *SelectDeviceWindow
	UsbMonitor         *monitor.USBMonitor
}

var once sync.Once
var instance *Gui

func GetInstance() *Gui {
	// myApp.NewWindow("Selected device")
	once.Do(func() {
		app := app.NewWithID("KeyPadManager")
		instance = &Gui{
			App:                app,
			UsbMonitor:         monitor.GetInstance(),
			SelectDeviceWindow: NewSelectDevice(app),
		}
		instance.MainWindow = instance.App.NewWindow("Selected device")
		instance.App.Settings().SetTheme(NewTheme(instance.App.Preferences()))
		instance.MainWindow.Hide()
		NewSysTrayMenu(instance).Start()
	})
	return instance
}

func (g Gui) ShowDeviceWindow() {
	g.SelectDeviceWindow.Show()
}

func (g Gui) ShowMainWindow(key string) {

}
