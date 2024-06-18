package gui

import (
	"cmp"
	"fmt"
	"main/devicelayout"
	"main/logger"
	"main/monitor"
	"main/txt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gousb"
)

type RawData struct {
	title     string
	navTitle  string
	button    *widget.Button
	body      *fyne.Container
	grid      *fyne.Container
	bndLength binding.ExternalInt
	stopChan  chan bool
	bndDta    []controlDataValue
	onceGrid  sync.Once
}

type controlDataValue struct {
	bndValueBin        binding.ExternalString
	bndValueHex        binding.ExternalString
	bndValueDec        binding.ExternalString
	bndValueDecAsFloat binding.Float
}

func NewRawData() NavigationItem {
	inst := &RawData{
		title:     txt.GetLabel("navi.rawDataTitle"),
		navTitle:  txt.GetLabel("navi.rawDataTitle"),
		bndLength: binding.BindInt(nil),
	}
	inst.buildBody()
	return inst
}

func (i *RawData) buildBody() {
	i.grid = container.NewGridWithColumns(4,
		widget.NewLabel(""),
		widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataVal"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataHex"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataDec"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	i.body = container.NewHBox(container.NewVBox(i.grid))
}

func (i *RawData) buildBindings(layout *devicelayout.DeviceLayoutConfig) {

	for ind, comp := range layout.Components {
		var val float64
		bndValue := controlDataValue{
			bndValueBin:        binding.BindString(nil),
			bndValueHex:        binding.BindString(nil),
			bndValueDec:        binding.BindString(nil),
			bndValueDecAsFloat: binding.BindFloat(&val),
		}

		i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataReceived"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		i.grid.Add(widget.NewLabelWithData(bndValue.bndValueBin))
		i.grid.Add(widget.NewLabelWithData(bndValue.bndValueHex))
		i.grid.Add(widget.NewLabelWithData(bndValue.bndValueDec))

		i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataType"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		i.grid.Add(widget.NewLabel(comp.Type))
		i.grid.Add(widget.NewLabel(""))
		i.grid.Add(widget.NewLabel(""))

		i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataByteNo"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		i.grid.Add(widget.NewLabel(fmt.Sprintf("%d", comp.ByteNumber)))
		i.grid.Add(widget.NewLabel(""))
		i.grid.Add(widget.NewLabel(""))

		i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawBitMask"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		comp.BitPosition = strings.Trim(comp.BitPosition, " ")
		if len(comp.BitPosition) > 0 {
			i.grid.Add(widget.NewLabel(comp.BitPosition))
			bitPosition, _ := strconv.Atoi(comp.BitPosition)
			i.grid.Add(widget.NewLabel(fmt.Sprintf("%02X", bitPosition)))
			i.grid.Add(widget.NewLabel(fmt.Sprintf("%d", bitPosition)))
		} else {
			i.grid.Add(widget.NewLabel(""))
			i.grid.Add(widget.NewLabel(""))
			i.grid.Add(widget.NewLabel(""))
		}

		i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataMin"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		i.grid.Add(widget.NewLabel(fmt.Sprintf("%d", comp.Min)))
		i.grid.Add(widget.NewLabel(""))
		i.grid.Add(widget.NewLabel(""))

		i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataMax"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		i.grid.Add(widget.NewLabel(fmt.Sprintf("%d", comp.Max)))
		i.grid.Add(widget.NewLabel(""))
		i.grid.Add(widget.NewLabel(""))

		if comp.Type == "dial" {
			i.grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.progress"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			i.grid.Add(widget.NewProgressBarWithData(bndValue.bndValueDecAsFloat))
			i.grid.Add(widget.NewLabel(""))
			i.grid.Add(widget.NewLabel(""))
		}

		if ind < len(layout.Components)-1 {
			i.grid.Add(widget.NewSeparator())
			i.grid.Add(widget.NewSeparator())
			i.grid.Add(widget.NewSeparator())
			i.grid.Add(widget.NewSeparator())
		}

		i.bndDta = append(i.bndDta, bndValue)
	}
}

func (i *RawData) refreshBindings(data []byte) {
	dataWithoutReportID := data[1:]
	for inx, val := range dataWithoutReportID {
		i.bndDta[inx].bndValueBin.Set(fmt.Sprintf("%08b", val))
		i.bndDta[inx].bndValueHex.Set(fmt.Sprintf("%02X", val))
		i.bndDta[inx].bndValueDec.Set(fmt.Sprintf("%d", val))
		i.bndDta[inx].bndValueDecAsFloat.Set(float64(val) / 100)
	}
}

func (i *RawData) setData(dev *monitor.ConnectedDevice) {
	slices.SortFunc(dev.DeviceLayoutConfig.Components, func(a, b devicelayout.Component) int {
		if n := cmp.Compare(a.ByteNumber, b.ByteNumber); n != 0 {
			return n
		}
		return 0
	})

	i.onceGrid.Do(func() {
		i.buildBindings(dev.DeviceLayoutConfig)
	})

	i.stopChan = make(chan bool)
	go func() {
		for {
			select {
			case <-i.stopChan:
				logger.Log.Debug("Stopping RawData")
				return
			default:
				i.refreshBindings(readUSB(dev.Device))
			}
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
