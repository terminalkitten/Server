package Websocket

import (
	"net/http"
	"github.com/gorilla/websocket"
)

var SocketUpgrade = websocket.Upgrader{
	ReadBufferSize:  1 << 14,
	WriteBufferSize: 1 << 18,

	// It's need to everyone can use the API, even on other website. ;)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start(w http.ResponseWriter, r *http.Request) {
	conn, err := SocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	listen(conn)
}
