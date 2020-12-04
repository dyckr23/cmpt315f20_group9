package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"

	rejson "github.com/nitishm/go-rejson"

	"github.com/gorilla/mux"

	"codenames/websock"
)

var pool *redis.Pool

var teams []string = []string{"red", "blue"}
var identities []string
var size int = 25

func middlewareLogWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source, _, _ := net.SplitHostPort(r.RemoteAddr)
		log.Println("request URI", r.RequestURI, "with method", r.Method, "from ip address", source)
		next.ServeHTTP(w, r)
	})
}

func getTest(w http.ResponseWriter, req *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("get", "test"))
	if err != nil {
		fmt.Println("handle errors better!")
	}

	fmt.Println(val)
}

// getRoom function handles get requests with a room code
// if room exists, return the room's current game state
// if room does not exist, create a new room and return a new game state
func getRoom(w http.ResponseWriter, r *http.Request) {
	// Establish a redis connection
	conn := pool.Get()
	defer conn.Close()
	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	// Get room code from path variables
	roomCode := mux.Vars(r)["roomCode"]
	if roomCode == "" {
		writeJSONResponse(w, http.StatusText(400), 400)
	}

	// Check if room code exists
	exists, err := redis.Bool(conn.Do("exists", roomCode))
	if err != nil {
		writeJSONResponse(w, err.Error(), 500)
		return
	}
	// If room exists, respond with current game state
	if exists {
		fmt.Printf("Room exists: %s\n", roomCode)

		valueJSON, err := redis.Bytes(rh.JSONGet(roomCode, "."))

		if err != nil {
			writeJSONResponse(w, err.Error(), 500)
			return
		}

		room := Room{}
		err = json.Unmarshal(valueJSON, &room)

		if err != nil {
			writeJSONResponse(w, err.Error(), 500)
		}

		json.NewEncoder(w).Encode(room)
		return
	}

	// If room does not exist
	fmt.Printf("Creating new room: %s\n", roomCode)

	// Fetch 25 random words from wordlist
	values, err := redis.Strings(conn.Do("SRANDMEMBER", "wordlist", size))

	if err != nil {
		writeJSONResponse(w, err.Error(), 500)
		return
	}

	rand.Seed(time.Now().Unix())
	// Randomly choose the starting team
	firstTeam := teams[rand.Intn(len(teams))]
	// Randomly choose the identities
	identityIndices := rand.Perm(size)
	// Create a list of Word objects
	var words []Word
	for i, word := range values {
		if identityIndices[i] != size-1 {
			words = append(words, Word{word, identities[identityIndices[i]], "false"})
		} else {
			words = append(words, Word{word, firstTeam, "false"})
		}
	}

	// Create Room object
	room := Room{roomCode, "ongoing", firstTeam, firstTeam, words}
	// Add new room to redis
	_, err = rh.JSONSet(room.RoomCode, ".", room)

	if err != nil {
		writeJSONResponse(w, err.Error(), 500)
		return
	}

	// Send Room object back as a response
	json.NewEncoder(w).Encode(room)

}

func writeJSONResponse(w http.ResponseWriter, message string, code int) {
	fmt.Println(message)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	response := map[string]interface{}{
		"message": message,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
	}
}

func makeRoom(w http.ResponseWriter, r *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	var payload Room
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Printf("Error! %s\n", err.Error())
	}

	fmt.Println(payload)

	// Check for pre-existing id????
	//roomCode := mux.Vars(r)["roomCode"]

	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	res, err := rh.JSONSet(payload.RoomCode, ".", payload)
	if err != nil {
		fmt.Printf("Error jsonset: %s\n", err.Error())
	}

	fmt.Println(res)

	//var test Room
	//decoder = json.NewDecoder(rh.JSONGet(payload.RoomCode, "."))
}

func serveWs(broker *websock.Broker, w http.ResponseWriter, r *http.Request) {
	conn, err := websock.Upgrade(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	client := &websock.Client{
		Conn:   conn,
		Broker: broker,
	}

	broker.Register <- client
	client.Read()
}

func main() {
	// Add 8 reds, 8 blues, 7 spectators, and 1 assassin to the identities slice
	for i := 1; i < 9; i++ {
		if i == 1 {
			identities = append(identities, "red", "blue", "assassin")
		} else {
			identities = append(identities, "red", "blue", "spectator")
		}
	}

	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	broker := websock.Newbroker()
	go broker.Run()

	router := mux.NewRouter()
	router.Use(middlewareLogWrapper)

	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.HandleFunc("/get", getTest).Methods("GET")
	subrouter.HandleFunc("/rooms/{roomCode}", getRoom).Methods("GET")
	subrouter.HandleFunc("/rooms", makeRoom).Methods("POST")

	router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(broker, w, r)
	})

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("htdocs")))

	webserver := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Println("Listening!")
	log.Fatal(webserver.ListenAndServe())
}
