package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/systray"
	"github.com/onyx-and-iris/voicemeeter/v2"
)

var (
	kernel32      = syscall.NewLazyDLL("kernel32.dll")
	attachConsole = kernel32.NewProc("AttachConsole")
	allocConsole  = kernel32.NewProc("AllocConsole")
)

const ATTACH_PARENT_PROCESS = ^uintptr(0) // -1 as uintptr

func tryAttachConsole() bool {
	ret, _, _ := attachConsole.Call(ATTACH_PARENT_PROCESS)
	return ret != 0
}

func redirectStdioToConsole() {
	stdout, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err == nil {
		os.Stdout = os.NewFile(uintptr(stdout), "CONOUT$")
	}
	stderr, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err == nil {
		os.Stderr = os.NewFile(uintptr(stderr), "CONOUT$")
	}
	log.SetOutput(io.MultiWriter(os.Stdout, os.Stderr))
}

var globalVMR *voicemeeter.Remote
var globalServer *server

func main() {
	foreground := tryAttachConsole()
	if foreground {
		redirectStdioToConsole()
		log.Println("Voicemeeter companion starting in foreground mode")
	}

	vm, k, err := detectKind()
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Logout()

	log.Printf("Connected to %s (%d strips, %d buses, %d buttons)", k.name, k.strips, k.buses, k.buttons)

	server := newServer(vm, k)
	globalServer = server

	go server.start()
	go startPolling(vm, server)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		systray.Quit()
	}()

	runTray()
}

func onExit() {
	log.Println("Shutting down")
	globalVMR.Logout()
	globalServer.shutdown()
	os.Exit(0)
}
