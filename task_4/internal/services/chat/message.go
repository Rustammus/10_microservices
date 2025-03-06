package chat

import "io"

type chatMsg struct {
	messageType int
	reader      io.Reader
}
