package screens

import (
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MacrosScreen struct {
	*fyne.Container
}

func NewMacrosScreen(_ *monitor.ConnectedDevice, _ *fyne.Window) NavigationItem {
	inst := &MacrosScreen{
		Container: container.NewStack(),
	}
	inst.buildContent()
	return inst
}

func (ss *MacrosScreen) buildContent() {
	ss.Container.Add(container.NewCenter(widget.NewLabel("Macros")))
}

func (ms *MacrosScreen) GetContent() *fyne.Container {
	return ms.Container
}

func (ms *MacrosScreen) Destroy() {}
