package main

import (
	"log"
	"net/http"
	"os"

	"main/config"
	"main/handlers"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	// Handle incoming traffic in a goroutine
	go handlers.ServerHandler(ws, config.OSVersion, config.AppVersion)
}

func main() {
	port, exists := os.LookupEnv("PORT")

	config.Init()

	if *config.Debug {
		log.Println("Debug mode enabled")
	}

	portNumber := "8080"
	if exists {
		portNumber = port
	}

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server on port " + portNumber)
	log.Fatal(http.ListenAndServe("localhost:"+portNumber, nil))
}
