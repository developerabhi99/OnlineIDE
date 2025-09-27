package main

import (
	"log"
	"net/http"

	"github.com/developerabhi99/onlineIDE/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", api.MainHandler).Methods("GET")

	r.HandleFunc("/ws", api.WebSocketHandler)
	r.HandleFunc("/files", api.GetFileTree)
	r.HandleFunc("/fileWatcher", api.FileTreeWatcher)
	r.HandleFunc("/fileCode/{path:.*}", api.GetFileCode) //{path:.*}-> this tell mux param may contain slashes.
	r.HandleFunc("/saveFile", api.SaveFileCode).Methods("POST")

	//confiruging cors

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React app URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Server running at 8080")
	err := http.ListenAndServe(":8080", handler)

	if err != nil {
		log.Fatal("Unable to Serve ", err.Error())
	}

}
