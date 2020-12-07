package datastore

import (
	"log"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"

	"codenames/structs"
)

var pool *redis.Pool

// SetPool gets a redis connection pool for sharing with other packages
func SetPool(extPool *redis.Pool) {
	pool = extPool
}

// GetPool returns a redis connection p
func GetPool() *redis.Pool {
	return pool
}

// UpdateGame is used by the rules processor to save the game state
func UpdateGame(game structs.Room) {
	conn := pool.Get()
	defer conn.Close()
	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	_, err := rh.JSONSet(game.RoomCode, ".", game)
	if err != nil {
		log.Fatalf("Fatal error: server cannot save game %s\n", game.RoomCode)
		return
	}
}

// NewGame creates a new game state when Start New game button is clicked
func NewGame(roomCode string) structs.Room {
	var teams []string = []string{"red", "blue"}
	var identities []string
	var size int = 25

	// Add 8 reds, 8 blues, 7 spectators, and 1 assassin to the identities slice
	for i := 1; i < 9; i++ {
		if i == 1 {
			identities = append(identities, "red", "blue", "assassin")
		} else {
			identities = append(identities, "red", "blue", "spectator")
		}
	}

	conn := pool.Get()
	defer conn.Close()
	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	// Fetch 25 random words from wordlist
	values, err := redis.Strings(conn.Do("SRANDMEMBER", "wordlist", size))
	if err != nil {
		log.Fatalln("Fatal error: could not make new game ", roomCode)
	}

	rand.Seed(time.Now().Unix())
	// Randomly choose the starting team
	firstTeam := teams[rand.Intn(len(teams))]
	// Randomly choose the identities
	identityIndices := rand.Perm(size)
	// Create a list of Word objects
	var words []structs.Word
	for i, text := range values {
		if identityIndices[i] != size-1 {
			word := structs.Word{
				Text:     text,
				Identity: identities[identityIndices[i]],
				Revealed: "false",
			}
			words = append(words, word)
		} else {
			word := structs.Word{
				Text:     text,
				Identity: firstTeam,
				Revealed: "false",
			}
			words = append(words, word)
		}
	}

	// Create Room object
	var room structs.Room
	room.RoomCode = roomCode
	room.Status = "ongoing"
	room.FirstTeam = firstTeam
	room.Turn = firstTeam
	if firstTeam == "blue" {
		room.BlueHidden = 9
		room.RedHidden = 8
	} else {
		room.BlueHidden = 8
		room.RedHidden = 9
	}
	room.Words = words

	return room
}
