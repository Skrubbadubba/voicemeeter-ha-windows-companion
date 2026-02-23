package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/onyx-and-iris/voicemeeter/v2"
)

const addr = ":27001"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type server struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
	vmr     *voicemeeter.Remote
	k       kind
}

func newServer(vm *voicemeeter.Remote, k kind) *server {
	return &server{
		clients: make(map[*websocket.Conn]struct{}),
		vmr:     vm,
		k:       k,
	}
}

func (s *server) start() {
	http.HandleFunc("/ws", s.handleConn)
	log.Printf("WebSocket server listening on ws://localhost%s/ws", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("WebSocket server error: %v", err)
	}
}

func (s *server) handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.clients[conn] = struct{}{}
	s.mu.Unlock()

	log.Printf("Client connected: %s", conn.RemoteAddr())

	// Protocol requires sending full state immediately on connect.
	if err := s.sendState(conn); err != nil {
		log.Printf("sendState error: %v", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		s.handleIncoming(conn, msg)
	}

	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()

	log.Printf("Client disconnected: %s", conn.RemoteAddr())
}

func (s *server) sendState(conn *websocket.Conn) error {
	msg := s.buildStateMsg()
	return writeJSON(conn, msg)
}

func (s *server) broadcast(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast marshal error: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("broadcast write error: %v", err)
			conn.Close()
			delete(s.clients, conn)
		}
	}
	log.Printf("Sent json update %v", msg)
}

func (s *server) handleIncoming(conn *websocket.Conn, raw []byte) {
	// Decode just the type field first.
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &base); err != nil {
		log.Printf("handleIncoming: bad JSON: %v", err)
		return
	}

	switch base.Type {
	case "set":
		var msg setMsg
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("handleIncoming: bad set message: %v", err)
			return
		}
		log.Printf("Recieved JSON: %v", msg)
		s.applySetMessage(msg)
	default:
		// Protocol says unknown types are silently ignored.
	}
}

func applyStripRouting(vm *voicemeeter.Remote, index int, param string, val bool) bool {
	switch param {
	case "a1":
		vm.Strip[index].SetA1(val)
	case "a2":
		vm.Strip[index].SetA2(val)
	case "a3":
		vm.Strip[index].SetA3(val)
	case "a4":
		vm.Strip[index].SetA4(val)
	case "a5":
		vm.Strip[index].SetA5(val)
	case "b1":
		vm.Strip[index].SetB1(val)
	case "b2":
		vm.Strip[index].SetB2(val)
	case "b3":
		vm.Strip[index].SetB3(val)
	default:
		return false
	}
	return true
}

func (s *server) applySetMessage(msg setMsg) {
	switch msg.Target {
	case "strip":
		if msg.Index < 0 || msg.Index >= len(s.vmr.Strip) {
			log.Printf("applySetMessage: strip index %d out of range", msg.Index)
			return
		}
		switch msg.Param {
		case "mute":
			val, ok := msg.Value.(bool)
			if !ok {
				log.Printf("applySetMessage: mute value must be bool")
				return
			}
			s.vmr.Strip[msg.Index].SetMute(val)
		case "gain":
			val, ok := msg.Value.(float64)
			if !ok {
				log.Printf("applySetMessage: gain value must be float")
				return
			}
			s.vmr.Strip[msg.Index].SetGain(val)
		default:
			if !(len(msg.Param) == 2 && (msg.Param[0] == 'a' || msg.Param[0] == 'b')) {
				log.Printf("applySetMessage: unknown strip param %q", msg.Param)
				return
			}
			val, ok := msg.Value.(bool)
			if !ok {
				log.Printf("applySetMessage: unknown strip param %q", msg.Param)
				return
			}

			if !applyStripRouting(s.vmr, msg.Index, msg.Param, val) {
				log.Printf("applySetMessage: unknown strip param %q", msg.Param)
				return
			}
			return
		}

	case "bus":
		if msg.Index < 0 || msg.Index >= len(s.vmr.Bus) {
			log.Printf("applySetMessage: bus index %d out of range", msg.Index)
			return
		}
		switch msg.Param {
		case "mute":
			val, ok := msg.Value.(bool)
			if !ok {
				log.Printf("applySetMessage: mute value must be bool")
				return
			}
			s.vmr.Bus[msg.Index].SetMute(val)
		case "gain":
			val, ok := msg.Value.(float64)
			if !ok {
				log.Printf("applySetMessage: gain value must be float")
				return
			}
			s.vmr.Bus[msg.Index].SetGain(val)
		default:
			log.Printf("applySetMessage: unknown bus param %q", msg.Param)
			return
		}

	default:
		log.Printf("applySetMessage: unknown target %q", msg.Target)
		return
	}

	// If update was succesfull, it will be picked up by polling and sent back to HA
	// No need for immediate confirmation
}

func (s *server) buildStateMsg() stateMsg {
	strips := make([]stripState, len(s.vmr.Strip))
	for i, strip := range s.vmr.Strip {
		strips[i] = stripState{
			Index:   i,
			Label:   strip.Label(),
			Mute:    strip.Mute(),
			Gain:    strip.Gain(),
			Virtual: i >= s.k.physicalStrips,
			A1:      strip.A1(),
			A2:      strip.A2(),
			A3:      strip.A3(),
			A4:      strip.A4(),
			A5:      strip.A5(),
			B1:      strip.B1(),
			B2:      strip.B2(),
			B3:      strip.B3(),
		}
	}

	buses := make([]busState, len(s.vmr.Bus))
	for i, bus := range s.vmr.Bus {
		buses[i] = busState{
			Index: i,
			Label: bus.Label(),
			Mute:  bus.Mute(),
			Gain:  bus.Gain(),
		}
	}

	return stateMsg{
		Type:   "state",
		Kind:   s.k.name,
		Strips: strips,
		Buses:  buses,
	}
}

func writeJSON(conn *websocket.Conn, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
