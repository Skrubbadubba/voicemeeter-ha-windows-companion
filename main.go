package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jwfergus/winsystray"
)

func main() {
	vm, k, err := detectKind()
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Logout()

	log.Printf("Connected to %s (%d strips, %d buses, %d buttons)", k.name, k.strips, k.buses, k.buttons)

	server := newServer(vm, k)

	go server.start()
	go startPolling(vm, server)

	ti, err := winsystray.NewTrayIcon()
	if err != nil {
		panic(err)
	}
	defer ti.Dispose()

	/*
		These can be called as frequently as necessary. Changing the
		icon quickly gives the illusion of animation.
	*/
	// ti.SetIconFromFile("icon.ico")
	ti.SetTooltip("おはよう世界！")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
