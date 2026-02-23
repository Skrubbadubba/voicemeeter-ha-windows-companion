package main

import "github.com/onyx-and-iris/voicemeeter/v2"

type stripCache struct {
	label string
	mute  bool
	solo  bool
	mono  bool
	gain  float64
	a1    bool
	a2    bool
	a3    bool
	a4    bool
	a5    bool
	b1    bool
	b2    bool
	b3    bool
}

type busCache struct {
	label string
	mute  bool
	mono  bool
	gain  float64
	eq    bool
}

type buttonCache struct {
	state     bool
	stateOnly bool
	trigger   bool
}

type vmCache struct {
	strips  []stripCache
	buses   []busCache
	buttons []buttonCache
}

func snapshot(vm *voicemeeter.Remote) vmCache {
	strips := make([]stripCache, len(vm.Strip))
	for i, s := range vm.Strip {
		strips[i] = stripCache{
			label: s.Label(),
			mute:  s.Mute(),
			solo:  s.Solo(),
			mono:  s.Mono(),
			gain:  s.Gain(),
			a1:    s.A1(),
			a2:    s.A2(),
			a3:    s.A3(),
			a4:    s.A4(),
			a5:    s.A5(),
			b1:    s.B1(),
			b2:    s.B2(),
			b3:    s.B3(),
		}
	}

	buses := make([]busCache, len(vm.Bus))
	for i, b := range vm.Bus {
		buses[i] = busCache{
			label: b.Label(),
			mute:  b.Mute(),
			mono:  b.Mono(),
			gain:  b.Gain(),
			eq:    b.Eq().On(),
		}
	}

	buttons := make([]buttonCache, len(vm.Button))
	for i, btn := range vm.Button {
		buttons[i] = buttonCache{
			state:     btn.State(),
			stateOnly: btn.StateOnly(),
			trigger:   btn.Trigger(),
		}
	}

	return vmCache{strips: strips, buses: buses, buttons: buttons}
}