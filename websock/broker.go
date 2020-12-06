package websock

import (
	"fmt"
	"log"

	"codenames/rules"
	"codenames/structs"
)

// Broker is the central messaging mechanism for a game instance
type Broker struct {
	Name       string
	Room       structs.Room
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan structs.Word
}

// Newbroker creates a broker for a new game
func Newbroker(name string, room structs.Room) *Broker {
	return &Broker{
		Name:       name,
		Room:       room,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan structs.Word),
	}
}

// Run behaviour for a broker goroutine, websocket channel client attach/remove,
// websocket Rx/Tx & rules processing
func (broker *Broker) Run() {
	for {
		select {
		case client := <-broker.Register:
			broker.Clients[client] = true
			log.Printf("Broker %s: disconnect, size %d ", broker.Name, len(broker.Clients))
			for client := range broker.Clients {
				client.Conn.WriteMessage(1, []byte("Connected"))
			}
			break
		case client := <-broker.Unregister:
			delete(broker.Clients, client)
			log.Printf("Broker %s: disconnect, size %d ", broker.Name, len(broker.Clients))
			for client := range broker.Clients {
				client.Conn.WriteMessage(1, []byte("Disconnected"))
			}
			break
		case move := <-broker.Broadcast:
			log.Printf("Broker: Move received: %+v\n", move)
			for client := range broker.Clients {
				//Process move according to game rules and update state
				rules.ProcessRules(move, broker.Room)
				if err := client.Conn.WriteJSON(broker.Room); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
