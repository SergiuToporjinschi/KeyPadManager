package screens

import (
	"main/devicelayout"
	"main/monitor"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type AppsScreen struct {
	*fyne.Container
	bndLength binding.ExternalInt
	bndData   binding.Bytes
	stopChan  chan struct{}
	closeOnce sync.Once
}

func NewAppsScreen(currentDevice *monitor.ConnectedDevice) NavigationItem {
	inst := &AppsScreen{
		stopChan:  make(chan struct{}),
		bndLength: binding.BindInt(nil),
		bndData:   binding.NewBytes(),
		Container: container.NewStack(),
	}
	inst.buildContent(currentDevice.DeviceDescriptor)
	return inst
}

func (as *AppsScreen) GetContent() *fyne.Container {
	return as.Container
}

func (as *AppsScreen) buildContent(_ *devicelayout.DeviceDescriptor) {
	as.Container.Add(container.NewCenter(widget.NewLabel("Apps")))
}

func (as *AppsScreen) Destroy() {
	as.closeOnce.Do(func() {
		close(as.stopChan)
	})
}
