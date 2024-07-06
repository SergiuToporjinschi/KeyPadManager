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

type ScriptsScreen struct {
	title     string
	navTitle  string
	button    *widget.Button
	body      *container.Scroll
	bndLength binding.ExternalInt
	stopChan  chan bool
	bndData   binding.Bytes
	onceGrid  sync.Once
}

func NewScriptsScreen() NavigationItem {
	inst := &ScriptsScreen{
		title:     txt.GetLabel("navi.scriptTitle"),
		navTitle:  txt.GetLabel("navi.scriptTitle"),
		bndLength: binding.BindInt(nil),
		bndData:   binding.NewBytes(),
	}
	inst.buildBody()
	return inst
}

func (ss *ScriptsScreen) buildBody() {
	ss.body = container.NewVScroll(container.New(layout.NewGridWrapLayout(fyne.NewSize(64, 64))))
}

func (ss *ScriptsScreen) GetContent(*monitor.ConnectedDevice) *container.Scroll {
	return ss.body
}
func (ss *ScriptsScreen) Destroy() {}
