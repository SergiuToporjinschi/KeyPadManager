package screens

import (
	"context"
	"fmt"
	"log/slog"
	"main/devicelayout"
	"main/devices"
	devkeyboardGui "main/devices/devkeyboard/gui"
	"main/monitor"
	"main/utility"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gousb"
)

var EventDeviceConnected map[string]devices.DeviceConstructor = map[string]devices.DeviceConstructor{
	"6001/1000": devkeyboardGui.New,
}

type DeviceScreen struct {
	*fyne.Container
	bndData    binding.Bytes
	bndDataStr binding.String
	stopChan   chan struct{}
	closeOnce  sync.Once
}

func NewDeviceScreen(currentDevice *monitor.ConnectedDevice) NavigationItem {
	slog.Debug("Creating NewDeviceScreen")
	inst := &DeviceScreen{
		stopChan:   make(chan struct{}),
		bndData:    binding.NewBytes(),
		bndDataStr: binding.NewString(),
		Container:  container.NewStack(),
	}
	//build content
	inst.buildContent(currentDevice.DeviceDescriptor)

	//start data monitoring
	go inst.monitorUSBData(currentDevice)

	return inst
}
func (ds *DeviceScreen) GetContent() *fyne.Container {
	return ds.Container
}

func (ds *DeviceScreen) buildContent(devDescriptor *devicelayout.DeviceDescriptor) {
	//draw device layout
	instObject := EventDeviceConnected[devDescriptor.Identifier.String()](ds.bndData, devDescriptor)

	//draw console
	console := container.NewGridWrap(fyne.NewSize(400, 90), widget.NewLabelWithData(ds.bndDataStr))

	ds.Container.Add(container.NewVSplit(container.NewCenter(instObject), console))
}

func (ds *DeviceScreen) monitorUSBData(dev *monitor.ConnectedDevice) {
	slog.Debug("Starting go routine DeviceScreen")
	for {
		select {
		case <-ds.stopChan:
			slog.Debug("Stopping maonitorUSBData")
			return
		default:
			if dev == nil {
				slog.Error("Device not connected")
				return
			}
			data := readUSB(dev.Device, &dev.DeviceDescriptor.Report)
			if len(data) == 0 || len(data[1:]) == 0 {
				continue
			}
			ds.bndData.Set(data)
			txt, _ := ds.bndDataStr.Get()
			ds.bndDataStr.Set(fmt.Sprintf("%s\n%s", utility.AsBinaryString(data), txt))
		}
	}
}

func (ds *DeviceScreen) Destroy() {
	slog.Debug("Destroying DeviceScreen")
	ds.closeOnce.Do(func() {
		close(ds.stopChan)
	})
}

func readUSB(dev *gousb.Device, devReport *devicelayout.Report) []byte {
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
	data := make([]byte, devReport.Size)

	timeoutCtx, cncFunc := context.WithTimeout(context.Background(), 10*time.Microsecond) //TODO make it configurable
	defer cncFunc()

	_, err = ep.ReadContext(timeoutCtx, data)
	if gousb.TransferCancelled == err {
		return nil
	}

	if err != nil {
		slog.Error("Could not read data:", "error", err)
		return nil
	}
	return data
}
