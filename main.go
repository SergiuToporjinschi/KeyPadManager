package main

import (
	"main/devicelayout"
	"main/logger"
	"main/monitoring"
)

func main() {
	logger.Init()
	logger.Log.Info("Starting the application")

	err := devicelayout.GetInstance().LoadConfig()
	if err != nil {
		return
	}
	err = monitoring.GetInstance().Start()
	if err != nil {
		return
	}
	select {}
	// g := gui.NewGUI()
	// g.OpenMain()
}
