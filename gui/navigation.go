package gui

import (
	"main/logger"
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Navigation struct {
	container.Split
	currentDevice *monitor.ConnectedDevice
}

func NewContentManager() *Navigation {
	s := &Navigation{
		Split: container.Split{
			Offset:     0.2, // Sensible default, can be overridden with SetOffset
			Horizontal: true,
			Trailing:   container.NewStack(),
		},
	}
	s.Split.Leading = NewMenuOptions(s.onMenuClicked)
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

func (c *Navigation) SetDevice(device *monitor.ConnectedDevice) {
	c.currentDevice = device
}

func (c *Navigation) onMenuClicked(content MainContent) {
	logger.Log.Debugf("Menu selection changed: device: %v; content: %v", c.currentDevice, content)
	c.Trailing.(*fyne.Container).Add(content.GetContent(c.currentDevice))
	// c.Split.Resize(fyne.NewSize(100, 100)) //TODO change size
	c.Refresh()
}
