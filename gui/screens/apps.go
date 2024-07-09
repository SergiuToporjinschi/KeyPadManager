package screens

import (
	"log/slog"
	resources "main/assets"
	"main/devicelayout"
	"main/monitor"
	"main/types"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type AppsScreen struct {
	*fyne.Container
	bndLength     binding.ExternalInt
	bndData       binding.Bytes
	stopChan      chan struct{}
	closeOnce     sync.Once
	parentWindow  *fyne.Window
	appList       types.UniSlice[string]
	list          *widget.List
	selectedIndex *int
}

func NewAppsScreen(currentDevice *monitor.ConnectedDevice, parentWindow *fyne.Window) NavigationItem {
	inst := &AppsScreen{
		stopChan:     make(chan struct{}),
		bndLength:    binding.BindInt(nil),
		bndData:      binding.NewBytes(),
		Container:    container.NewStack(),
		parentWindow: parentWindow,
		appList:      types.NewUniSlice[string](),
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
		widget.NewToolbarAction(resources.ResTrashBinPng, as.removeApp),
		widget.NewToolbarSpacer(),
	)

	as.list = widget.NewList(
		func() int {
			return as.appList.Size()
		},
		func() fyne.CanvasObject { return widget.NewLabel("test") },
		func(i int, item fyne.CanvasObject) {
			app := as.appList.Get(i)
			if app == "" {
				return
			}
			item.(*widget.Label).SetText(app)
		})
	as.list.OnSelected = func(i int) {
		as.selectedIndex = &i
	}
	as.list.OnUnselected = func(i int) {
		as.selectedIndex = nil
	}
	as.Container.Add(container.NewBorder(toolbar, nil, nil, nil, container.NewPadded(as.list)))
}

func (as *AppsScreen) Destroy() {
	as.closeOnce.Do(func() {
		close(as.stopChan)
	})

}

func (as *AppsScreen) removeApp() {
	if as.selectedIndex != nil {
		as.appList.RemoveByIndex(*as.selectedIndex)
		as.list.Refresh()
	}
}

func (as *AppsScreen) selectExe() {
	dia := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			slog.Error("Error opening file: ", "error", err)
			return
		}
		file := strings.TrimLeft(reader.URI().String(), "file://")
		file = strings.ReplaceAll(file, "/", "\\")
		as.addApp([]string{file})
	}, *as.parentWindow)
	dia.Resize(fyne.NewSize(800, 600))
	dia.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/*"}))
	dia.Show()
}

func (as *AppsScreen) addApp(exePaths []string) {
	as.appList.AddAll(exePaths...)
	as.list.Refresh()
}

func (as *AppsScreen) selectProcess() {
	dia := NewSelectProcessDialog(types.NewStringSetWithValues(as.appList...), as.parentWindow)
	dia.SetOnClose(func(selection types.StringSet, confirmed bool) {
		if confirmed {
			slog.Debug("Selected processes: ", "list", dia.GetSelection())
			as.addApp(dia.GetSelection().Keys())
		} else {
			slog.Debug("Cancelled")
		}
	})
	dia.Show()

}
