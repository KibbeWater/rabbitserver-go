package main

import (
	"log"
	"net/http"
	"os"

	"main/handlers"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

var (
	AppVersion = "undefined"
	OSVersion  = "undefined"
)

func validateEnv() {
	os_version, exists := os.LookupEnv("OS_VERSION")
	if !exists {
		log.Fatal("OS_VERSION environment variable not set")
	}

	app_version, exists := os.LookupEnv("APP_VERSION")
	if !exists {
		log.Fatal("APP_VERSION environment variable not set")
	}

	if os_version == "" {
		log.Fatal("OS_VERSION environment variable is empty")
	}

	if app_version == "" {
		log.Fatal("APP_VERSION environment variable is empty")
	}

	AppVersion = app_version
	OSVersion = os_version
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	// Handle incoming traffic in a goroutine
	go handlers.ServerHandler(ws, OSVersion, AppVersion)
}

func main() {
	port, exists := os.LookupEnv("PORT")

	validateEnv()
	println("OS_VERSION: " + OSVersion)
	println("APP_VERSION: " + AppVersion)

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
