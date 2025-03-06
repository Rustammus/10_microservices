package chat

import (
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 << 10
)

// client is a middleman between the websocket connection and the hub.
type client struct {
	hub *hub

	name string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan chatMsg
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println("client.readPump:  setReadLimit err: ", err)
	}

	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// run receive cycle
	for {
		msgType, reader, err := c.conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("recv msg from: %v msgType: %v", c.name, msgType)

		// create pipe and send to hub
		pr, pw := io.Pipe()
		c.hub.broadcast <- chatMsg{
			messageType: msgType,
			reader:      pr,
		}

		g, err := io.Copy(pw, reader)
		if err != nil {
			log.Printf("client.readPump: io.Copy err: %v", err)
			break
		}
		log.Printf("client.readPump: copyed %d bytes", g)

		err = pw.Close()
		if err != nil {
			log.Printf("client.readPump: pw.Close err: %v", err)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		// on send
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// get io.Writer
			writer, err := c.conn.NextWriter(message.messageType)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				c.hub.unregister <- c
				break
			}

			_, err = io.Copy(writer, message.reader)
			if err != nil {
				log.Printf("client.writePump: io.Copy err: %v", err)
			}

			err = writer.Close()
			if err != nil {
				log.Printf("client.writePump: writer.Close err: %v", err)
			}

		// pinging
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error on ping user: %s err: %v", c.name, err)
				return
			}
		}
	}
}

// register new client in hub
func newClient(hub *hub, name string, conn *websocket.Conn) {
	c := &client{hub: hub, name: name, conn: conn, send: make(chan chatMsg, 256)}
	c.hub.register <- c

	// start receiving and sending
	go c.writePump()
	go c.readPump()
}
