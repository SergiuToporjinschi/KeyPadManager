package screens

import (
	resources "main/assets"
	"main/monitor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type NavigationItem interface {
	GetContent(*monitor.ConnectedDevice) *container.Scroll
	Destroy()
}
type Navigation struct {
	Initilizer func() NavigationItem
	Icon       fyne.Resource
	Title      string
	Inst       NavigationItem
}

var NaviIndex = map[string][]string{
	"":                     {"device", "mapping", "navi.managementTitle"},
	"device":               {},
	"mapping":              {},
	"navi.managementTitle": {"apps", "macros", "scripts"},
}

var NavigationContent = map[string]*Navigation{
	"device": {
		Initilizer: NewDeviceScreen,
		Icon:       resources.ResDevicePng,
		Title:      "navi.deviceTitle",
	},
	"apps": {
		Initilizer: NewAppsScreen,
		Icon:       resources.ResApplicationPng,
		Title:      "navi.appsTitle",
	},
	"macros": {
		Initilizer: NewMacrosScreen,
		Icon:       resources.ResMacrosPng,
		Title:      "navi.macrosTitle",
	},
	"mapping": {
		Initilizer: NewMacrosScreen,
		Icon:       resources.ResMacrosPng,
		Title:      "navi.mappingTitle",
	},
	"scripts": {
		Initilizer: NewScriptsScreen,
		Icon:       resources.ResScriptPng,
		Title:      "navi.scriptTitle",
	},
}
