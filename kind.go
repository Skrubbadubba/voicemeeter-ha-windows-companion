package main

import (
	"fmt"
	"log"

	"github.com/onyx-and-iris/voicemeeter/v2"
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
		return nil, err
	}
	return vm, nil
}
