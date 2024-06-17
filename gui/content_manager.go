package gui

import (
	"image/color"
	"main/logger"
	"main/monitor"
	"main/txt"
	"main/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type NavigationItem interface {
	GetTitle() string
	GetNavTitle() string
	GetContent(*monitor.ConnectedDevice) *fyne.Container
	Destroy()
}

var naviIndex = map[string][]string{
	"":                        {"navi.deviceBranchTitle", "navi.profileBranchTitle"},
	"navi.deviceBranchTitle":  {"info", "rawdata"},
	"navi.profileBranchTitle": {"something1", "something2"},
}

type Navigation struct {
	naviInitilizer func() NavigationItem
	naviInst       NavigationItem
}

var navigationContent = map[string]*Navigation{
	"info": {
		naviInitilizer: NewDeviceInfo,
	},
	"rawdata": {
		naviInitilizer: NewRawData,
	},
}

type ContentManager struct {
	container.Split
	navi           *widget.Tree
	currentDevice  *monitor.ConnectedDevice
	currentNavItem NavigationItem
}

func NewContentManager() *ContentManager {
	s := &ContentManager{
		Split: container.Split{
			Offset:     0.2, // Sensible default, can be overridden with SetOffset
			Horizontal: true,
			Trailing:   container.NewStack(),
		},
	}
	s.buildNavigation()
	s.Split.Leading = s.navi
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

func (c *ContentManager) SetDevice(device *monitor.ConnectedDevice) {
	c.currentDevice = device
}

func (c *ContentManager) OnMainWindowHide() {
	c.currentNavItem.Destroy()
}

func (c *ContentManager) buildNavigation() {
	c.navi = widget.NewTreeWithStrings(naviIndex)
	c.navi.ChildUIDs = c.childUIDs
	c.navi.IsBranch = c.isBranch
	c.navi.CreateNode = c.createNode
	c.navi.UpdateNode = c.updateNode
	c.navi.OnSelected = c.onSelected

}

func (c *ContentManager) childUIDs(uid string) []string {
	return naviIndex[uid]
}

func (c *ContentManager) isBranch(uid string) bool {
	children, ok := naviIndex[uid]

	return ok && len(children) > 0
}

func (c *ContentManager) createNode(branch bool) fyne.CanvasObject {
	return widget.NewLabel(txt.GetLabel("common.notDefined"))
}

func (c *ContentManager) updateNode(uid string, branch bool, obj fyne.CanvasObject) {
	if branch {
		obj.(*widget.Label).SetText(txt.GetLabel(uid))
		return
	}
	item, found := navigationContent[uid]
	if found {
		if item.naviInst == nil {
			item.naviInst = item.naviInitilizer()
		}
		obj.(*widget.Label).SetText(item.naviInst.GetNavTitle())
	} else {
		logger.Log.Warnf("Unknown UID: %v", uid)
	}
}

func (n *ContentManager) onSelected(uid string) {
	item, found := navigationContent[uid]
	if found {
		mainContent := n.Trailing.(*fyne.Container)
		if n.currentNavItem != nil {
			n.currentNavItem.Destroy()
		}
		mainContent.RemoveAll()
		mainContent.Add(
			container.NewBorder(
				newTitleText(item.naviInst.GetTitle()),
				nil,
				nil,
				nil,
				item.naviInst.GetContent(n.currentDevice),
			),
		)
		n.currentNavItem = item.naviInst

	}
}

func newTitleText(text string) *canvas.Text {
	r := utility.NewSizeableText(text, 20)
	r.TextStyle.Bold = true
	return utility.NewSizeableColorText(text, 20, color.NRGBA{R: 0xFE, G: 0x58, B: 0x62, A: 0xFF})
}
