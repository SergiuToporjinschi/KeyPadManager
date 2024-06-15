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

// func onReady(label *widget.Label) func() {
// 	return func() {
// 		data, err := os.ReadFile("icon.ico")
// 		if err != nil {
// 			log.Fatalf("Failed to read file: %v", err)
// 		}
// 		// Set the icon of the systray
// 		systray.SetIcon(data) // iconData should be a byte slice containing the icon data

// 		// Set the tooltip of the systray
// 		systray.SetTooltip("System Tray Example")

// 		// Add menu items to the systray
// 		mHello := systray.AddMenuItem("Hello", "Say Hello")
// 		mQuit := systray.AddMenuItem("Quit", "Quit the application")

// 		// Handle menu item clicks
// 		go func() {
// 			for {
// 				select {
// 				case <-mHello.ClickedCh:
// 					label.SetText("Hello menu item clicked")
// 				case <-mQuit.ClickedCh:
// 					systray.Quit()
// 					return
// 				}
// 			}
// 		}()
// 	}
// }
