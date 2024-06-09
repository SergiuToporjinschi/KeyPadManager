package usb

import (
	"github.com/google/gousb"
)

type USBHandler struct {
	device *gousb.Device
	vid    gousb.ID
	pid    gousb.ID
	ctx    *gousb.Context
}

// NewUSBHandler creates a new USB handler
func NewUSBHandler(VID gousb.ID, PID gousb.ID) *USBHandler {
	return nil
}
