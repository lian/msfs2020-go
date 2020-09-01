package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type ReceiveMessage struct {
	Message    []byte
	Connection *Connection
}

type Connection struct {
	socket    *Websocket
	conn      *websocket.Conn
	Send      chan []byte
	SendQueue chan []byte
}

func (c *Connection) Run() {
	go c.readPump()
	go c.writer()
	c.writePump()
}

func (c *Connection) SendPacket(data map[string]interface{}) {
	buf, _ := json.Marshal(data)
	c.Send <- buf
}

func (c *Connection) SendError(target string, msg string) {
	pkt := map[string]string{"target": target, "type": "error", "message": msg}
	buf, _ := json.Marshal(pkt)
	c.Send <- buf
}

func (c *Connection) readPump() {
	defer func() {
		c.socket.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				//log.Printf("error: %v", err)
			}
			//log.Printf("error: %v", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.socket.ReceiveMessages <- ReceiveMessage{
			Message:    message,
			Connection: c,
		}
	}
}

func (c *Connection) writer() {
	var buf bytes.Buffer
	flushticker := time.NewTicker(time.Millisecond * 16)

	defer func() {
		flushticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				close(c.SendQueue)
				return
			}
			buf.Write(message)
			buf.WriteString("\n")

		case <-flushticker.C:
			message, err := buf.ReadBytes('\n')
			if err == nil {
				c.SendQueue <- message
			}
		}
	}
}

func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendQueue:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				fmt.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
