package main

import "time"

// Room represents a room in which a game takes place
type Room struct {
	RoomCode  string   `json:"roomCode,omitempty"`
	Status    string   `json:"status,omitempty"`
	FirstTeam string   `json:"firstTeam,omitempty"`
	Turn      string   `json:"turn,omitempty"`
	Words     []Word   `json:"words,omitempty"`
	Players   []Player `json:"players,omitempty"`
	Logs      []Log    `json:"logs,omitempty"`
}

// Word represents the 25 words in a game
type Word struct {
	Text     string `json:"text,omitempty"`
	Identity string `json:"identity,omitempty"`
	Revealed bool   `json:"revealed,omitempty"`
}

// Player represents a player in a game
type Player struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}

// Log represents a log message in a game
type Log struct {
	Text      string    `json:"text,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// REF
// https://github.com/gomodule/redigo/blob/master/redis/pubsub.go

// Subscription represents a subscribe or unsubscribe notification
type Subscription struct {
	Type    string
	Channel string
	Subs    int
}
