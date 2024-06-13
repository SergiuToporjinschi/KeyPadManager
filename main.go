package main

import (
	"log"
	"main/devicelayout"
	"main/gui"
	"main/logger"
	"os"

	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

func main() {
	logger.Init()
	logger.Log.Info("Starting the application")

	err := devicelayout.GetInstance().LoadConfig()
	if err != nil {
		return
	}
	// mon := monitoring.GetInstance()

	// mon.AddMonitorEvent("MainFunc", func(event string) {
	// 	logger.Log.Infof("MAIN: Monitor connected: %v", event)
	// })

	g := gui.GetInstance()
	g.ShowDeviceWindow()
	g.App.Run()
	// err = mon.Start()
	// if err != nil {
	// 	return
	// }

	// g := gui.NewGUI()
	// g.OpenMain()
	// fyne.StaticResourceFromPath("devKeypad.png")
	// myApp := app.New()
	// myApp.Storage()
	// myApp.Settings().SetTheme(NewTheme(myApp.Preferences()))

	// // selectdevice.GetInstance().NewWindow("Selected device")

	// deviceWindow := myApp.NewWindow("Selected device")
	// menuWindow := myApp.NewWindow("menu ")

	// // crd.Resize(fyne.NewSize(2800, 2600))
	// img := canvas.NewImageFromFile("devKeypad.png")
	// img.FillMode = canvas.ImageFillOriginal

	// card := widget.NewCard(
	// 	"S", // Title
	// 	"x", // Description
	// 	img, // Content
	// )

	// paddedCard := container.NewVBox(
	// 	// layout.NewSpacer(), // Top padding
	// 	container.NewHBox(
	// 		// layout.NewSpacer(), // Left padding
	// 		card,
	// 		// layout.NewSpacer(), // Right padding
	// 	),
	// 	// layout.NewSpacer(), // Bottom padding
	// )
	// deviceWindow.SetOnClosed(func() {
	// })
	// deviceWindow.SetPadded(true)
	// deviceWindow.SetContent(
	// 	paddedCard,
	// )
	// mon.AddDeviceEvent("MainFunc", func(event string, device monitoring.ConnectedDevice) {
	// 	if event == "disconnected" {
	// 		deviceWindow.Show()
	// 		menuWindow.Hide()
	// 	}
	// })
	// systray.Run(onReady(widget.NewLabel("Key Pad Manager")), func() {

	// })
	// deviceWindow.Show()
	// menuWindow.Hide()
	// myApp.Run()

}
func onReady(label *widget.Label) func() {
	return func() {
		data, err := os.ReadFile("icon.ico")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		// Set the icon of the systray
		systray.SetIcon(data) // iconData should be a byte slice containing the icon data

		// Set the tooltip of the systray
		systray.SetTooltip("System Tray Example")

		// Add menu items to the systray
		mHello := systray.AddMenuItem("Hello", "Say Hello")
		mQuit := systray.AddMenuItem("Quit", "Quit the application")

		// Handle menu item clicks
		go func() {
			for {
				select {
				case <-mHello.ClickedCh:
					label.SetText("Hello menu item clicked")
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		}()
	}
}
