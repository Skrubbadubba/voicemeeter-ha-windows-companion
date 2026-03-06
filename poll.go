package main

import (
	"log"
	"sync"

	"github.com/onyx-and-iris/voicemeeter/v2"
)

type Poller struct {
	mu           sync.Mutex
	isPolling    bool
	quit         chan struct{}
	disconnected chan struct{}
	vmr          *voicemeeter.Remote
}

func NewPoller(vmr *voicemeeter.Remote) *Poller {
	return &Poller{
		vmr: vmr,
	}
}

func (p *Poller) start(srv *server) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isPolling {
		log.Println("poller already running")
		return
	}

	p.isPolling = true
	p.quit = make(chan struct{})
	events := make(chan string, 8)
	p.vmr.Register(events)
	p.vmr.EventAdd("pdirty", "mdirty")

	cache := snapshot(p.vmr)
	log.Println("Listening for changes...")

	go func() {
		for {
			select {
			case <-p.quit:
				log.Println("Poller stopped.")
				return

			case event, ok := <-events:
				if !ok {
					log.Printf("Poller disconnected")
					return
				}

				switch event {
				case "pdirty":
					fresh := snapshot(p.vmr)
					broadcastDiff(srv, cache, fresh)
					cache = fresh
				case "mdirty":
					fresh := snapshot(p.vmr)
					diffButtons(cache.buttons, fresh.buttons)
					cache = fresh
				}
			}
		}
	}()
}

func (p *Poller) stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isPolling {
		return
	}

	close(p.quit)

	p.vmr.EventRemove("pdirty", "mdirty")
	p.isPolling = false
}

func broadcastDiff(srv *server, old, new vmCache) {
	for i := range new.strips {
		o, n := old.strips[i], new.strips[i]
		if o.mute != n.mute {
			log.Printf("Strip %d: mute %v → %v", i, o.mute, n.mute)
			srv.broadcast(updateMsg{Type: "update", Target: "strip", Index: i, Param: "mute", Value: n.mute})
		}
		if o.gain != n.gain {
			log.Printf("Strip %d: gain %.1f → %.1f dB", i, o.gain, n.gain)
			srv.broadcast(updateMsg{Type: "update", Target: "strip", Index: i, Param: "gain", Value: n.gain})
		}
		if o.label != n.label {
			log.Printf("Strip %d: label %q → %q", i, o.label, n.label)
		}
		for _, r := range []struct {
			param string
			o, n  bool
		}{
			{"a1", o.a1, n.a1},
			{"a2", o.a2, n.a2},
			{"a3", o.a3, n.a3},
			{"a4", o.a4, n.a4},
			{"a5", o.a5, n.a5},
			{"b1", o.b1, n.b1},
			{"b2", o.b2, n.b2},
			{"b3", o.b3, n.b3},
		} {
			if r.o != r.n {
				log.Printf("Strip %d: %s %v → %v", i, r.param, r.o, r.n)
				srv.broadcast(updateMsg{Type: "update", Target: "strip", Index: i, Param: r.param, Value: r.n})
			}
		}
	}

	for i := range new.buses {
		o, n := old.buses[i], new.buses[i]
		if o.mute != n.mute {
			log.Printf("Bus %d: mute %v → %v", i, o.mute, n.mute)
			srv.broadcast(updateMsg{Type: "update", Target: "bus", Index: i, Param: "mute", Value: n.mute})
		}
		if o.gain != n.gain {
			log.Printf("Bus %d: gain %.1f → %.1f dB", i, o.gain, n.gain)
			srv.broadcast(updateMsg{Type: "update", Target: "bus", Index: i, Param: "gain", Value: n.gain})
		}
	}
}

func diffButtons(old, new []buttonCache) {
	for i := range new {
		o, n := old[i], new[i]
		if o.state != n.state {
			log.Printf("Button %d: state %v → %v", i, o.state, n.state)
		}
	}
}
