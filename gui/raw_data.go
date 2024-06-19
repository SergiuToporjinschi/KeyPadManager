package gui

import (
	"cmp"
	"fmt"
	"image/color"
	"main/devicelayout"
	"main/logger"
	"main/monitor"
	"main/txt"
	"main/utility"
	"slices"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gousb"
)

type RawData struct {
	title     string
	navTitle  string
	button    *widget.Button
	body      *container.Scroll
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
	i.body = container.NewVScroll(container.New(layout.NewGridWrapLayout(fyne.NewSize(450, 300))))
}

func (i *RawData) newTitleText(text string) *canvas.Text {
	return utility.NewSizeableColorText(text, 15, color.NRGBA{R: 0xFE, G: 0x58, B: 0x62, A: 0xFF})
}

func (i *RawData) buildBindings(layout *devicelayout.DeviceLayoutConfig) {

	for ind, comp := range layout.Components {
		grid := container.NewGridWithColumns(4,
			widget.NewLabel(""),
			NewCustomLocaleLabel("cont.rawDataVal", &style{Bold: true}),
			NewCustomLocaleLabel("cont.rawDataHex", &style{Bold: true}),
			NewCustomLocaleLabel("cont.rawDataDec", &style{Bold: true}),
		)

		var val float64
		bndValue := controlDataValue{
			bndValueBin:        binding.BindString(nil),
			bndValueHex:        binding.BindString(nil),
			bndValueDec:        binding.BindString(nil),
			bndValueDecAsFloat: binding.BindFloat(&val),
		}

		grid.Add(NewCustomLocaleLabel("cont.rawDataReceived", &style{Bold: true}))
		// grid.Add(widget.NewLabelWithStyle(txt.GetLabel("cont.rawDataReceived"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		grid.Add(widget.NewLabelWithData(bndValue.bndValueBin))
		grid.Add(widget.NewLabelWithData(bndValue.bndValueHex))
		grid.Add(widget.NewLabelWithData(bndValue.bndValueDec))

		grid.Add(NewCustomLocaleLabel("cont.rawDataType", &style{Bold: true}))
		grid.Add(widget.NewLabel(comp.Type))
		grid.Add(widget.NewLabel(""))
		grid.Add(widget.NewLabel(""))
		grid.Add(NewCustomLocaleLabel("cont.rawDataByteNo", &style{Bold: true}))
		grid.Add(widget.NewLabel(fmt.Sprintf("%d", comp.ByteNumber)))
		grid.Add(widget.NewLabel(""))
		grid.Add(widget.NewLabel(""))
		grid.Add(NewCustomLocaleLabel("cont.rawBitMask", &style{Bold: true}))
		comp.BitPosition = strings.Trim(comp.BitPosition, " ")
		if len(comp.BitPosition) > 0 {
			grid.Add(widget.NewLabel(comp.BitPosition))
			bitPosition, _ := strconv.Atoi(comp.BitPosition)
			grid.Add(widget.NewLabel(fmt.Sprintf("%02X", bitPosition)))
			grid.Add(widget.NewLabel(fmt.Sprintf("%d", bitPosition)))
		} else {
			grid.Add(widget.NewLabel(""))
			grid.Add(widget.NewLabel(""))
			grid.Add(widget.NewLabel(""))
		}

		grid.Add(NewCustomLocaleLabel("cont.rawDataMin", &style{Bold: true}))
		grid.Add(widget.NewLabel(fmt.Sprintf("%d", comp.Min)))
		grid.Add(widget.NewLabel(""))
		grid.Add(widget.NewLabel(""))

		grid.Add(NewCustomLocaleLabel("cont.rawDataMax", &style{Bold: true}))
		grid.Add(widget.NewLabel(fmt.Sprintf("%d", comp.Max)))
		grid.Add(widget.NewLabel(""))
		grid.Add(widget.NewLabel(""))
		contrImg := container.NewStack()
		if comp.Type == "dial" {
			contrImg.Add(widget.NewProgressBarWithData(bndValue.bndValueDecAsFloat))
		} else if comp.Type == "button" {
			contrImg.Add(widget.NewButton(txt.GetLabel("cont.buttonTestLabel"), func() {}))
		}

		i.body.Content.(*fyne.Container).Add(
			container.NewBorder(
				container.NewStack(
					i.newTitleText(comp.Name),
				),
				contrImg,
				nil,
				nil,
				grid),
		)

		if ind < len(layout.Components)-1 {
			widget.NewSeparator()
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

func (i *RawData) GetContent(dev *monitor.ConnectedDevice) *container.Scroll {
	i.setData(dev)
	return i.body
}

func (i *RawData) GetButton() *widget.Button {
	return i.button
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
