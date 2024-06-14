package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type NavOptions struct {
	widget.Accordion
}

func NewMenuOptions(onClick func(MainContent)) *NavOptions {
	menu := &NavOptions{}
	menu.ExtendBaseWidget(menu)
	devInf := NewDeviceInfo()
	inf := widget.NewButton("Info", func() { onClick(devInf) })
	// inf := widget.NewButton("Info", func() { onClick(NewDeviceInfo()) })
	//.......
	// inf := widget.NewButton("Info", func() { onClick(NewDeviceInfo()) })
	menu.Append(widget.NewAccordionItem("Device", container.NewVBox(inf)))
	return menu
}
