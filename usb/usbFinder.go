package usb

import (
	"main/logger"

	"github.com/google/gousb"
)

const (
	VID gousb.ID = 0x239A
	PID gousb.ID = 0x80F4
)

func FindDevices() ([]Device, error) {
	ctx := gousb.NewContext()
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == VID && desc.Product == PID
	})

	var deviceList = make([]Device, len(devices))
	for i, device := range devices {
		deviceList[i] = Device{
			Info:   *getInfo(device),
			Device: device,
		}
	}
	return deviceList, err
}

func getInfo(device *gousb.Device) *DevInfo {
	manufacturer, err := device.Manufacturer()
	if err != nil {
		logger.Log.Warnf("Error reading manufacturer: %v\n", err)
	}

	product, err := device.Product()
	if err != nil {
		logger.Log.Warnf("Error reading product: %v\n", err)
	}

	serialNumber, err := device.SerialNumber()
	if err != nil {
		logger.Log.Warnf("Error reading serial number: %v\n", err)
	}

	return &DevInfo{
		VID:          VID,
		PID:          PID,
		Manufacturer: manufacturer,
		Product:      product,
		SerialNumber: serialNumber,
	}
}
