package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	ws := &websocket.Dialer{
		TLSClientConfig: nil,
	}

	conn, resp, err := ws.Dial("ws://localhost:8080/chat/connect?name=gotest", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code", resp.StatusCode)
	}
	defer conn.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wg.Add(2)
	go get(ctx, conn)
	go send(ctx, conn)
	wg.Wait()
}

func get(ctx context.Context, conn *websocket.Conn) {
f:
	for {
		select {
		case <-ctx.Done():
			break f
		default:
			_, r, err := conn.ReadMessage()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(r))
		}
	}
	wg.Done()
}

func send(_ context.Context, conn *websocket.Conn) {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		if err != nil {
			log.Println(err)
			return
		}
	}

	wg.Done()
}
