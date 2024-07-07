package screens

import (
	"fmt"
	"log/slog"
	resources "main/assets"
	"main/devicelayout"
	"main/monitor"
	"main/types"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppsScreen struct {
	*fyne.Container
	bndLength    binding.ExternalInt
	bndData      binding.Bytes
	stopChan     chan struct{}
	closeOnce    sync.Once
	parentWindow *fyne.Window
	appList      types.StringSet
}

func NewAppsScreen(currentDevice *monitor.ConnectedDevice, parentWindow *fyne.Window) NavigationItem {
	inst := &AppsScreen{
		stopChan:     make(chan struct{}),
		bndLength:    binding.BindInt(nil),
		bndData:      binding.NewBytes(),
		Container:    container.NewStack(),
		parentWindow: parentWindow,
		appList:      types.NewStringSet(),
	}
	inst.buildContent(currentDevice.DeviceDescriptor)
	return inst
}

func (as *AppsScreen) GetContent() *fyne.Container {

	return as.Container
}

func (as *AppsScreen) buildContent(_ *devicelayout.DeviceDescriptor) {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(resources.ResWindowsPng, as.selectProcess),
		widget.NewToolbarAction(resources.ResSearchExePng, as.selectExe),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(resources.ResTrashBinPng, as.removeApp),
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() { fmt.Println("New") }),
		widget.NewToolbarAction(theme.AccountIcon(), func() { fmt.Println("New") }),
		widget.NewToolbarSpacer(),
	)

	as.Container.Add(container.NewBorder(toolbar, nil, nil, nil, nil))
}

func (as *AppsScreen) Destroy() {
	as.closeOnce.Do(func() {
		close(as.stopChan)
	})
}

func (as *AppsScreen) removeApp() {

}

func (as *AppsScreen) selectExe() {
	dia := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			slog.Error("Error opening file: ", "error", err)
			return
		}
		as.addApp(reader.URI().String())
	}, *as.parentWindow)
	dia.Resize(fyne.NewSize(800, 600))
	dia.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/*"}))
	dia.Show()
}

func (as *AppsScreen) addApp(exePath string) {
	slog.Debug("Adding app: ", "path", exePath)
}

func (as *AppsScreen) selectProcess() {
	dia := NewSelectProcessDialog()
	dia.Show(as.appList, func(selected bool) {
		as.appList.AddAll(dia.GetSelection().Keys()...)
	}, as.parentWindow)

	// listAllRunningProcesses()
}
