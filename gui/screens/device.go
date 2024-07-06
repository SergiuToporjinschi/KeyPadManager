package screens

import (
	"log/slog"
	"main/devicelayout"
	"main/devices"
	devkeyboardGui "main/devices/devkeyboard/gui"
	"main/monitor"
	"main/txt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gousb"
)

var EventDeviceConnected map[string]devices.DeviceConstructor = map[string]devices.DeviceConstructor{
	"6001/1000": devkeyboardGui.New,
}

type DeviceScreen struct {
	title     string
	navTitle  string
	button    *widget.Button
	body      *container.Scroll
	bndLength binding.ExternalInt
	stopChan  chan bool
	bndData   binding.Bytes
	onceGrid  sync.Once
}

func NewDeviceScreen() NavigationItem {
	inst := &DeviceScreen{
		title:     txt.GetLabel("navi.deviceTitle"),
		navTitle:  txt.GetLabel("navi.deviceTitle"),
		bndLength: binding.BindInt(nil),
		bndData:   binding.NewBytes(),
	}
	inst.buildBody()
	return inst
}

func (ds *DeviceScreen) buildBody() {
	ds.body = container.NewVScroll(container.New(layout.NewGridWrapLayout(fyne.NewSize(64, 64))))
}

func (ds *DeviceScreen) buildBindings(devDescriptor *devicelayout.DeviceDescriptor) {
	instObject := EventDeviceConnected[devDescriptor.Identifier.String()](ds.bndData, devDescriptor)

	ds.body.Content.(*fyne.Container).Add(instObject)
}

func (ds *DeviceScreen) refreshBindings(data []byte, devDesc *devicelayout.DeviceDescriptor) {
	if len(data) == 0 || len(data[1:]) == 0 {
		return
	}
	if devDesc == nil {
		slog.Error("Device layout is not loaded")
		return
	}
	ds.bndData.Set(data)
}

func (ds *DeviceScreen) setData(dev *monitor.ConnectedDevice) {
	ds.onceGrid.Do(func() {
		ds.buildBindings(dev.DeviceDescriptor)
	})

	ds.stopChan = make(chan bool)
	go func() {
		for {
			select {
			case <-ds.stopChan:
				slog.Debug("Stopping RawData")
				return
			default:
				ds.refreshBindings(readUSB(dev.Device), dev.DeviceDescriptor)
			}
		}
	}()
}

func (ds *DeviceScreen) GetContent(dev *monitor.ConnectedDevice) *container.Scroll {
	ds.setData(dev)
	return ds.body
}

func (ds *DeviceScreen) GetButton() *widget.Button {
	return ds.button
}

func (ds *DeviceScreen) Destroy() {
	slog.Debug("Destroying RawData")
	select {
	case ds.stopChan <- true:
	default:
	}
}

func readUSB(dev *gousb.Device) []byte {
	cfg, err := dev.Config(1)
	if err != nil {
		slog.Error("Could not get config:", "error", err)
		return nil
	}
	defer cfg.Close()

	intf, err := cfg.Interface(3, 0)
	if err != nil {
		slog.Error("Could not get interface:", "error", err)
		return nil
	}

	defer intf.Close()

	// Setup the endpoint
	ep, err := intf.InEndpoint(4) // 1 is the endpoint number
	if err != nil {
		slog.Error("Could not get endpoint:", "error", err)
		return nil
	}

	// Read data from the endpoint
	data := make([]byte, 4)
	_, err = ep.Read(data)
	if err != nil {
		slog.Error("Could not read data:", "error", err)
		return nil
	}
	return data
}
