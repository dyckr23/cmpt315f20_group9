package structs

import (
	"time"
)

// Room represents a room in which a game takes place
type Room struct {
	RoomCode  string `json:"roomCode,omitempty"`
	Status    string `json:"status,omitempty"`
	FirstTeam string `json:"firstTeam,omitempty"`
	Turn      string `json:"turn,omitempty"`
	Words     []Word `json:"words,omitempty"`
}

// Word represents the 25 words in a game
type Word struct {
	Text     string `json:"text,omitempty"`
	Identity string `json:"identity,omitempty"`
	Revealed string `json:"revealed,omitempty"`
}

// Player represents a player in a game
/*type Player struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}*/

// Log represents a log message in a game
type Log struct {
	Text      string    `json:"text,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// Move represents a turn taken in a game
type Move struct {
	RoomCode string `json:"roomCode,omitempty"`
	Tile     Word   `json:"tile,omitempty"`
}
