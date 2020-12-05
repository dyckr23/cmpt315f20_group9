package websock

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Upgrade turns the HTTP request in a Websocket connection
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return wsConn, err
	}
	return wsConn, nil
}
