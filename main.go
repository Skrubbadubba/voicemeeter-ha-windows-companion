package main

import (
	"log"
	"os"

	"os/signal"

	// "slices"
	"syscall"

	"fyne.io/systray"
	"github.com/onyx-and-iris/voicemeeter/v2"
)

var globalVMR *voicemeeter.Remote
var globalServer *server

func main() {
	runCore()
}

func runCore() {
	vmr, k, err := detectKind()
	if err != nil {
		log.Fatal(err)
	}
	globalVMR = vmr

	log.Printf("Connected to %s (%d strips, %d buses, %d buttons)", k.name, k.strips, k.buses, k.buttons)

	server := newServer(vmr, k)
	if server == nil {
		log.Printf("server pointer is nil")
	}
	globalServer = server

	go server.start()
	go startPolling(vmr, server)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("Control Event registered")
		systray.Quit()
	}()

	runTray()
}

func onExit() {
	log.Printf("shutting down ws server")
	globalServer.shutdown()
	log.Println("Shutting down...")

	// Dont call logout on voicemeeter here, it panics. Maybe a bug in the library
	os.Exit(0)
}
