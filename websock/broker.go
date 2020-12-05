package websock

import (
	"codenames/structs"
	"fmt"
)

type Broker struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan structs.Message
}

func Newbroker() *Broker {
	return &Broker{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan structs.Message),
	}
}

func (broker *Broker) Run() {
	for {
		select {
		case client := <-broker.Register:
			broker.Clients[client] = true
			fmt.Println("Number of clients: ", len(broker.Clients))
			for client, _ := range broker.Clients {
				client.Conn.WriteJSON(structs.Message{Type: 1, Body: "New client"})
			}
			break
		case client := <-broker.Unregister:
			delete(broker.Clients, client)
			fmt.Println("Size of Connection broker: ", len(broker.Clients))
			for client, _ := range broker.Clients {
				client.Conn.WriteJSON(structs.Message{Type: 1, Body: "Client disconnect"})
			}
			break
		case message := <-broker.Broadcast:
			fmt.Println("Relaying broadcast message")
			for client, _ := range broker.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
