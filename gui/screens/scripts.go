package screens

import (
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ScriptsScreen struct {
	*fyne.Container
}

func NewScriptsScreen(_ *monitor.ConnectedDevice, _ *fyne.Window) NavigationItem {
	inst := &ScriptsScreen{
		Container: container.NewStack(),
	}
	inst.buildContent()
	return inst
}

func (ss *ScriptsScreen) buildContent() {
	ss.Container.Add(container.NewCenter(widget.NewLabel("Scripts")))
}

func (ss *ScriptsScreen) GetContent() *fyne.Container {
	return ss.Container
}
func (ss *ScriptsScreen) Destroy() {}
