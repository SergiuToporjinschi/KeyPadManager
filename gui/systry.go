package gui

import (
	"log/slog"
	resources "main/assets"
	"main/monitor"
	"main/txt"

	"fyne.io/fyne/v2"
	"fyne.io/systray"
)

type SysTrayMenu struct {
	gui           *Gui
	mDevConMenus  map[string]*systray.MenuItem
	mSelectDevice *systray.MenuItem
	mDevConMain   *systray.MenuItem
	mQuit         *systray.MenuItem
}

func NewSysTrayMenu(gui *Gui) *SysTrayMenu {
	return &SysTrayMenu{
		gui:          gui,
		mDevConMenus: make(map[string]*systray.MenuItem),
	}
}

func (s SysTrayMenu) Start() {
	systray.RunWithExternalLoop(s.onReady, s.onExit)
}

func (s SysTrayMenu) onExit() {
	slog.Info("Exiting the application")

	//Stop usb monitoring
	s.gui.UsbMonitor.Stop()

	//Close the app
	s.gui.App.Quit()
}

func (s SysTrayMenu) onReady() {
	// Set the icon of the systray
	systray.SetIcon(resources.ResLogoIco.StaticContent) // iconData should be a byte slice containing the icon data

	// Set the tooltip of the systray
	systray.SetTooltip(txt.GetLabel("app.title"))

	// Add menu items to the systray
	s.mSelectDevice = systray.AddMenuItem(txt.GetLabel("try.selDevice"), txt.GetLabel("try.selDeviceTT"))
	s.mDevConMain = systray.AddMenuItem(txt.GetLabel("try.conDevices"), txt.GetLabel("try.conDevicesTT"))
	s.mQuit = systray.AddMenuItem(txt.GetLabel("try.quit"), txt.GetLabel("try.quitTT"))

	// Add connected devices to the menu

	s.gui.UsbMonitor.AddDeviceEvent("SysTrayConnectedDevice", s.onDeviceConnection)

	// Handle menu item clicks
	go s.sysTrayActionListener()

	//start monitor
	s.gui.UsbMonitor.Start()
}

func (s SysTrayMenu) sysTrayActionListener() {
	for {
		select {
		case <-s.mSelectDevice.ClickedCh:
			slog.Info("selDevOpen: clicked")
			s.gui.ShowDeviceWindow()
		case <-s.mQuit.ClickedCh:
			slog.Info("Quit menu systry clicked")
			systray.Quit()
			fyne.CurrentApp().Quit()
			return
		}
	}
}

func (s SysTrayMenu) onDeviceConnection(event string, device *monitor.ConnectedDevice) {
	item, found := s.mDevConMenus[device.Identifier.String()]

	if event == monitor.EventDeviceConnected {
		newItem := s.mDevConMain.AddSubMenuItem(device.Identifier.Product, device.Identifier.StringDetailed())
		s.mDevConMenus[device.Identifier.String()] = newItem

	} else if event == monitor.EventDeviceDisconnected && found {
		item.Remove()
	}
}
