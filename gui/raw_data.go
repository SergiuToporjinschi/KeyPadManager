package gui

import (
	"cmp"
	"fmt"
	"main/devicelayout"
	"main/logger"
	"main/monitor"
	"main/txt"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gousb"
)

type RawData struct {
	title           string
	navTitle        string
	button          *widget.Button
	body            *fyne.Container
	bndData         binding.ExternalString
	bndLength       binding.ExternalInt
	bndColumnHeader []binding.ExternalString
	bndRowHeader    []binding.ExternalString
	stopChan        chan bool
}

func NewRawData(onNavClick func(NavigationItem)) *RawData {
	inst := &RawData{
		title:           txt.GetLabel("navi.rawDataTitle"),
		navTitle:        txt.GetLabel("navi.rawDataTitle"),
		bndData:         binding.BindString(nil),
		bndLength:       binding.BindInt(nil),
		bndColumnHeader: make([]binding.ExternalString, 4),
		bndRowHeader:    make([]binding.ExternalString, 6),
	}
	for i := range inst.bndColumnHeader {
		inst.bndColumnHeader[i] = binding.BindString(nil)
	}
	for i := range inst.bndRowHeader {
		inst.bndRowHeader[i] = binding.BindString(nil)
	}
	inst.buildButton(onNavClick)
	inst.buildBody()
	return inst
}

func (i *RawData) buildButton(onNavClick func(NavigationItem)) {
	i.button = widget.NewButton(i.title, func() {
		onNavClick(i)
	})
}

func (i *RawData) buildBody() {
	t := widget.NewTable(
		func() (rows int, cols int) {
			return len(i.bndRowHeader), len(i.bndColumnHeader)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},

		func(cell widget.TableCellID, o fyne.CanvasObject) {

		},
	)

	t.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("HHH")
	}
	t.ShowHeaderColumn = true
	t.ShowHeaderRow = true
	// t.StickyColumnCount = 1
	t.StickyRowCount = 1
	t.UpdateHeader = func(cell widget.TableCellID, o fyne.CanvasObject) {
		if cell.Row == -1 {
			val, err := i.bndColumnHeader[cell.Col].Get()
			if err != nil {
				logger.Log.Warn("Error getting column header", val)
			}
			o.(*widget.Label).SetText(val)
		} else if cell.Col == -1 {
			val, err := i.bndRowHeader[cell.Row].Get()
			if err != nil {
				logger.Log.Warn("Error getting row header")
			}
			o.(*widget.Label).SetText(val)
		}
	}

	i.body = container.NewStack(t)
}

func (i *RawData) setData(dev *monitor.ConnectedDevice) {
	slices.SortFunc(dev.DeviceLayoutConfig.Components, func(a, b devicelayout.Component) int {
		if n := cmp.Compare(a.ByteNumber, b.ByteNumber); n != 0 {
			return n
		}
		return 0
	})

	i.bndColumnHeader[0].Set(txt.GetLabel(""))
	i.bndColumnHeader[2].Set(txt.GetLabel("cont.rawDataVal"))
	i.bndColumnHeader[1].Set(txt.GetLabel("cont.rawDataHex"))
	i.bndColumnHeader[3].Set(txt.GetLabel("cont.rawDataDec"))

	i.bndRowHeader[0].Set(txt.GetLabel(""))
	i.bndRowHeader[1].Set(txt.GetLabel("cont.rawDataReceived"))
	i.bndRowHeader[2].Set(txt.GetLabel("cont.rawDataType"))
	i.bndRowHeader[3].Set(txt.GetLabel("cont.rawDataBytePos"))
	i.bndRowHeader[4].Set(txt.GetLabel("cont.rawDataMin"))
	i.bndRowHeader[5].Set(txt.GetLabel("cont.rawDataMax"))
	i.body.Refresh()
	// dev.DeviceLayoutConfig.Components[cell.Row].Name

	// for _, v := range dev.DeviceLayoutConfig.Components {
	// 	i.body.Add(utility.NewTitleLabel(v.Name))
	// 	i.body.Add(utility.NewTitleLabel(v.Type)) //TODO Maybe mapping it to a better name
	// 	i.body.Add()
	// }
	// widget.NewTableWithHeaders()
	// i.body.Add(widget.NewLabelWithData(i.bindingData))
	i.stopChan = make(chan bool)
	go func() {
		for {
			select {
			case <-i.stopChan:
				logger.Log.Debug("Stopping RawData")
				return
			default:
				data := readUSB(dev.Device)
				fmt.Printf("s %v\n", data)
			}
			// i.bndLength.Set(len(dev.DeviceLayoutConfig.Components))
			// i.bndData.Set()

		}
	}()
}

func (i *RawData) GetContent(dev *monitor.ConnectedDevice) *fyne.Container {
	i.setData(dev)
	return i.body
}

func (i *RawData) GetButton() *widget.Button {
	return i.button
}

func (i *RawData) GetTitle() string {
	return i.title
}

func (i *RawData) GetNavTitle() string {
	return i.navTitle
}

func (i *RawData) Destroy() {
	logger.Log.Debug("Destroying RawData")
	i.stopChan <- true
}

func readUSB(dev *gousb.Device) []byte {
	cfg, err := dev.Config(1)
	if err != nil {
		fmt.Println("Could not get config:", err)
		return nil
	}
	defer cfg.Close()

	intf, err := cfg.Interface(3, 0)
	if err != nil {
		fmt.Println("Could not get interface:", err)
		return nil
	}

	defer intf.Close()

	// Setup the endpoint
	ep, err := intf.InEndpoint(4) // 1 is the endpoint number
	if err != nil {
		fmt.Println("Could not set up endpoint:", err)
		return nil
	}

	// Read data from the endpoint
	data := make([]byte, 5) // 64 is the size of the data buffer
	_, err = ep.Read(data)
	if err != nil {
		fmt.Println("Could not read data:", err)
		return nil
	}
	return data
}
