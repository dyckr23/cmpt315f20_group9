package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	rejson "github.com/nitishm/go-rejson"

	"github.com/gorilla/mux"

	"codenames/structs"

	"codenames/websock"
)

var pool *redis.Pool

var games map[string]*websock.Broker

//!!!
var devFlag bool

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

		room := structs.Room{}
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
	var words []structs.Word
	for i, word := range values {
		if identityIndices[i] != size-1 {
			words = append(words, structs.Word{word, identities[identityIndices[i]], "false"})
		} else {
			words = append(words, structs.Word{word, firstTeam, "false"})
		}
	}

	// Create Room object
	var room structs.Room
	//room := structs.Room{roomCode, "ongoing", firstTeam, firstTeam, words}
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

	var payload structs.Room
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

func serveWs(w http.ResponseWriter, r *http.Request) {
	var broker *websock.Broker
	wConn, err := websock.Upgrade(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	roomCode := mux.Vars(r)["roomCode"]
	if roomCode == "" {
		err = errors.New("websocket error: missing room code")
		writeJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	rConn := pool.Get()
	defer rConn.Close()

	exists, err := redis.Bool(rConn.Do("exists", roomCode))
	if err != nil {
		writeJSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		err = errors.New("websocket error: cannot find room")
		writeJSONResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(rConn)

	//!!!
	/*if !devFlag {
		devFlag = true
		testJSON, _ := redis.Bytes(rh.JSONGet("test-room-rd", "."))
		var testState structs.Room
		err = json.Unmarshal(testJSON, &testState)
		log.Printf("TEST state loaded: %+v\n", testState)
		broker = websock.Newbroker("test-room-rd", testState)
		games["test-room-rd"] = broker
		go broker.Run()
	}*/
	//!!!

	if _, ok := games[roomCode]; ok {
		log.Println("Found broker!")
		broker = games[roomCode]
	} else {
		roomJSON, err := redis.Bytes(rh.JSONGet(roomCode, "."))
		if err != nil {
			writeJSONResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		roomState := structs.Room{}
		err = json.Unmarshal(roomJSON, &roomState)

		log.Printf("State loaded: %+v\n", roomState)

		broker = websock.Newbroker(roomCode, roomState)
		games[roomCode] = broker
		go broker.Run()
	}

	client := &websock.Client{
		Conn:   wConn,
		Broker: broker,
	}

	broker.Register <- client
	client.Read()
}

func serveGamePage(w http.ResponseWriter, r *http.Request) {
	/*fs := http.FileServer(http.Dir("./htdocs"))
	fmt.Println("Go here!")
	fmt.Println(r.URL.Path)
	r.URL.Path = "game.html"

	fs.ServeHTTP(w, r)*/

	sendTo := "htdocs/game.html"
	http.ServeFile(w, r, sendTo)
}

func serveGame(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("request URI", r.RequestURI)
		log.Println("request URL", r.URL)

		conn := pool.Get()
		defer conn.Close()
		rh := rejson.NewReJSONHandler()
		rh.SetRedigoClient(conn)

		roomCode := r.RequestURI
		roomCode = strings.Trim(roomCode, "/")

		log.Println(roomCode)

		exists, err := redis.Bool(conn.Do("exists", roomCode))
		if err != nil {
			writeJSONResponse(w, err.Error(), 500)
			return
		}

		log.Println(exists)

		if exists {
			fmt.Printf("Room exists: %s\n", roomCode)

			redir := new(url.URL)
			redir.Path = "/game.html"
			r.URL = redir
		}

		next.ServeHTTP(w, r)
	})
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

	games = make(map[string]*websock.Broker)
	router := mux.NewRouter()
	router.Use(middlewareLogWrapper)

	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.HandleFunc("/rooms/{roomCode}", getRoom).Methods("GET")
	//subrouter.HandleFunc("/rooms", makeRoom).Methods("POST")

	router.HandleFunc("/websocket/{roomCode}", serveWs)
	//router.HandleFunc("/{id}", serveGamePage)
	//router.PathPrefix("/test-room-2").Handler(http.StripPrefix("/test-room-2", serveGame))
	router.Use(serveGame)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("htdocs")))

	webserver := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Println("Listening!")
	log.Fatal(webserver.ListenAndServe())
}
