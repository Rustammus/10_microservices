package chat

import (
	"io"
	"log"
)

// hub maintains the set of active clients and broadcasts messages to the clients
type hub struct {
	// Registered clients.
	clients map[*client]bool

	// Inbound chatMsg from the clients.
	broadcast chan chatMsg

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client
}

// create new hub
func newHub() *hub {
	return &hub{
		broadcast:  make(chan chatMsg),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

// run hub
func (h *hub) run() {
	for {
		select {
		// on register
		case client := <-h.register:
			h.clients[client] = true

		// on unregister
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		// on broadcast
		case message := <-h.broadcast:
			pipeWriters := make([]*io.PipeWriter, 0)

			// Prepare pipes for every client
			for client := range h.clients {

				// create pipe for every client
				pr, pw := io.Pipe()
				pipeWriters = append(pipeWriters, pw)

				// send chatMsg to client chan
				select {
				case client.send <- chatMsg{messageType: message.messageType, reader: pr}:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

			// pars pipes to io.Writer slice
			writers := make([]io.Writer, 0, len(pipeWriters))
			for _, pw := range pipeWriters {
				writers = append(writers, pw)
			}

			// writer to all pipes writers
			multiWriter := io.MultiWriter(writers...)

			// run copy
			go func() {
				_, err := io.Copy(multiWriter, message.reader)
				if err != nil {
					log.Println("error on copy to multiWriter")
				}

				// Close all pipeWriters
				for _, wr := range pipeWriters {
					err = wr.Close()
					if err != nil {
						log.Println("hub.run: pipeWriter.Close err: ", err)
					}
				}
			}()
		}
	}
}
