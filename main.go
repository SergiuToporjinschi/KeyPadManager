package main

import (
	"main/gui"
	"main/logger"
)

func main() {
	logger.Init()
	logger.Log.Info("Starting the application")
	g := gui.NewGUI()
	g.OpenMain()
}
