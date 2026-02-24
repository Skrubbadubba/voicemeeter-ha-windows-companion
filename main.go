package main

import (
	"log"
	"os"

	"os/signal"

	// "slices"
	"syscall"

	// "fyne.io/systray"
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
	globalServer = server

	go server.start()
	go startPolling(vmr, server)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("Control Event registered")
		onExit()
	}()

	runTray()
}

func onExit() {
	log.Println("Shutting down...")
	if globalVMR != nil {
		globalVMR.Logout()
	}
	if globalServer != nil {
		globalServer.shutdown()
	}
	os.Exit(0)
}
