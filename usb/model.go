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
	Info   DevInfo
	Device *gousb.Device
}
