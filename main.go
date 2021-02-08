package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}
var room *Room

func checkOrigin(r *http.Request) bool {
	return true
}

func join(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["Name"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	fmt.Printf("Welcome %s\n", name)

	newUser := NewUser(name, conn)
	room.AddUser(newUser)
}

func broadcastLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)

	for range ticker.C {
		room.BroadcastState()
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func main() {
	fmt.Println("Listening on port 8080")

	room = NewRoom()
	go broadcastLoop()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/join/{Name}", join)

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Length", "Date", "X-Content-Type-Options"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST"})

	log.Fatal(http.ListenAndServe(
		*addr,
		handlers.CORS(
			allowedOrigins,
			allowedHeaders,
			allowedMethods,
			handlers.AllowCredentials())(router)))
}
