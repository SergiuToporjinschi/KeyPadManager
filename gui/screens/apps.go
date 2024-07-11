package screens

import (
	"encoding/json"
	"log/slog"
	resources "main/assets"
	"main/devicelayout"
	"main/monitor"
	"main/types"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Application struct {
	Name    string
	ExePath string
}

func (a *Application) Key() string {
	return a.ExePath
}

func (a *Application) Value() Application {
	return *a
}

type AppsScreen struct {
	*fyne.Container
	bndLength     binding.ExternalInt
	bndData       binding.Bytes
	stopChan      chan struct{}
	closeOnce     sync.Once
	parentWindow  *fyne.Window
	appList       types.UniAnySlice[string, Application]
	list          *widget.List
	selectedIndex *int
	appListMutex  sync.Mutex
}

func NewAppsScreen(currentDevice *monitor.ConnectedDevice, parentWindow *fyne.Window) NavigationItem {
	inst := &AppsScreen{
		stopChan:     make(chan struct{}),
		bndLength:    binding.BindInt(nil),
		bndData:      binding.NewBytes(),
		Container:    container.NewStack(),
		parentWindow: parentWindow,
		appList:      loadAppList(),
	}
	inst.buildContent(currentDevice.DeviceDescriptor)
	return inst
}

func (as *AppsScreen) GetContent() *fyne.Container {
	return as.Container
}
func (c *AppsScreen) Dragged(event *fyne.DragEvent) {
	// Implement logic during dragging over the container if needed
}

// DragEnd is called when the drag event ends
func (c *AppsScreen) DragEnd() {
	// Implement drop logic here
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Drop Event",
		Content: "Item dropped on custom container",
	})
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
		func() fyne.CanvasObject { return container.NewHBox(widget.NewLabel("name"), widget.NewLabel("test")) },
		func(i int, item fyne.CanvasObject) {
			app := as.appList.GetByIndex(i).Value()
			if app.ExePath == "" {
				return
			}
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(app.Name)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(app.ExePath)
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
	as.appListMutex.Lock()
	defer as.appListMutex.Unlock()
	if as.selectedIndex != nil {
		as.appList.RemoveByIndex(*as.selectedIndex)
		as.list.Refresh()
		go as.saveAppList()
	}
}

func (as *AppsScreen) addApp(exePaths []string) {
	as.appListMutex.Lock()
	defer as.appListMutex.Unlock()
	for _, exePath := range exePaths {
		fileName := filepath.Base(exePath)
		ext := filepath.Ext(fileName)
		nameWithoutExt := strings.TrimSuffix(fileName, ext)
		as.appList.Add(&Application{Name: nameWithoutExt, ExePath: exePath})
	}
	as.list.Refresh()
	go as.saveAppList()
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
	dia.SetFilter((storage.NewExtensionFileFilter([]string{".exe", ".com"})))
	dia.Show()
}

func (as *AppsScreen) selectProcess() {
	dia := NewSelectProcessDialog(types.NewStringSetWithValues(as.appList.Keys()...), as.parentWindow)
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

func loadAppList() types.UniAnySlice[string, Application] {
	var appList = types.NewUniAnySlice[string, Application]()
	content, err := os.ReadFile("apps.json")
	if os.IsNotExist(err) {
		slog.Warn("File not found: ", "error", err)
		return appList
	} else if err != nil {
		slog.Error("Error opening file: ", "error", err)
		return appList
	}
	var apps []Application
	err = json.Unmarshal(content, &apps)
	if err != nil {
		slog.Error("Error unmarshalling app list: ", "error", err)
		return appList
	}
	pairs := make([]types.PairKeyValue[string, Application], len(apps))
	for i, app := range apps {
		pairs[i] = &app
	}
	appList.AddAll(pairs...)
	return appList

}

func (as *AppsScreen) saveAppList() {
	slog.Debug("Saving app list started")
	as.appListMutex.Lock()
	defer as.appListMutex.Unlock()
	jsonApps, err := json.MarshalIndent(as.appList, "", "  ")
	if err != nil {
		slog.Error("Error marshalling app list: ", "error", err)
		return
	}

	file, err := os.OpenFile("apps.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Error opening file: ", "error", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(string(jsonApps))
	if err != nil {
		slog.Error("Error writing to file: ", "error", err)
		return
	}
	slog.Debug("Saving app list finished")
}
