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
	mon := monitoring.GetInstance()

	mon.AddDeviceEvent("MainFunc", func(event string, device monitoring.ConnectedDevice) {
		logger.Log.Infof("MAIN: Device connected: %s : %v", event, device)
	})

	mon.AddMonitorEvent("MainFunc", func(event string) {
		logger.Log.Infof("MAIN: Monitor connected: %v", event)
	})

	err = mon.Start()
	if err != nil {
		return
	}
	select {}
	// g := gui.NewGUI()
	// g.OpenMain()
}
