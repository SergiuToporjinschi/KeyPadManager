package gui

import (
	"main/logger"
	"main/monitor"
	"os"

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
	logger.Log.Info("Exiting the application")

	//Stop usb monitoring
	s.gui.UsbMonitor.Stop()

	//Close the app
	s.gui.App.Quit()
}

func (s SysTrayMenu) onReady() {
	data, err := s.getIcon()
	if err != nil {
		logger.Log.Warnf("Error reading icon file for systray: %v", err)
	} else {
		// Set the icon of the systray
		systray.SetIcon(data) // iconData should be a byte slice containing the icon data
	}

	// Set the tooltip of the systray
	systray.SetTooltip("Keypad Manager")

	// Add menu items to the systray
	s.mSelectDevice = systray.AddMenuItem("Select device", "Open device selection window")
	s.mDevConMain = systray.AddMenuItem("Connected devices", "Show connected devices")
	s.mQuit = systray.AddMenuItem("Quit", "Quit the application")

	// Add connected devices to the menu

	s.gui.UsbMonitor.AddDeviceEvent("SysTrayConnectedDevice", s.onDeviceConnection)

	// Handle menu item clicks
	go s.sysTrayActionListener()

	//start monitor
	s.gui.UsbMonitor.Start()
}

func (SysTrayMenu) getIcon() ([]byte, error) {
	//TODO move it in resources
	data, err := os.ReadFile("icon.ico")
	if err != nil {
		return data, err
	}
	return data, nil
}

func (s SysTrayMenu) sysTrayActionListener() {
	for {
		select {
		case <-s.mSelectDevice.ClickedCh:
			logger.Log.Info("selDevOpen: clicked")
			s.gui.ShowDeviceWindow()
		case <-s.mQuit.ClickedCh:
			logger.Log.Info("Quit menu systry clicked")
			systray.Quit()
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
