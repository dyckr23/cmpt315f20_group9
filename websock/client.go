package websock

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"

	"codenames/structs"
)

// Client is a websocket remote endpoint, a frontend
type Client struct {
	Conn   *websocket.Conn
	Broker *Broker
}

// Clients attached to concurrent brokers will listen for and read messages
// sent to them from the frontend
func (c *Client) Read() {
	defer func() {
		c.Broker.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, buf, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var move structs.Word
		decoder := json.NewDecoder(strings.NewReader(string(buf)))
		err = decoder.Decode(&move)

		c.Broker.Broadcast <- move
		//log.Printf("Client: Move received: %+v\n", move)
	}
}
