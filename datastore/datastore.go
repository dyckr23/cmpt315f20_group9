package datastore

import (
	"log"

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
