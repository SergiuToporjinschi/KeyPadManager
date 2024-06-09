package gui

import (
	"fmt"
	"main/logger"
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SelectDevice struct {
	container *fyne.Container
	gui       *GUI
}

func NewSelectDevice(gui *GUI) *SelectDevice {
	return &SelectDevice{
		gui: gui,
	}
}
func (s *SelectDevice) Build(device *usb.Device) *fyne.Container {
	deviceList, err := usb.FindDevices()
	if err != nil {
		logger.Log.Errorf("Error finding devices: %v", err)
	}
	content := container.NewHBox()
	container.NewCenter(content)

	if (len(deviceList)) == 0 {
		content.Add(widget.NewLabel("No devices found!"))
	} else {
		for _, device := range deviceList {
			content.Add(NewDeviceButton(&device, s.onSelect))
			// content.Add(NewDeviceCard(&device, s.onSelect))
		}
	}
	s.container = content
	return container.NewCenter(content)
}

func (s *SelectDevice) onSelect(device *usb.Device) {
	logger.Log.Infof("Selected device: %v", device)
	s.gui.device = device
	s.gui.deviceSelected()
}

func (s *SelectDevice) Destroy() {
	if s.container != nil {
		s.container.Hide()
	}
}

type DeviceButton struct {
	widget.Button
	Device *usb.Device
}

func NewDeviceButton(device *usb.Device, onClick func(*usb.Device)) *DeviceButton {
	button := &DeviceButton{}
	button.ExtendBaseWidget(button)
	button.SetText(fmt.Sprintf("%s - %s", device.Info.Product, device.Info.Manufacturer))
	button.Device = device
	button.OnTapped = func() { onClick(device) }
	return button
}

type DeviceCard struct {
	widget.Card
	Device *usb.Device
}

func NewDeviceCard(device *usb.Device, onClick func(*usb.Device)) *DeviceCard {
	card := &DeviceCard{}
	card.ExtendBaseWidget(card)
	card.SetTitle(fmt.Sprintf("%s - %s", device.Info.Product, device.Info.Manufacturer))
	card.SetSubTitle(device.Info.SerialNumber)
	card.Device = device
	// card.OnTapped = func() { onClick(device) }
	return card
}
