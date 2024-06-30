package gui

import (
	"log/slog"
	resources "main/assets"
	"main/monitor"
	"main/txt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type NavigationItem interface {
	GetContent(*monitor.ConnectedDevice) *container.Scroll
	Destroy()
}

var naviIndex = map[string][]string{
	"":                        {"navi.deviceBranchTitle", "navi.profileBranchTitle"},
	"navi.deviceBranchTitle":  {"info", "rawdata"},
	"navi.profileBranchTitle": {"something1", "something2"},
}

type Navigation struct {
	naviInitilizer func() NavigationItem
	naviIcon       fyne.Resource
	naviTitle      string
	naviInst       NavigationItem
}

var navigationContent = map[string]*Navigation{
	"info": {
		naviInitilizer: NewDeviceInfo,
		naviIcon:       resources.ResInfoPng,
		naviTitle:      "navi.deviceInfoTitle",
	},
	"rawdata": {
		naviInitilizer: NewRawData,
		naviIcon:       resources.ResMatrixPng,
		naviTitle:      "navi.rawDataTitle",
	},
}

type ContentManager struct {
	container.Split
	navi           *widget.Tree
	currentDevice  *monitor.ConnectedDevice
	currentNavItem NavigationItem
	mutex          sync.Mutex
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.currentDevice = device
}

func (c *ContentManager) OnMainWindowHide() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.currentNavItem != nil {
		c.currentNavItem.Destroy()
	}
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
	return container.NewHBox(widget.NewIcon(nil), widget.NewLabel(txt.GetLabel("common.notDefined")))
}

func (c *ContentManager) updateNode(uid string, branch bool, obj fyne.CanvasObject) {
	box := obj.(*fyne.Container)
	if branch {
		box.Objects[1].(*widget.Label).SetText(txt.GetLabel(uid))
		return
	}
	item, found := navigationContent[uid]
	if found {
		if item.naviInst == nil {
			item.naviInst = item.naviInitilizer()
		}
		box.Objects[1].(*widget.Label).SetText(txt.GetLabel(item.naviTitle))
		if item.naviIcon != nil {
			box.Objects[0] = widget.NewIcon(item.naviIcon)
		}
	} else {
		slog.Warn("Unknown UID", "UUID", uid)
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
		icn := canvas.NewImageFromResource(item.naviIcon)
		icn.FillMode = canvas.ImageFillOriginal
		mainContent.Add(
			container.NewBorder(
				container.NewHBox(
					icn,
					NewTitleLocaleText(item.naviTitle)),
				nil,
				nil,
				nil,
				item.naviInst.GetContent(n.currentDevice),
			),
		)
		n.currentNavItem = item.naviInst

	}
}
