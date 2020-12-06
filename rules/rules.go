package rules

import (
	"log"
	"strings"

	"codenames/structs"
)

// ProcessRules applies the rules of the game using the word that was chosen
// as a move, to the game state contained within the session broker
func ProcessRules(move structs.Word, game structs.Room) structs.Room {
	if strings.Contains(game.Status, "win!") {
		//log.Fatalf("Fatal error: referenced finished game %s", game.RoomCode)
		log.Printf("Error: referenced finished game %s", game.RoomCode)
	}
	if move.Revealed == "true" {
		//log.Fatalf("Fatal error: room %s sent flipped card", game.RoomCode)
		log.Printf("Error: room %s sent flipped card", game.RoomCode)
	}
	// check for end turn signal
	if move.Text == "end turn" && move.Identity == "control" {
		if game.Turn == "blue" {
			game.Turn = "red"
			return game
		}
		game.Turn = "blue"
		return game
	}

	for i, v := range game.Words {
		if v.Text == move.Text {
			v.Revealed = "true"
			game.Words[i] = v
			break
		}
	}

	switch move.Identity {
	case "assassin":
		if game.Turn == "blue" {
			game.Status = "red win!"
			return game
		}
		game.Status = "blue win!"
		return game
	case "spectator":
		if game.Turn == "blue" {
			game.Turn = "red"
			return game
		}
		game.Turn = "blue"
		return game
	case "blue":
		game.BlueHidden--
		if game.BlueHidden == 0 {
			game.Status = "blue win!"
			return game
		} else if game.Turn == "blue" {
			return game
		} else if game.Turn == "red" {
			game.Turn = "blue"
			return game
		}
	case "red":
		game.RedHidden--
		if game.RedHidden == 0 {
			game.Status = "red win!"
			return game
		} else if game.Turn == "red" {
			return game
		} else if game.Turn == "blue" {
			game.Turn = "red"
			return game
		}
	}
	return game
}
