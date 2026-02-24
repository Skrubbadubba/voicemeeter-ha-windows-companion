package main

import (
	_ "embed"
	"log"

	"fyne.io/systray"
)

//go:embed iconwin.ico
var iconData []byte

var Version = "dev" // overridden at build time

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

	mVersion := systray.AddMenuItem("Protocol Ver.: "+Version, "")
	mVersion.Disable()

	go func() {
		for range mQuit.ClickedCh {
			log.Println("Quit requested from tray")
			systray.Quit()
		}
	}()
}
