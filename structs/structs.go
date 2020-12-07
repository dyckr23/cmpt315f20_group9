package structs

// Room represents a room in which a game takes place
type Room struct {
	RoomCode   string `json:"roomCode,omitempty"`
	Status     string `json:"status,omitempty"`
	FirstTeam  string `json:"firstTeam,omitempty"`
	Turn       string `json:"turn,omitempty"`
	BlueHidden int    `json:"blueHidden"`
	RedHidden  int    `json:"redHidden"`
	Words      []Word `json:"words,omitempty"`
}

// Word represents one of the 25 words in a game
type Word struct {
	Text     string `json:"text,omitempty"`
	Identity string `json:"identity,omitempty"`
	Revealed string `json:"revealed,omitempty"`
}
