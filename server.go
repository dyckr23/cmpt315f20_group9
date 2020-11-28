package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"

	rejson "github.com/nitishm/go-rejson"

	"github.com/gorilla/mux"
)

var pool *redis.Pool

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
// If room does not exist --> create new room --> response with Room as json
// If room exists --> check room's status --> response with Room as json, or 403
func getRoom(w http.ResponseWriter, r *http.Request) {
	// Log request
	logRequest(w, r)

	// Establish a redis connection
	conn := pool.Get()
	defer conn.Close()

	// Get id from path variables
	id := mux.Vars(r)["id"]
	if id == "" {
		writeJSONResponse(w, http.StatusText(400), 400)
	}

	//rh := rejson.NewReJSONHandler()
	// rh.JSON

}

func makeRoom(w http.ResponseWriter, r *http.Request) {
	logRequest(w, r)

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

func main() {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	router := mux.NewRouter()
	router.Use(middlewareLogWrapper)

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	subrouter.HandleFunc("/get", getTest).Methods("GET")

	subrouter.HandleFunc("/rooms/{roomCode}", getRoom).Methods("GET")

	subrouter.HandleFunc("/rooms", makeRoom).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("htdocs")))

	webserver := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Listening!")
	log.Fatal(webserver.ListenAndServe())
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.RequestURI)
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
