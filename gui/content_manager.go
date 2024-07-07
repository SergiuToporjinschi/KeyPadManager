package gui

import (
	"log/slog"
	"main/gui/screens"
	"main/monitor"
	"main/txt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ContentManager struct {
	container.Split
	navi           *widget.Tree
	currentDevice  *monitor.ConnectedDevice
	currentNavItem screens.NavigationItem
	mutex          sync.Mutex
	lastSelected   string
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

	if c.lastSelected == "" {
		c.lastSelected = screens.GetFirstNaviIndexID()
	}
	c.navi.FocusGained()
	c.navi.Select(c.lastSelected)
	c.onSelected(c.lastSelected)
}

func (c *ContentManager) OnMainWindowHide() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mainContent := c.Trailing.(*fyne.Container)
	mainContent.RemoveAll()
	if c.currentNavItem != nil {
		c.currentNavItem.Destroy()
		c.currentNavItem = nil
	}
}

func (c *ContentManager) buildNavigation() {
	c.navi = widget.NewTreeWithStrings(screens.NaviIndex)
	c.navi.ChildUIDs = c.childUIDs
	c.navi.IsBranch = c.isBranch
	c.navi.CreateNode = c.createNode
	c.navi.UpdateNode = c.updateNode
	c.navi.OnSelected = c.onSelected
	c.navi.OnUnselected = c.onUnselected
}

func (c *ContentManager) childUIDs(uid string) []string {
	return screens.NaviIndex[uid]
}

func (c *ContentManager) isBranch(uid string) bool {
	children, ok := screens.NaviIndex[uid]

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
	item, found := screens.NavigationContent[uid]
	if found {
		// if item.Inst == nil {
		// 	item.Inst = item.Initilizer()
		// }
		box.Objects[1].(*widget.Label).SetText(txt.GetLabel(item.Title))
		if item.Icon != nil {
			box.Objects[0] = widget.NewIcon(item.Icon)
		}
	} else {
		slog.Warn("Unknown UID", "UUID", uid)
	}
}
func (n *ContentManager) onUnselected(uid string) {
	mainContent := n.Trailing.(*fyne.Container)
	mainContent.RemoveAll()
	if n.currentNavItem != nil {
		n.currentNavItem.Destroy()
		n.currentNavItem = nil
	}
}

func (n *ContentManager) onSelected(uid string) {
	n.lastSelected = uid
	item, found := screens.NavigationContent[uid]
	if found {
		naviContentInst := item.Initilizer(n.currentDevice)
		mainContent := n.Trailing.(*fyne.Container)
		if n.currentNavItem != nil {
			n.currentNavItem.Destroy()
		}
		mainContent.RemoveAll()
		icn := canvas.NewImageFromResource(item.Icon)
		icn.FillMode = canvas.ImageFillOriginal
		mainContent.Add(
			container.NewBorder(
				container.NewHBox(
					icn,
					NewTitleLocaleText(item.Title)),
				nil,
				nil,
				nil,
				naviContentInst.GetContent(),
			),
		)
		n.currentNavItem = naviContentInst
	}
}
