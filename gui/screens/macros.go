package screens

import (
	"main/monitor"
	"main/txt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MacrosScreen struct {
	title     string
	navTitle  string
	button    *widget.Button
	body      *container.Scroll
	bndLength binding.ExternalInt
	stopChan  chan bool
	bndData   binding.Bytes
	onceGrid  sync.Once
}

func NewMacrosScreen() NavigationItem {
	inst := &MacrosScreen{
		title:     txt.GetLabel("navi.macrosTitle"),
		navTitle:  txt.GetLabel("navi.macrosTitle"),
		bndLength: binding.BindInt(nil),
		bndData:   binding.NewBytes(),
	}
	inst.buildBody()
	return inst
}

func (as *MacrosScreen) buildBody() {
	as.body = container.NewVScroll(container.New(layout.NewGridWrapLayout(fyne.NewSize(64, 64))))
}

func (as *MacrosScreen) GetContent(*monitor.ConnectedDevice) *container.Scroll {
	return as.body
}
func (as *MacrosScreen) Destroy() {}
