package screens

import (
	"log"
	resources "main/assets"
	"main/gui/widgets"
	"main/txt"
	"main/types"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/v4/process"
)

type Process struct {
	Name    string
	ExePath string
}

type SelectProcessDialog struct {
	*dialog.CustomDialog
	list            *widget.List
	body            fyne.CanvasObject
	processList     []Process
	currentSelProcs types.StringSet
	oldSelProcs     types.StringSet
	callBack        func(selection types.StringSet, confirmed bool)
	searchBar       *widgets.StringSearchBar[Process]
}

func NewSelectProcessDialog(oldSelProcs types.StringSet, parent *fyne.Window) *SelectProcessDialog {
	inst := &SelectProcessDialog{
		body:            container.NewStack(),
		processList:     []Process{},
		currentSelProcs: types.NewStringSet(),
		oldSelProcs:     oldSelProcs,
	}
	inst.CustomDialog = dialog.NewCustom(txt.GetLabel("apps.selProcTitle"), "", inst.body, *parent)

	inst.CustomDialog.Resize(fyne.NewSize(1100, 600))
	inst.buildContent()
	inst.populateProcessList()
	inst.Refresh()
	return inst
}

func (spd *SelectProcessDialog) buildContent() {
	//create confirm button
	confirmBtn := widget.NewButton(txt.GetLabel("btn.generalSelect"), func() {
		spd.CustomDialog.Hide()
		if spd.callBack != nil {
			spd.callBack(spd.currentSelProcs, true)
		}
	})

	//update confirm button
	confirmBtn.Importance = widget.HighImportance
	if spd.currentSelProcs.IsEmpty() {
		confirmBtn.Disable()
	} else {
		confirmBtn.Enable()
	}

	//set buttons
	spd.SetButtons([]fyne.CanvasObject{
		widget.NewButton(txt.GetLabel("btn.generalCancel"), func() {
			spd.CustomDialog.Hide()
			if spd.callBack != nil {
				spd.callBack(spd.currentSelProcs, false)
			}
		}), confirmBtn,
	})

	//create list
	spd.list = widget.NewList(func() int {
		return len(spd.processList)
	}, func() fyne.CanvasObject {
		return container.NewHBox(
			widget.NewCheck("", nil),
			widget.NewLabel("Name"),
			widget.NewLabel("Path"),
		)
	}, func(id int, item fyne.CanvasObject) {
		process := spd.processList[id]

		isInTheList := spd.oldSelProcs.Contains(process.ExePath)
		isSel := spd.currentSelProcs.Contains(process.ExePath)

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
				spd.currentSelProcs.Add(spd.processList[id].ExePath)
			} else {
				spd.currentSelProcs.Remove(spd.processList[id].ExePath)
			}

			if spd.currentSelProcs.IsEmpty() {
				confirmBtn.Disable()
			} else {
				confirmBtn.Enable()
			}
			check.FocusLost()
		}
	})

	//create toolbar
	spd.body.(*fyne.Container).Add(
		container.NewBorder(
			spd.getToolbar(),
			nil, nil, nil,
			spd.list,
		))
}

func (spd *SelectProcessDialog) SetOnClose(callback func(selection types.StringSet, confirmed bool)) {
	spd.callBack = callback
}

func (spd *SelectProcessDialog) getToolbar() *fyne.Container {

	spd.searchBar = widgets.NewStringSearchBar(spd.processList, func(proc Process, inputStr string) bool {
		return strings.Contains(strings.ToLower(proc.ExePath), strings.ToLower(inputStr)) || strings.Contains(strings.ToLower(proc.Name), strings.ToLower(inputStr))
	}, func(poz int) {
		spd.list.ScrollTo(poz)
		spd.list.Select(poz)
	})

	refreshButton := widget.NewButtonWithIcon(
		"", resources.ResRefreshPng,
		func() {
			spd.populateProcessList()
			spd.Refresh()
		},
	)

	entry, button := spd.searchBar.GetControls()
	return container.NewBorder(nil, nil, nil, container.NewHBox(button, refreshButton), entry)
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
	return spd.currentSelProcs
}

func (spd *SelectProcessDialog) populateProcessList() {
	spd.processList = spd.processList[:0]
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
		spd.processList = append(spd.processList, proc)
		set.Add(exe)
	}
	spd.searchBar.SetSearchList(spd.processList)
}
