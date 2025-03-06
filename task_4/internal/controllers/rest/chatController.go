package rest

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func (s *API) initChatController(r *httprouter.Router) {
	r.GET("/chat/connect", s.chatConnect)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 << 10,
	WriteBufferSize: 1024 << 10,
}

func (s *API) chatConnect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("%s connected", name)
	_ = s.s.Chat.Register(name, conn)
}
