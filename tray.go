package main

import (
	_ "embed"
	"fmt"
	"log"

	"fyne.io/systray"
)

//go:embed ha-vm-square-source.ico
var iconData []byte

func runTray(reconnect func()) {
	systray.Run(onReady(reconnect), onExit)
}

func onReady(reconnect func()) func() {
	return func() {
		systray.SetIcon(iconData)
		systray.SetTooltip("Voicemeeter Companion")

		mStatus := systray.AddMenuItem("Voicemeeter Companion", "")
		mStatus.Disable()

		systray.AddSeparator()

		mQuit := systray.AddMenuItem("Quit", "Stop the companion app")

		mReconnect := systray.AddMenuItem("Reonnect", "Try reconnecting to voicemeeter")

		mVersion := systray.AddMenuItem(fmt.Sprintf("Protocol Ver.: %s", PROTOCOL_VER), "")
		mVersion.Disable()

		go func() {
			for {
				select {
				case <-mReconnect.ClickedCh:
					reconnect()
				case <-mQuit.ClickedCh:
					log.Println("Quit requested from tray")
					systray.Quit()
				}

			}
		}()

	}
}
