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
	log.Println()
	fmt.Printf("Rules got move: %+v\n", move)
	log.Println()
	fmt.Printf("In game state: %+v\n", game)
	log.Println()

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
	case "blue":
		game.BlueHidden--
		if game.BlueHidden == 0 {
			game.Status = "blue win!"
			return
		} else if game.Turn == "blue" {
			return
		} else if game.Turn == "red" {
			game.Turn = "blue"
			return
		}
	case "red":
		game.RedHidden--
		if game.RedHidden == 0 {
			game.Status = "red win!"
			return
		} else if game.Turn == "red" {
			return
		} else if game.Turn == "blue" {
			game.Turn = "red"
			return
		}
	}
}
