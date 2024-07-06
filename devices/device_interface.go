package devices

import (
	"main/devicelayout"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type DeviceConstructor func(binding.Bytes, *devicelayout.DeviceDescriptor) fyne.CanvasObject
type DeviceInterface interface {
}
