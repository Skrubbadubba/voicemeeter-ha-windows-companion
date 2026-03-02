package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/onyx-and-iris/voicemeeter/v2"
	"golang.org/x/sys/windows/registry"
)

type kind struct {
	name           string
	strips         int
	buses          int
	buttons        int
	physicalStrips int
}

var kinds = map[string]kind{
	"basic":  {name: "basic", strips: 3, buses: 2, buttons: 80, physicalStrips: 2},
	"banana": {name: "banana", strips: 5, buses: 5, buttons: 80, physicalStrips: 3},
	"potato": {name: "potato", strips: 8, buses: 8, buttons: 80, physicalStrips: 5},
}

func tryConnectVM(timeout time.Duration) (*voicemeeter.Remote, kind, error) {
	var vmr *voicemeeter.Remote
	var k kind
	var err error

	deadline := time.Now().Add(0 * time.Minute)

	for {
		vmr, k, err = detectKind()
		if err == nil && vmr != nil {
			break
		}

		if time.Now().After(deadline) {
			log.Printf("Failed to connect to Voicemeeter after %s.", timeout.String())
			break
		}

		log.Printf("Voicemeeter not found, retrying in 10s... (Error: %v)", err)
		time.Sleep(10 * time.Second)
	}
	return vmr, k, err
}

func detectKind() (*voicemeeter.Remote, kind, error) {
	vm, err := connect("potato")
	if err != nil {
		return nil, kind{}, fmt.Errorf("initial connection failed: %w", err)
	}

	detected := vm.Type()
	log.Printf("Detected Voicemeeter kind: %s", detected)

	k, ok := kinds[detected]
	if !ok {
		return nil, kind{}, fmt.Errorf("unrecognised Voicemeeter kind: %q", detected)
	}

	if detected != "potato" {
		vm.Logout()
		vm, err = connect(detected)
		if err != nil {
			return nil, kind{}, fmt.Errorf("reconnect as %s failed: %w", detected, err)
		}
	}

	return vm, k, nil
}

func connect(kindName string) (*voicemeeter.Remote, error) {
	vm, err := voicemeeter.NewRemote(kindName, 0)
	if err != nil {
		return nil, err
	}
	if err := vm.Login(); err != nil {
		rawLogout()
		return nil, err
	}
	return vm, nil
}

func rawLogout() {
	// 1. Replicate the library's internal behavior: Find Voicemeeter in the Windows Registry
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\VB:Voicemeeter {17359A74-1236-5467}`, registry.QUERY_VALUE)
	if err != nil {
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\VB:Voicemeeter {17359A74-1236-5467}`, registry.QUERY_VALUE)
		if err != nil {
			return
		}
	}
	defer k.Close()

	uninst, _, err := k.GetStringValue("UninstallString")
	if err != nil {
		return
	}

	dir := filepath.Dir(uninst)
	dllName := "VoicemeeterRemote.dll"
	if runtime.GOARCH == "amd64" {
		dllName = "VoicemeeterRemote64.dll"
	}
	dllPath := filepath.Join(dir, dllName)

	dll := syscall.NewLazyDLL(dllPath)
	logoutProc := dll.NewProc("VBVMR_Logout")

	logoutProc.Call()
}
