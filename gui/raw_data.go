package gui

import (
	"fmt"
	"image/color"
	"log/slog"
	resources "main/assets"
	"main/devicelayout"
	"main/gui/widgets"
	"main/monitor"
	"main/txt"
	"main/utility"
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
	bndData   []controlBindings
	onceGrid  sync.Once
}

type controlBindings struct {
	bndValueBin     binding.ExternalString
	bndValueHex     binding.ExternalString
	bndValueDec     binding.ExternalString
	bndValueFocused binding.Bool
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

func (rd *RawData) buildBody() {
	rd.body = container.NewVScroll(container.New(layout.NewGridWrapLayout(fyne.NewSize(450, 180))))
}

func (rd *RawData) newTitleText(text string) *canvas.Text {
	return utility.NewSizeableColorText(text, 15, color.NRGBA{R: 0xFE, G: 0x58, B: 0x62, A: 0xFF})
}

func (rd *RawData) buildBindings(layout *devicelayout.DeviceLayoutConfig) {
	for ind, comp := range layout.Components {
		var grid *fyne.Container
		var bindings controlBindings
		if comp.Type == "button" {
			grid, bindings = rd.buildButtonInfoGrid(comp)
		} else if comp.Type == "encoder" {
			grid, bindings = rd.buildEncoderInfoGrid(comp)
		}
		rd.bndData = append(rd.bndData, bindings)
		rd.body.Content.(*fyne.Container).Add(
			container.NewBorder(
				container.NewStack(
					rd.newTitleText(fmt.Sprintf("%s (%s)", comp.Name, comp.Type)),
				),
				nil,
				nil,
				nil,
				grid),
		)

		if ind < len(layout.Components)-1 {
			widget.NewSeparator()
		}
	}
}

func (rd *RawData) buildButtonInfoGrid(_ devicelayout.Component) (*fyne.Container, controlBindings) {
	grid := container.NewGridWithColumns(4,
		widget.NewLabel(""),
		NewCustomLocaleLabel("cont.rawDataVal", &style{Bold: true}),
		NewCustomLocaleLabel("cont.rawDataHex", &style{Bold: true}),
		NewCustomLocaleLabel("cont.rawDataDec", &style{Bold: true}),
	)
	bndValue := controlBindings{
		bndValueBin:     binding.BindString(nil),
		bndValueHex:     binding.BindString(nil),
		bndValueDec:     binding.BindString(nil),
		bndValueFocused: binding.BindBool(nil),
	}
	grid.Add(NewCustomLocaleLabel("cont.rawDataReceived", &style{Bold: true}))
	grid.Add(widget.NewLabelWithData(bndValue.bndValueBin))
	grid.Add(widget.NewLabelWithData(bndValue.bndValueHex))
	grid.Add(widget.NewLabelWithData(bndValue.bndValueDec))

	keySwitch := widgets.NewKeySwitchControl(resources.ResButtonGrayPng, resources.ResButtonPng, bndValue.bndValueFocused)
	return container.NewVBox(grid, keySwitch), bndValue
}

func (rd *RawData) buildEncoderInfoGrid(_ devicelayout.Component) (*fyne.Container, controlBindings) {
	grid := container.NewGridWithColumns(4,
		widget.NewLabel(""),
		NewCustomLocaleLabel("cont.rawDataVal", &style{Bold: true}),
		NewCustomLocaleLabel("cont.rawDataHex", &style{Bold: true}),
		NewCustomLocaleLabel("cont.rawDataDec", &style{Bold: true}),
	)

	bindings := controlBindings{
		bndValueBin:     binding.BindString(nil),
		bndValueHex:     binding.BindString(nil),
		bndValueDec:     binding.BindString(nil),
		bndValueFocused: binding.BindBool(nil),
	}

	grid.Add(NewCustomLocaleLabel("cont.rawDataReceived", &style{Bold: true}))

	img := widgets.NewKnobControlImage(resources.ResKnobPng, bindings.bndValueFocused)

	return container.NewVBox(grid, img), bindings
}

func (rd *RawData) refreshBindings(data []byte, layout *devicelayout.DeviceLayoutConfig) {
	if len(data) == 0 || len(data[1:]) == 0 {
		return
	}
	if layout == nil {
		slog.Error("Device layout is not loaded")
		return
	}
	for inx, comp := range layout.Components {
		hexStr, decStr, byteStr, value := getControlValues(comp, data[1:])
		slog.Debug("Data read from USB", "data", fmt.Sprintf("%s %s %s %s", comp.Name, decStr, hexStr, byteStr))
		rd.bndData[inx].bndValueBin.Set(byteStr)
		rd.bndData[inx].bndValueHex.Set(hexStr)
		rd.bndData[inx].bndValueDec.Set(decStr)
		rd.bndData[inx].bndValueFocused.Set(value&comp.Value != 0)
	}
}

func (rd *RawData) setData(dev *monitor.ConnectedDevice) {

	rd.onceGrid.Do(func() {
		rd.buildBindings(dev.DeviceLayoutConfig)
	})

	rd.stopChan = make(chan bool)
	go func() {
		for {
			select {
			case <-rd.stopChan:
				slog.Debug("Stopping RawData")
				return
			default:
				rd.refreshBindings(readUSB(dev.Device), dev.DeviceLayoutConfig)
			}
		}
	}()
}

func (rd *RawData) GetContent(dev *monitor.ConnectedDevice) *container.Scroll {
	rd.setData(dev)
	return rd.body
}

func (rd *RawData) GetButton() *widget.Button {
	return rd.button
}

func (rd *RawData) Destroy() {
	slog.Debug("Destroying RawData")
	select {
	case rd.stopChan <- true:
	default:
	}
}

func readUSB(dev *gousb.Device) []byte {
	cfg, err := dev.Config(1)
	if err != nil {
		slog.Error("Could not get config:", "error", err)
		return nil
	}
	defer cfg.Close()

	intf, err := cfg.Interface(3, 0)
	if err != nil {
		slog.Error("Could not get interface:", "error", err)
		return nil
	}

	defer intf.Close()

	// Setup the endpoint
	ep, err := intf.InEndpoint(4) // 1 is the endpoint number
	if err != nil {
		slog.Error("Could not get endpoint:", "error", err)
		return nil
	}

	// Read data from the endpoint
	data := make([]byte, 4)
	_, err = ep.Read(data)
	if err != nil {
		slog.Error("Could not read data:", "error", err)
		return nil
	}
	return data
}

func getControlValues(layoutComp devicelayout.Component, data []byte) (string, string, string, int) {
	byteVal := make([]byte, len(data[layoutComp.Bytes[0]:layoutComp.Bytes[1]+1]))
	copy(byteVal, data[layoutComp.Bytes[0]:layoutComp.Bytes[1]+1])

	var value int
	if layoutComp.Bytes[1]-layoutComp.Bytes[0]+1 > 1 { //more than one byte
		if layoutComp.Endianess == "big" {
			for i := 0; i < len(byteVal); i++ {
				value |= int(byteVal[i]) << (8 * (len(byteVal) - 1 - i))
			}
		} else if layoutComp.Endianess == "little" {
			for i := 0; i < len(byteVal); i++ {
				value |= int(byteVal[i]) << (8 * i)
			}
			byteTemp := byteVal[0]
			byteVal[0] = byteVal[1]
			byteVal[1] = byteTemp
		} else {
			slog.Warn("Endianess not specified")
		}
	} else { //one byte
		if layoutComp.ByteType == "signed" {
			value = int(int8(byteVal[0]))
		} else if layoutComp.ByteType == "unsigned" {
			value = int(uint8(byteVal[0]))
		} else {
			slog.Warn("Byte type not specified")
		}
	}

	if layoutComp.Value != 0 {
		byteStr := utility.FormatAsBinary(utility.AbsInt(value)&layoutComp.Value, layoutComp.Bytes[1]-layoutComp.Bytes[0])
		hexStr := fmt.Sprintf("0x%02X", utility.AbsInt(value)&layoutComp.Value)
		decStr := fmt.Sprintf("%d", value&layoutComp.Value)
		return hexStr, decStr, byteStr, value
	} else {
		byteStr := utility.FormatAsBinary(int(uint8(value)), layoutComp.Bytes[1]+1-layoutComp.Bytes[0])
		hexStr := fmt.Sprintf("0x%02X", int(uint8(value)))
		decStr := fmt.Sprintf("%d", value)
		return hexStr, decStr, byteStr, value
	}
}
