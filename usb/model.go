package usb

import "github.com/google/gousb"

type DevInfo struct {
	VID          gousb.ID
	PID          gousb.ID
	Manufacturer string
	Product      string
	SerialNumber string
}

type Device struct {
	Info          DevInfo
	Device        *gousb.Device
	HIDDescriptor *HIDDescriptor
}

type HIDDescriptor struct {
	ReportDescriptor []byte
	UsagePage        uint16
	Usage            uint16
	LogicalMin       int16
	LogicalMax       int16
	ReportSize       uint16
	ReportCount      uint16
	ReportID         uint16
}
