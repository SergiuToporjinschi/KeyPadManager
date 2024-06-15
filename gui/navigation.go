package gui

import (
	"main/monitor"
	"main/txt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type NavigationItem interface {
	GetTitle() string
	GetNavTitle() string
	GetButton() *widget.Button
	GetContent(*monitor.ConnectedDevice) *fyne.Container
	Destroy()
}
type Navigation struct {
	widget.Accordion
}

func NewNavigation(onClick func(NavigationItem)) *Navigation {
	nav := &Navigation{}
	nav.ExtendBaseWidget(nav)

	inf := NewDeviceInfo(onClick)
	rawData := NewRawData(onClick)

	nav.Append(
		widget.NewAccordionItem(txt.GetLabel("navi.deviceSecTitle"),
			container.NewVBox(
				inf.GetButton(),
				rawData.GetButton(),
			),
		),
	)
	return nav
}
