package screens

import (
	"fmt"
	"log"
	"main/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/v4/process"
)

type Process struct {
	Name    string
	ExePath string
}

type SelectProcessDialog struct {
	*dialog.ConfirmDialog
	list        *widget.List
	body        fyne.CanvasObject
	processList binding.UntypedList
	selProcs    types.StringSet
	oldSelProcs types.StringSet
}

func NewSelectProcessDialog() *SelectProcessDialog {
	inst := &SelectProcessDialog{
		body:        container.NewStack(),
		processList: binding.NewUntypedList(),
		selProcs:    types.NewStringSet(),
		oldSelProcs: types.NewStringSet(),
	}

	inst.buildContent()
	return inst
}

type MySearchToolbar struct {
	*widget.ToolbarAction
}

func (s *MySearchToolbar) ToolbarObject() fyne.CanvasObject {
	button := widget.NewButtonWithIcon("", s.ToolbarAction.Icon, s.ToolbarAction.OnActivated)
	// button.Importance = LowImportance
	entr := widget.NewEntry()
	entr.Resize(fyne.NewSize(400, 30))
	cnt := container.NewBorder(nil, nil, nil, button, entr)
	cnt.Resize(fyne.NewSize(400, 30))
	return cnt
}

func (spd *SelectProcessDialog) buildContent() {
	spd.list = widget.NewListWithData(spd.processList, func() fyne.CanvasObject {
		return container.NewHBox(
			widget.NewCheck("", nil),
			widget.NewLabel("Name"),
			widget.NewLabel("Path"),
		)
	}, func(id binding.DataItem, item fyne.CanvasObject) {
		process := getProcess(id)
		isInTheList := spd.oldSelProcs.Contains(process.ExePath)
		isSel := spd.selProcs.Contains(process.ExePath)

		check := widget.NewCheck("", nil)
		item.(*fyne.Container).Objects[0] = check
		check.Checked = isSel
		if isInTheList {
			check.Checked = true
			check.Disable()
		}
		check.Refresh()

		item.(*fyne.Container).Objects[1].(*widget.Label).SetText(process.Name)
		item.(*fyne.Container).Objects[2].(*widget.Label).SetText(process.ExePath)

		check.OnChanged = func(selected bool) {
			if selected {
				spd.selProcs.Add(getProcess(id).ExePath)
			} else {
				spd.selProcs.Remove(getProcess(id).ExePath)
			}
			check.FocusLost()
		}
	})
	toolbar := widget.NewToolbar(
		&MySearchToolbar{widget.NewToolbarAction(theme.SearchIcon(), func() { fmt.Println("search") })},
	)
	toolbar.Resize(fyne.NewSize(1800, 30))
	toolbar.Refresh()
	spd.body.(*fyne.Container).Add(container.NewBorder(toolbar, nil, nil, nil, spd.list))
}

func (spd *SelectProcessDialog) Show(oldSelProcs types.StringSet, callBack func(bool), parent *fyne.Window) {
	spd.ConfirmDialog = dialog.NewCustomConfirm("Select Process", "Select", "Cancel", spd.body, callBack, *parent)
	spd.oldSelProcs = oldSelProcs
	spd.ConfirmDialog.Resize(fyne.NewSize(1800, 600))
	spd.populateProcessList()
	spd.ConfirmDialog.Show()
}

func getProcess(id binding.DataItem) *Process {
	proc, err := id.(binding.Untyped).Get()
	if err != nil {
		log.Fatalf("Error getting process: %v", err)
	}

	rezz, ok := proc.(Process)
	if !ok {
		return nil
	}
	return &rezz
}

func (spd *SelectProcessDialog) GetSelection() types.StringSet {
	return spd.selProcs
}

func (spd *SelectProcessDialog) populateProcessList() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error getting processes: %v", err)
	}
	set := types.NewStringSet()

	// Iterate over the list of processes and print their names
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			// log.Printf("Error getting process name: %v", err)
			continue
		}
		exe, err := p.Exe()
		if err != nil || exe == "" || set.Contains(exe) {
			// log.Printf("Error getting process exe: %v", err)
			continue
		}
		proc := Process{Name: name, ExePath: exe}
		spd.processList.Append(proc)
		set.Add(exe)
	}
}
