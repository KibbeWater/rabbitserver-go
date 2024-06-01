package handlers

import (
	"encoding/json"
	"log"
	"main/interfaces"
	"strings"

	"github.com/Noooste/azuretls-client"
	"github.com/gorilla/websocket"
)

func HandleRabbit(rabbit *azuretls.Websocket, ws *websocket.Conn, loggedIn *bool) {
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
			*loggedIn = true
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
		} else if strings.Contains(message, "{\"assistantResponseDevice\":") {
			var assistantResponse interfaces.AssistantDeviceResponse
			err = json.Unmarshal(bytes, &assistantResponse)
			if err != nil {
				log.Println("error unmarshalling assistant device response:", err)
				continue
			}

			response := interfaces.AudioMessageResponse{
				Type: "audio",
				Data: struct {
					Text  string `json:"text"`
					Audio string `json:"audio"`
				}{
					Text:  assistantResponse.Kernel.AssistantResponseDevice.Text,
					Audio: assistantResponse.Kernel.AssistantResponseDevice.Audio,
				},
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshalling audio message response:", err)
				continue
			}

			ws.WriteMessage(1, responseBytes)
		} else if strings.Contains(message, "{\"speechRecognized\":") {
			var speechResponse interfaces.RabbitSpeechResponse
			err = json.Unmarshal(bytes, &speechResponse)
			if err != nil {
				log.Println("error unmarshalling speech response:", err)
				continue
			}

			if !speechResponse.SpeechRecognized.Recognized {
				// Write a message response
				response := interfaces.MessageResponse{
					Type: "message",
					Data: "Sorry, I didn't catch that. Could you repeat it?",
				}
				responseBytes, err := json.Marshal(response)
				if err != nil {
					log.Println("error marshalling message response:", err)
				}

				ws.WriteMessage(1, responseBytes)
				continue
			}
		} else {
			log.Println("unknown message type:", message)
		}
	}
}
