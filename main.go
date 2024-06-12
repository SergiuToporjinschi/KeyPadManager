package main

import (
	devicelayout "main/deviceLayout"
	"main/gui"
	"main/logger"
)

func main() {
	logger.Init()
	logger.Log.Info("Starting the application")

	err := devicelayout.GetInstance().LoadConfig()
	if err != nil {
		return
	}

	g := gui.NewGUI()
	g.OpenMain()
}
