package rules

import (
	"fmt"
	"log"
	"strings"

	"codenames/structs"
)

// ProcessRules applies the rules of the game using the word that was chosen
// as a move, to the game state contained within the session broker
func ProcessRules(move structs.Word, game structs.Room) {
	fmt.Printf("Rules got move: %+v\n", move)
	fmt.Printf("In game state: %+v\n", game)

	if strings.Contains(game.Status, "win!") {
		log.Fatalf("Fatal error: referenced finished game %s", game.RoomCode)
	}
	if move.Revealed == "true" {
		log.Fatalf("Fatal error: room %s sent flipped card", game.RoomCode)
	}
	move.Revealed = "true"

	switch move.Identity {
	case "assassin":
		if game.Turn == "blue" {
			game.Status = "red win!"
			return
		}
		game.Status = "blue win!"
		return
	case "spectator":
		if game.Turn == "blue" {
			game.Turn = "red"
			return
		}
		game.Turn = "blue"
		return
	}

}
