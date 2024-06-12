package gui

import (
	"main/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type menuChildElements map[string]MainContent
type menuOptions map[string]menuChildElements

func (g *GUI) getMenuConfig() map[string]menuChildElements {
	return map[string]menuChildElements{
		"Device": {
			"Info":      NewInfo(),
			"RawValues": NewRawValues(g),
		},
		"Profiles": {
			"Management":   NewProfile(g),
			"Running apps": NewRunningApps(g),
		},
	}
}

func (g *GUI) buidMenu(mainContent *fyne.Container) *widget.Accordion {
	menuItems := []*widget.AccordionItem{}
	for parentItem := range g.menuOptions {
		parentBtns := []fyne.CanvasObject{}
		logger.Log.Debugf("parentItem: %v", parentItem)
		logger.Log.Debugf("g.menuOptions[%v]: ", parentBtns)

		childContainer := container.NewVBox()
		for childItem := range g.menuOptions[parentItem] {
			childContainer.Add(widget.NewButton(childItem, func() { g.menuSelection(parentItem, childItem, mainContent) }))
		}
		menuItems = append(menuItems, widget.NewAccordionItem(parentItem, childContainer))

		//Add parent
		widget.NewAccordionItem(parentItem, container.NewVBox(parentBtns...))
	}
	return widget.NewAccordion(menuItems...)
}

func (g *GUI) menuSelection(parentMenu string, childMenu string, mainContent *fyne.Container) {
	for parentItem := range g.menuOptions {
		for childItem := range g.menuOptions[parentItem] {
			if parentItem == parentMenu && childItem == childMenu {
				continue
			}
			g.menuOptions[parentItem][childItem].Destroy()
		}
	}

	mainContent.RemoveAll()
	mainContent.Add(g.menuOptions[parentMenu][childMenu].Build(g.device))
	mainContent.Refresh()
}
