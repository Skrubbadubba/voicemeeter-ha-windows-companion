package main

import (
	_ "embed"
	"fmt"
	"log"

	"fyne.io/systray"
)

//go:embed ha-vm-square-source.ico
var iconData []byte

func runTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTooltip("Voicemeeter Companion")

	mStatus := systray.AddMenuItem("Voicemeeter Companion", "")
	mStatus.Disable()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Stop the companion app")

	mVersion := systray.AddMenuItem(fmt.Sprintf("Protocol Ver.: %s", PROTOCOL_VER), "")
	mVersion.Disable()

	go func() {
		for range mQuit.ClickedCh {
			log.Println("Quit requested from tray")
			systray.Quit()
		}
	}()
}
