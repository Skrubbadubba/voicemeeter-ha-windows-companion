package main

import (
	"log"
	"os"
	"time"

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
	vmr, k, err := tryConnectVM(0 * time.Minute)
	if err == nil && vmr != nil {
		log.Println("Successfully connected to Voicemeeter!")
	}

	globalVMR = vmr

	log.Printf("Connected to %s (%d strips, %d buses, %d buttons)", k.name, k.strips, k.buses, k.buttons)

	server := newServer(vmr, k)
	if server == nil {
		log.Printf("server pointer is nil")
	}
	globalServer = server
	go server.start()

	isPolling := false
	if vmr != nil {
		go startPolling(vmr, server)
		isPolling = true
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("Control Event registered")
		systray.Quit()
	}()

	runTray(func() {
		if isPolling {
			log.Println("already polling")
			return
		}
		vmr, _, err := tryConnectVM(2 * time.Minute)
		if err != nil {
			log.Printf("reconnection failed")
			return
		}
		go startPolling(vmr, server)
		isPolling = true

	})
}

func onExit() {
	log.Printf("shutting down ws server")
	globalServer.shutdown()
	log.Println("Shutting down...")

	// Dont call logout on voicemeeter here, it panics. Maybe a bug in the library
	os.Exit(0)
}
