package gui

import (
	"main/logger"
	"main/monitor"
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DeviceList struct {
	widget.BaseWidget
	container *fyne.Container
	combo     *widget.Select
	// deviceList  map[string]*usb.Device
	connectedDevices map[string]monitor.ConnectedDevice
	onSelection      func(*usb.Device)
}

func NewDeviceList(onSelectionChange func(*usb.Device)) *DeviceList {
	logger.Log.Debug("Building device selection widget")
	item := &DeviceList{
		onSelection: onSelectionChange,
	}
	item.container = container.NewHBox([]fyne.CanvasObject{}...)

	// item.deviceList = make(map[string]*usb.Device)
	item.combo = widget.NewSelect(item.getDeviceNameList(), item.selectionChanged)

	//set min size for the combo box
	item.container.Add(item.combo)
	item.ExtendBaseWidget(item)
	item.container.Refresh()
	return item
}

func (dl *DeviceList) getDeviceNameList() []string {
	// devListDesc := []string{}
	// for key := range dl.deviceList {
	// 	devListDesc = append(devListDesc, key)
	// }
	return []string{"Device 1", "Device 2", "Device 3"}
}

func (item *DeviceList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(item.container))
}

func (dl *DeviceList) selectionChanged(s string) {
	logger.Log.Debugf("Selected device changed: %v", s)
	// dl.onSelection(dl.deviceList[s])
}
