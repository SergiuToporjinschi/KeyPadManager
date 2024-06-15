package gui

import (
	"image/color"
	"main/monitor"
	"main/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type ContentManager struct {
	container.Split
	currentDevice *monitor.ConnectedDevice
}

func NewContentManager() *ContentManager {
	s := &ContentManager{
		Split: container.Split{
			Offset:     0.2, // Sensible default, can be overridden with SetOffset
			Horizontal: true,
			Trailing:   container.NewStack(),
		},
	}
	s.Split.Leading = NewNavigation(s.onMenuClicked)
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

func (c *ContentManager) SetDevice(device *monitor.ConnectedDevice) {
	c.currentDevice = device
}

func (c *ContentManager) onMenuClicked(navItem NavigationItem) {
	c.Trailing.(*fyne.Container).Add(
		container.NewVBox(
			newTitleText(navItem.GetTitle()),
			navItem.GetContent(c.currentDevice),
		),
	)
	c.Refresh()
}

func newTitleText(text string) *canvas.Text {
	r := utility.NewSizeableText(text, 20)
	r.TextStyle.Bold = true
	return utility.NewSizeableColorText(text, 20, color.NRGBA{R: 0xFE, G: 0x58, B: 0x62, A: 0xFF})
}
