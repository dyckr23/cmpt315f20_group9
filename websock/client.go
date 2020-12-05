package websock

import (
	"encoding/json"
	"log"
	"strings"

	"codenames/structs"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Broker *Broker
}

func (c *Client) Read() {
	defer func() {
		c.Broker.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		//message := structs.Message{Type: messageType, Body: string(p)}
		var move structs.Word
		decoder := json.NewDecoder(strings.NewReader(string(p)))
		err = decoder.Decode(&move)

		c.Broker.Broadcast <- move
		log.Printf("Message Received: %+v\n", move)
	}
}
