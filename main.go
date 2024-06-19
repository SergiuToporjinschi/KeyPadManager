package main

import (
	"main/devicelayout"
	"main/gui"
	"main/logger"
	"main/txt"
)

//go:generate $GOPATH\bin\fyne bundle -prefix Lang -package txt -o txt/bundles.go assets/langs
//go:generate $GOPATH\bin\fyne bundle -prefix Res -package resources -o assets/bundled.go assets/files

func main() {
	logger.Init()
	logger.Log.Info("Starting the application")

	err := devicelayout.GetInstance().LoadConfig()
	if err != nil {
		return
	}

	txt.GetInstance().SetLanguage("en") //TODO get language from config
	g := gui.GetInstance()
	g.ShowDeviceWindow()
	g.App.Run()
}
