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
	vmr, k, err := tryConnectVM(2 * time.Minute)
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

	poller := NewPoller(vmr)
	if vmr != nil {
		poller.start(server)
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("Control Event registered")
		systray.Quit()
	}()

	runTray(
		// callback for reconnecting
		func() {
			if poller.isPolling {
				log.Println("Already polling, stopping")
				poller.stop()
			}
			if vmr != nil {
				log.Printf("Voicemeeter already connected, logging out")
				if err := vmr.Logout(); err != nil {
					log.Fatalf("Fatal: could not logout from voicemeeter")
				}
			}
			vmr, k, err := tryConnectVM(2 * time.Minute)
			if err != nil || vmr == nil {
				log.Println("Could not reconnect to Voicemeeter")
				return
			}
			log.Printf("Reconnected to voicemeeter, starting polling and sending state messages")
			poller.vmr = vmr
			server.vmr = vmr
			server.k = k
			for conn := range server.clients {
				server.sendState(conn)
			}
		})
}

func onExit() {
	log.Printf("shutting down ws server")
	globalServer.shutdown()
	log.Println("Shutting down...")

	// Dont call logout on voicemeeter here, it panics. Maybe a bug in the library
	os.Exit(0)
}
