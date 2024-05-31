package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"main/interfaces"
	"main/rabbit"

	"github.com/Noooste/azuretls-client"
	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

var (
	AppVersion = "undefined"
	OSVersion  = "undefined"
)

func handleRabbit(rabbit *azuretls.Websocket, ws *websocket.Conn) {
	for {
		_, bytes, err := rabbit.ReadMessage()
		if err != nil {
			log.Println(err)
			rabbit.Close()
			break
		}

		message := string(bytes)
		if strings.Contains(message, "{\"initialize\"") {
			response := interfaces.LogonResponse{
				Type: "logon",
				Data: "success",
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshalling logon response:", err)
				continue
			}
			err = ws.WriteMessage(websocket.TextMessage, responseBytes)
			if err != nil {
				log.Println("error writing logon response:", err)
				continue
			}
		} else if strings.Contains(message, "{\"assistantResponse\":") {
			// demarshal the message into AssistantResponse
			var assistantResponse interfaces.AssistantResponse
			err = json.Unmarshal(bytes, &assistantResponse)
			if err != nil {
				log.Println("error unmarshalling assistant response:", err)
				continue
			}

			// Create a MessageResponse
			response := interfaces.MessageResponse{
				Type: "message",
				Data: assistantResponse.Kernel.AssistantResponse,
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshalling message response:", err)
				continue
			}

			ws.WriteMessage(1, responseBytes)
		}
	}
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	// Handle incoming traffic in a goroutine
	go func() {
		var rabbitConnection *azuretls.Websocket = nil
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				if rabbitConnection != nil {
					rabbitConnection.Close()
				}
				break
			}

			var msg interfaces.APIRequest
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("error unmarshalling message:", err)
				continue
			}

			switch msg.Type {
			case "logon":
				request := interfaces.LogonRequest{}
				err = json.Unmarshal([]byte(message), &request)
				if err != nil {
					log.Println("error unmarshalling logon data:", err)
					continue
				}

				rabbitConnection = rabbit.SpawnRabbitConnection(request.Data.IMEI, request.Data.AccountKey, OSVersion, AppVersion)

				go handleRabbit(rabbitConnection, ws)
			case "message":
				if rabbitConnection == nil {
					log.Println("message received before logon")
					continue
				}

				// Unmarshal the message
				data := interfaces.MessageRequest{}
				err = json.Unmarshal([]byte(message), &data)
				if err != nil {
					log.Println("error unmarshalling message data:", err)
					continue
				}

				// Send the message to the rabbit connection
				err = rabbitConnection.WriteJSON(rabbit.GenerateMessagePayload(data.Message))
				if err != nil {
					log.Println("error writing message to rabbit:", err)
					continue
				}
			}
		}
	}()
}

func main() {
	port, exists := os.LookupEnv("PORT")

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
