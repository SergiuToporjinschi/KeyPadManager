package usb

import (
	"fmt"
	"log/slog"

	"github.com/google/gousb"
)

const (
	VID gousb.ID = 0x6001
	PID gousb.ID = 0x1000
)

func FindDevices() ([]Device, error) {
	ctx := gousb.NewContext()
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == VID && desc.Product == PID
	})

	var deviceList = make([]Device, len(devices))
	for i, device := range devices {

		info, descriptor := getInfo(device)

		deviceList[i] = Device{
			Info:          *info,
			HIDDescriptor: descriptor,
			Device:        device,
		}
	}
	return deviceList, err
}

func getInfo(device *gousb.Device) (*DevInfo, *HIDDescriptor) {
	manufacturer, err := device.Manufacturer()
	if err != nil {
		slog.Warn("Error reading manufacturer", "error", err)
	}

	product, err := device.Product()
	if err != nil {
		slog.Warn("Error reading product", "error", err)
	}

	serialNumber, err := device.SerialNumber()
	if err != nil {
		slog.Warn("Error reading serial number", "error", err)
	}
	descriptor := getDescriptor(device)
	return &DevInfo{
			VID:          VID,
			PID:          PID,
			Manufacturer: manufacturer,
			Product:      product,
			SerialNumber: serialNumber,
		}, &HIDDescriptor{
			ReportDescriptor: descriptor,
		}
}

func getDescriptor(device *gousb.Device) []byte {
	cfg, err := device.Config(1)
	if err != nil {
		slog.Error("Could not get config", "error", err)
	}
	defer cfg.Close()

	// Iterate over all interfaces
	var interfaceNum int
	found := false

	for _, intf := range cfg.Desc.Interfaces {
		for _, setting := range intf.AltSettings {
			if setting.Class == gousb.ClassHID {
				interfaceNum = setting.Number
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		slog.Error("No HID interface found")
	}

	descriptor := make([]byte, 256)

	reqType := gousb.ControlIn | gousb.ControlInterface
	req := uint8(0x06)      // GET_DESCRIPTOR
	value := uint16(0x2200) // HID Report Descriptor
	index := uint16(interfaceNum)

	length, err := device.Control(uint8(reqType), req, value, index, descriptor)
	if err != nil {
		slog.Error("Control request failed", "error", err)
	}

	descriptor = descriptor[:length]
	parseHIDDescriptor(descriptor)
	return descriptor
}

func parseHIDDescriptor(descriptor []byte) {
	// Simplistic parser for demonstration
	i := 0
	for i < len(descriptor) {
		prefix := descriptor[i]
		size := int(prefix & 0x03)
		if size == 3 {
			size = 4
		}

		tag := prefix & 0xFC
		data := descriptor[i+1 : i+1+size]

		switch tag {
		case 0x04: // Usage Page
			fmt.Printf("%02x Usage Page: %02X\n", tag, data[0])
		case 0x08: // Usage
			fmt.Printf("%02x Usage: %02X\n", tag, data[0])
		case 0xA0: // Collection
			fmt.Printf("%02x Collection: %02X\n", tag, data[0])
		case 0xC0: // End Collection
			fmt.Printf("%02x End Collection\n", tag)
		case 0x14: // Logical Minimum
			fmt.Printf("%02x Logical Minimum: %d\n", tag, data[0])
		case 0x24: // Logical Maximum
			fmt.Printf("%02x Logical Maximum: %d\n", tag, data[0])
		case 0x74: // Report Size
			fmt.Printf("%02x Report Size: %d\n", tag, data[0])
		case 0x94: // Report Count
			fmt.Printf("%02x Report Count: %d\n", tag, data[0])
		case 0x84: // Report ID
			fmt.Printf("%02x Report ID: %d\n", tag, data[0])
		case 0x80: // Input
			fmt.Printf("%02x Input: %02X\n", tag, data[0])
		default:
			fmt.Printf("Tag: %02X Data: %v\n", tag, data)
		}

		i += 1 + size
	}
}
