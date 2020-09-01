package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Websocket struct {
	connections     map[*Connection]bool
	broadcast       chan []byte
	register        chan *Connection
	unregister      chan *Connection
	receive         chan []byte
	ReceiveMessages chan ReceiveMessage
	NewConnection   chan ReceiveMessage
}

func New() *Websocket {
	ws := &Websocket{
		broadcast:       make(chan []byte, 256),
		register:        make(chan *Connection),
		unregister:      make(chan *Connection),
		connections:     make(map[*Connection]bool),
		ReceiveMessages: make(chan ReceiveMessage, 256),
		NewConnection:   make(chan ReceiveMessage, 256),
	}
	go ws.Run()

	return ws
}

func (s *Websocket) ConnectionCount() int {
	return len(s.connections)
}

func (s *Websocket) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &Connection{
		socket:    s,
		conn:      conn,
		Send:      make(chan []byte, 256),
		SendQueue: make(chan []byte),
	}
	s.register <- c

	c.Run()
}

func (s *Websocket) Broadcast(pkt map[string]interface{}) {
	buf, _ := json.Marshal(pkt)
	s.broadcast <- buf
}

func (h *Websocket) Run() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("new browser connection")
			h.connections[c] = true
			h.NewConnection <- ReceiveMessage{Connection: c}
		case c := <-h.unregister:
			fmt.Println("remove browser connection")
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.Send)
			}
		case packet := <-h.broadcast:
			for c := range h.connections {
				c.Send <- packet
				/*
					select {
					case c.Send <- packet:
					default:
						close(c.Send)
						delete(h.connections, c)
					}
				*/
			}
		}
	}
}
