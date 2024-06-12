package monitoring

import (
	"fmt"
	"main/devicelayout"
	"main/logger"
	"sync"
	"time"

	"github.com/google/gousb"
)

type Monitor struct {
	ctx              *gousb.Context
	stopChan         chan struct{}
	connectedDevices map[string]ConnectedDevice
}

type ConnectedDevice struct {
	*gousb.Device
	*devicelayout.DeviceLayoutConfig
}

var instance *Monitor
var once sync.Once

func GetInstance() *Monitor {
	once.Do(func() {
		instance = &Monitor{}
		instance.connectedDevices = make(map[string]ConnectedDevice)
	})
	return instance
}

func (m *Monitor) Start() error {
	if m.stopChan != nil {
		logger.Log.Debug("Monitoring already started")
		return nil
	}

	logger.Log.Debug("Starting the USB monitor")

	m.ctx = gousb.NewContext()

	m.stopChan = make(chan struct{})

	go m.monitorDevices()

	return nil
}

func (m *Monitor) Stop() {
	close(m.stopChan)
	defer m.ctx.Close()
}

func (m *Monitor) monitorDevices() {

	for {
		select {
		case <-m.stopChan:
			logger.Log.Debug("Stopping device monitoring...")
			return
		default:
			foundDevices, err := m.listHIDDevices()
			if err != nil {
				logger.Log.Errorf("Error listing HID devices: %v", err)
				continue
			}

			// Check for new devices
			for key, val := range foundDevices {
				_, found := m.connectedDevices[key]
				if !found {
					m.connectedDevices[key] = val
					logger.Log.Infof("Device connected: %s %v", key, val)
				}
			}

			// Check for removed devices
			for key, dev := range m.connectedDevices {
				_, found := foundDevices[key]
				if !found {
					logger.Log.Infof("Device disconnected: %v", dev)
					delete(m.connectedDevices, key)
				}
			}

		}
		time.Sleep(2 * time.Second)
	}
}
func (m *Monitor) listHIDDevices() (map[string]ConnectedDevice, error) {
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
		key := fmt.Sprintf("%s/%s", dev.Desc.Vendor.String(), dev.Desc.Product.String())
		conf, found := knwonDevices.FindLayoutByKey(key)
		if found {
			devices[key] = ConnectedDevice{
				Device:             dev,
				DeviceLayoutConfig: conf,
			}
		} else {
			logger.Log.Warnf("Device with VID: %v and PID: %v not found in device layout config", dev.Desc.Vendor, dev.Desc.Product)
		}

	}

	return devices, nil
}
