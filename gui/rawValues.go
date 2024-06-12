package gui

import (
	"log"
	"main/logger"
	"main/usb"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gousb"
)

type RawValues struct {
	container *fyne.Container
	gui       *GUI
}

func NewRawValues(gui *GUI) *RawValues {
	return &RawValues{
		gui: gui,
	}
}

func (i *RawValues) Build(device *usb.Device) *fyne.Container {
	i.container = container.NewVBox(widget.NewLabel("Profile"))
	i.readUSBValues(device)
	i.container.Hide()
	return i.container
}
func (i *RawValues) DeviceSelectionChanged(dev *usb.Device) {
	logger.Log.Infof("RawValues: Device selection changed %v", dev)
}

func (i *RawValues) Destroy() {
	if i.container != nil {
		i.container.Hide()
	}
}

func (i *RawValues) readUSBValues(device *usb.Device) {
	if device == nil {
		return
	}
	//read values coming from usb device.device
	const configNum = 1
	conf, err := device.Device.Config(configNum)
	if err != nil {
		log.Fatalf("Could not set config %v: %v", configNum, err)
	}
	logger.Log.Infof("conf: %v", conf)
	c, err := conf.Interface(3, 0)
	logger.Log.Infof("c: %v", c.Setting.Endpoints[0].Number)
	// Find the HID interrupt endpoint
	var endpoint *gousb.EndpointDesc
	for _, ep := range c.Setting.Endpoints {
		if ep.Direction == gousb.EndpointDirectionIn && ep.TransferType == gousb.TransferTypeInterrupt {
			endpoint = &ep
			break
		}
	}

	if endpoint == nil {
		log.Fatalf("Could not find an interrupt IN endpoint")
	}

	epIn, err := c.InEndpoint(endpoint.Number)
	if err != nil {
		log.Fatalf("Could not open IN endpoint %v: %v", endpoint.Number, err)
	}

	// Read data from the endpoint
	data := make([]byte, 64) // Adjust the buffer size to your device's report size
	for {
		n, err := epIn.Read(data)
		if err != nil {
			log.Printf("Read error: %v", err)
			time.Sleep(1 * time.Second) // Retry after some time
			continue
		}

		logger.Log.Infof("Read %d bytes: % X", n, data[:n])
	}
	// if err := device.Device.SetConfig(configNum); err != nil {
	//     log.Fatalf("Could not set config %v: %v", configNum, err)
	// }

}
