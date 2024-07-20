package screens

import (
	"fmt"
	"log/slog"
	resources "main/assets"
	"main/datahandlers"
	"main/devicelayout"
	"main/gui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MappingContainer struct {
	*fyne.Container
	parent      *fyne.Window
	appSelector *widgets.ComboBox[datahandlers.Application]
	device      *devicelayout.DeviceDescriptor

	keyList             *widget.Table
	keyListSelectedCell *widget.TableCellID

	macroList   *widget.List
	scriptTable *widget.List
}

func NewMappingContainer(parent *fyne.Window, device *devicelayout.DeviceDescriptor) *MappingContainer {
	inst := &MappingContainer{
		Container: container.NewStack(),
		parent:    parent,
		device:    device,
	}
	inst.buildContent()
	return inst
}

func (mc *MappingContainer) buildContent() {
	mc.buildAppSelector()
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(resources.ResApplicationPng, mc.recordKey),
		widget.NewToolbarSpacer(),
	)
	mc.Add(container.NewBorder(container.NewVBox(toolbar, mc.appSelector), nil, nil, nil, container.NewStack(mc.getItemList())))
	mc.Refresh()
}

func (mc *MappingContainer) recordKey() {
	slog.Debug("Record key")
	slog.Debug("app", "value", mc.appSelector.GetSelectedValue())
	slog.Debug("keys", "mc.keyListSelectedCell", mc.keyListSelectedCell, "key", datahandlers.GetKeysHandlerInstance().GetByIndex(mc.keyListSelectedCell.Row))
}

func (mc *MappingContainer) buildAppSelector() {
	mc.appSelector = widgets.NewCombobox(datahandlers.GetAppHandlerInstance().GetAppList())
}

func (mc *MappingContainer) getItemList() fyne.CanvasObject {
	mc.buildKeyTable()
	mc.buildMacroTable()
	mc.buildScriptTable()

	keyList := widget.NewAccordionItem("Keys", mc.keyList)
	keyList.Open = true

	macroList := widget.NewAccordionItem("Macros", mc.macroList)
	macroList.Open = false

	scriptList := widget.NewAccordionItem("Scripts", mc.scriptTable)
	scriptList.Open = false

	acc := widget.NewAccordion(keyList, macroList, scriptList)
	acc.MultiOpen = true
	return acc
}

func (mc *MappingContainer) buildKeyTable() {
	keylist := datahandlers.GetKeysHandlerInstance()
	headers := map[int]string{0: "Name", 1: "ASCII"}

	mc.keyList = widget.NewTable( //TODO change it to select row by implementing custom
		func() (int, int) {
			return keylist.Size(), 3
		},
		func() fyne.CanvasObject {
			cell := widget.NewLabel("ID")
			return cell
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			key := keylist.GetByIndex(id.Row)

			nameCell := cell.(*widget.Label)
			if id.Col == 0 {
				nameCell.SetText(key.Name)
			} else if id.Col == 1 {
				nameCell.SetText(fmt.Sprintf("%d", key.Key))
			} else if id.Col == 2 {
				keyId := datahandlers.GetMappingHandlerInstance().GetKeyForDev(mc.device.Identifier.String(), mc.appSelector.GetSelectedValue().ExePath, key.ID)
				nameCell.SetText(fmt.Sprintf("%d", keyId))
			}
		},
	)
	mc.keyList.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Header")
	}
	mc.keyList.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		if id.Row == -1 {
			template.(*widget.Label).SetText(headers[id.Col])
		}
	}
	mc.keyList.OnSelected = func(id widget.TableCellID) {
		mc.keyListSelectedCell = &id
		slog.Debug("Selected cell", "id", id)
	}
	mc.keyList.OnUnselected = func(id widget.TableCellID) {
		mc.keyListSelectedCell = nil
		slog.Debug("UnSelected cell", "id", id)
	}
	mc.keyList.ShowHeaderRow = true
	mc.keyList.Refresh()
}

func (mc *MappingContainer) buildMacroTable() {
	mc.macroList = widget.NewList(
		func() int {
			return 10
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Path"),
			)
		},
		func(id int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText("Name")
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("Path")
		},
	)
}

func (mc *MappingContainer) buildScriptTable() {
	mc.scriptTable = widget.NewList(
		func() int {
			return 10
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Path"),
			)
		},
		func(id int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText("Name")
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("Path")
		},
	)
}

// return
