package main

import (
	_ "embed"
	"log"

	"fyne.io/systray"
)

//go:embed iconwin.ico
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

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				log.Println("Quit requested from tray")
				systray.Quit()
			}
		}
	}()
}
