package screens

import (
	"log/slog"
	resources "main/assets"
	"main/monitor"

	"fyne.io/fyne/v2"
)

type NavigationItem interface {
	GetContent() *fyne.Container
	Destroy()
}
type Navigation struct {
	Initilizer func(*monitor.ConnectedDevice) NavigationItem
	Icon       fyne.Resource
	Title      string
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

func GetFirstNaviIndexID() string {
	root, found := NaviIndex[""]
	if found {
		return root[0]
	} else {
		slog.Error("No root navigation item found")
	}
	return ""
}
