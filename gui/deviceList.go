package gui

import (
	"fmt"
	"main/logger"
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DeviceList struct {
	widget.BaseWidget
	container   *fyne.Container
	btn         *widget.Button
	combo       *widget.Select
	deviceList  map[string]*usb.Device
	onSelection func(*usb.Device)
}

func NewDeviceList(onSelectionChange func(*usb.Device)) *DeviceList {
	logger.Log.Debug("Building device selection widget")
	item := &DeviceList{
		onSelection: onSelectionChange,
	}
	item.container = container.NewHBox([]fyne.CanvasObject{}...)

	item.btn = widget.NewButton("Refresh", item.searchDevices)
	item.container.Add(item.btn)

	item.deviceList = make(map[string]*usb.Device)
	item.combo = widget.NewSelect(item.getDeviceNameList(), item.selectionChanged)

	//set min size for the combo box
	item.container.Add(item.combo)
	item.ExtendBaseWidget(item)
	item.container.Refresh()
	return item
}

func (dl *DeviceList) getDeviceNameList() []string {
	devListDesc := []string{}
	for key := range dl.deviceList {
		devListDesc = append(devListDesc, key)
	}
	return devListDesc
}

func (item *DeviceList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(item.container))
}

func (dl *DeviceList) selectionChanged(s string) {
	logger.Log.Debugf("Selected device changed: %v", s)
	dl.onSelection(dl.deviceList[s])
}

func (dl *DeviceList) searchDevices() {
	logger.Log.Debug("Searching for devices")

	for k := range dl.deviceList {
		delete(dl.deviceList, k)
	}

	if err := dl.updateDeviceList(); err != nil {
		logger.Log.Errorf("Error updatsing device list: %v", err)
		return
	}
	dl.combo.Options = dl.getDeviceNameList()
	dl.combo.Refresh()
}

func (dl *DeviceList) updateDeviceList() error {
	deviceList, err := usb.FindDevices()
	if err != nil {
		logger.Log.Errorf("Error finding devices: %v", err)
		return err
	}
	if len(deviceList) == 0 {
		return nil
	}
	for _, device := range deviceList {
		key := fmt.Sprintf("%s %s", device.Info.Manufacturer, device.Info.Product)
		dl.deviceList[key] = &device
	}
	return nil
}
