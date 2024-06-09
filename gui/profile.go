package gui

import (
	"main/usb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Profile struct {
	container *fyne.Container
	gui       *GUI
}

func NewProfile(gui *GUI) *Profile {
	return &Profile{
		gui: gui,
	}
}

func (i *Profile) Init() {
}

func (i *Profile) Build(device *usb.Device) *fyne.Container {
	i.container = container.NewVBox(widget.NewLabel("Profile"))
	i.container.Hide()
	return i.container
}

func (i *Profile) Destroy() {
	if i.container != nil {
		i.container.Hide()
	}
}
