package main

import (
	"log/slog"
	"main/datahandlers"
	"main/devicelayout"
	"main/gui"
	"main/logger"
	"main/txt"
)

//go:generate $GOPATH\bin\fyne bundle -prefix Res -package resources -o assets/bundled.go assets/files
//go:generate $GOPATH\bin\fyne bundle -prefix Lang -package txt -o txt/bundles.go assets/langs

func main() {
	slog.SetDefault(logger.NewSLogger())
	slog.Info("Starting the application")

	go datahandlers.GetAppHandlerInstance().LoadAppList()
	go datahandlers.GetMacrosHandlerInstance().LoadMacroList()
	go datahandlers.GetKeysHandlerInstance().GenerateKeyList()
	go datahandlers.GetMappingHandlerInstance().LoadMapping()

	err := devicelayout.GetInstance().LoadConfig()
	if err != nil {
		return
	}

	txt.GetInstance().SetLanguage("en") //TODO get language from config
	g := gui.GetInstance()
	g.ShowDeviceWindow()
	g.App.Run()
}
