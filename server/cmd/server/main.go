package main

import (
	"log"
	"net/http"

	"github.com/developerabhi99/onlineIDE/api"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", api.MainHandler).Methods("GET")

	r.HandleFunc("/ws", api.WebSocketHandler)

	log.Println("Server running at 8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal("Unable to Serve ", err.Error())
	}

}
