package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/gorilla/mux"
)

var pool *redis.Pool

func getTest(w http.ResponseWriter, req *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("get", "test"))
	if err != nil {
		fmt.Println("handle errors better!")
	}

	fmt.Println(val)
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
	//router.Use(middlewareLogWrapper)
	//subrouter := router.PathPrefix("/api/v1").Subrouter()

	// POST => Create
	//subrouter.HandleFunc("/pastes", pasteCreate).Methods("POST")
	//subrouter.HandleFunc("/pastes/{stub}/reports", reportSubmit).Methods("POST")
	// GET => Read
	//subrouter.HandleFunc("/pastes", pasteRead).Methods("GET")
	// PUT => Update
	//subrouter.HandleFunc("/pastes/{stub}", pasteUpdate).Methods("PUT")
	// DELETE => Delete
	//subrouter.HandleFunc("/pastes/{stub}", pasteDelete).Methods("DELETE")

	//router.HandleFunc("/pastes/{stub}/reports", reportRead).Methods("GET")
	//router.HandleFunc("/pastes/{stub}", tmplRead).Methods("GET")
	//router.HandleFunc("/pastes", pasteBrowse).Methods("GET")
	router.HandleFunc("/get", getTest).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist")))

	webserver := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	log.Println("Listening!")
	log.Fatal(webserver.ListenAndServe())
}
