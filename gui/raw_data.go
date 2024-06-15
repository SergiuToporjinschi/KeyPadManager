package gui

import (
	"main/monitor"
	"main/txt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type RawData struct {
	title    string
	navTitle string
	button   *widget.Button
	body     *fyne.Container
}

func NewRawData(onNavClick func(NavigationItem)) *RawData {

	inst := &RawData{
		title:    txt.GetLabel("navi.rawDataTitle"),
		navTitle: txt.GetLabel("navi.rawDataTitle"),
	}
	inst.buildButton(onNavClick)
	inst.buildBody()
	return inst
}

func (i *RawData) buildButton(onNavClick func(NavigationItem)) {
	i.button = widget.NewButton(i.title, func() {
		onNavClick(i)
	})
}

func (i *RawData) buildBody() {
	i.body = container.New(layout.NewFormLayout())
}

func (i *RawData) setData(dev *monitor.ConnectedDevice) {
}

func (i *RawData) GetContent(dev *monitor.ConnectedDevice) *fyne.Container {
	i.setData(dev)
	return i.body
}

func (i *RawData) GetButton() *widget.Button {
	return i.button
}

func (i *RawData) GetTitle() string {
	return i.title
}

func (i *RawData) GetNavTitle() string {
	return i.navTitle
}

func (i *RawData) Destroy() {
	//TODO
}
