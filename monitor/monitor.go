package monitor

import (
	"fmt"
	"log/slog"
	"main/devicelayout"
	"sync"
	"time"

	"github.com/google/gousb"
)

const (
	EventDeviceConnected    string = "connected"
	EventDeviceDisconnected string = "disconnected"
)

// TODO use MUTEX to protect the all maps
type USBMonitor struct {
	ctx              *gousb.Context
	stopChan         chan struct{}
	connectedDevices map[string]ConnectedDevice

	deviceEvents    chan deviceEvent
	deviceListeners map[string]func(string, *ConnectedDevice)

	monitorEvents    chan string
	monitorListeners map[string]func(string)
}

type ConnectedDevice struct {
	*gousb.Device
	*devicelayout.DeviceLayoutConfig
}

type deviceEvent struct {
	event  string
	device ConnectedDevice
}

var instance *USBMonitor
var once sync.Once

func GetInstance() *USBMonitor {
	once.Do(func() {
		instance = &USBMonitor{}
		instance.deviceEvents = make(chan deviceEvent)
		instance.deviceListeners = make(map[string]func(string, *ConnectedDevice))

		instance.monitorEvents = make(chan string)
		instance.monitorListeners = make(map[string]func(string), 10)

		instance.connectedDevices = make(map[string]ConnectedDevice)
	})
	return instance
}

func (m *USBMonitor) Start() error {
	if m.stopChan != nil {
		slog.Debug("Monitoring already started")
		return nil
	}

	slog.Debug("Starting the USB monitor")

	m.ctx = gousb.NewContext()
	m.stopChan = make(chan struct{})

	go m.eventListener()

	go m.monitorDevices()
	m.monitorEvents <- "start"

	return nil
}

func (m *USBMonitor) Stop() {
	close(m.stopChan)
	m.monitorEvents <- "stop"
	defer m.ctx.Close()
}

func (m *USBMonitor) AddMonitorEvent(name string, callback func(string)) {
	m.monitorListeners[name] = callback
}

func (m *USBMonitor) RemoveMonitorEvent(name string) {
	delete(m.monitorListeners, name)
}

func (m *USBMonitor) AddDeviceEvent(name string, callback func(string, *ConnectedDevice)) {
	m.deviceListeners[name] = callback
}

func (m *USBMonitor) RemoveDeviceEvent(name string) {
	delete(m.deviceListeners, name)
}

func (m *USBMonitor) eventListener() {
	for {
		select {

		case monitorEvent := <-m.monitorEvents:
			slog.Debug("Monitor event received", "monitorEvent", monitorEvent)
			m.callMonitorListeners(monitorEvent)
			if monitorEvent == "monitorStop" {
				return
			}
			continue
		case deviceEvent := <-m.deviceEvents:
			slog.Debug("Device event received", "deviceEvent", deviceEvent)
			m.callDeviceListeners(&deviceEvent)
			continue
		case <-m.stopChan:
			slog.Debug("Stopping event listener")
			return
		}
	}
}

func (m *USBMonitor) monitorDevices() {
	// monitorLoop:
	for {
		select {
		case <-m.stopChan:
			slog.Debug("Stopping device monitoring...")
			m.monitorEvents <- "monitorStop"
			// break monitorLoop
			return
		default:
			foundDevices, err := m.listHIDDevices()
			if err != nil {
				slog.Error("Error listing HID devices", "error", err)
				continue
			}

			// Check for new devices
			for key, val := range foundDevices {
				_, found := m.connectedDevices[key]
				if !found {
					m.connectedDevices[key] = val
					m.deviceEvents <- deviceEvent{"connected", val}
					slog.Info("Device connected", "key", key, "value", val)
				}
			}

			// Check for removed devices
			for key, dev := range m.connectedDevices {
				_, found := foundDevices[key]
				if !found {
					delete(m.connectedDevices, key)
					m.deviceEvents <- deviceEvent{"disconnected", dev}
					slog.Info("Device disconnected ", "device", dev)
				}
			}

		}
		time.Sleep(2 * time.Second)
	}
}

func (m *USBMonitor) listHIDDevices() (map[string]ConnectedDevice, error) {
	devices := make(map[string]ConnectedDevice)
	knwonDevices := devicelayout.GetInstance()
	// List devices
	devs, err := m.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		_, found := knwonDevices.FindLayout(uint16(desc.Vendor), uint16(desc.Product))
		return found
	})

	if err != nil {
		return nil, err
	}

	for _, dev := range devs {
		key, conDevice := m.newConnectedDevice(dev, knwonDevices)
		if conDevice != nil {
			devices[key] = *conDevice
		}
	}

	return devices, nil
}

func (*USBMonitor) newConnectedDevice(dev *gousb.Device, knwonDevices *devicelayout.DeviceLayout) (string, *ConnectedDevice) {
	key := fmt.Sprintf("%s/%s", dev.Desc.Vendor.String(), dev.Desc.Product.String())
	conf, found := knwonDevices.FindLayoutByKey(key)
	if found {

		prod, err := dev.Product()
		if err == nil {
			conf.Identifier.Product = prod
		} else {
			slog.Warn("Error getting product name from phisical device", "error", err)
		}

		man, err := dev.Manufacturer()
		if err == nil {
			conf.Identifier.Manufacturer = man
		} else {
			slog.Warn("Error getting manufacturer name from phisical device", "error", err)
		}

		serial, err := dev.SerialNumber()
		if err == nil {
			conf.Identifier.SerialNumber = serial
		} else {
			slog.Warn("Error getting serial number from phisical device", "error", err)
		}

		return key, &ConnectedDevice{
			Device:             dev,
			DeviceLayoutConfig: conf,
		}
	} else {
		slog.Warn("Device not found in layout config", "VID", dev.Desc.Vendor, "PID", dev.Desc.Product)
	}
	return key, nil
}

func (m *USBMonitor) callMonitorListeners(event string) {
	for _, listener := range m.monitorListeners {
		listener(event)
	}
}

func (m *USBMonitor) callDeviceListeners(event *deviceEvent) {
	for _, listener := range m.deviceListeners {
		listener(event.event, &event.device)
	}
}
