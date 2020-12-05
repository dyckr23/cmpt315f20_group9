package websock

import (
	"log"

	"codenames/structs"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Broker *Broker
}

func (c *Client) Read() {
	defer func() {
		c.Broker.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := structs.Message{Type: messageType, Body: string(p)}
		c.Broker.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)
	}
}
