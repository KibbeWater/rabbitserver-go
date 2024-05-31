package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"main/interfaces"
	"main/rabbit"

	"github.com/Noooste/azuretls-client"
	"github.com/gorilla/websocket"
)

func ServerHandler(ws *websocket.Conn, OSVersion string, AppVersion string) {
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

			go HandleRabbit(rabbitConnection, ws)
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
		case "ptt":
			if rabbitConnection == nil {
				log.Println("ptt received before logon")
				continue
			}

			// Unmarshal the message
			data := interfaces.PTTRequest{}
			err = json.Unmarshal([]byte(message), &data)
			if err != nil {
				log.Println("error unmarshalling ptt data:", err)
				continue
			}

			// Send the message to the rabbit connection
			if data.Data.Active {
				err = rabbitConnection.WriteJSON(map[string]interface{}{
					"kernel": map[string]interface{}{
						"voiceActivity": map[string]interface{}{
							"imageBase64": data.Data.Image,
							"state":       "pttButtonPressed",
						},
					},
				})
			} else {
				err = rabbitConnection.WriteJSON(map[string]interface{}{
					"kernel": map[string]interface{}{
						"voiceActivity": map[string]interface{}{
							"imageBase64": data.Data.Image,
							"state":       "pttButtonReleased",
						},
					},
				})
			}
			if err != nil {
				log.Println("error writing image to rabbit:", err)
				continue
			}
		case "audio":
			if rabbitConnection == nil {
				log.Println("audio received before logon")
				continue
			}

			// Unmarshal the message
			data := interfaces.AudioRequest{}
			err = json.Unmarshal([]byte(message), &data)
			if err != nil {
				log.Println("error unmarshalling audio data:", err)
				continue
			}

			// Convert the base64 audio to a byte array
			audioData, err := base64.StdEncoding.DecodeString(data.Data)
			if err != nil {
				log.Println("error decoding base64 audio:", err)
				continue
			}

			// Send the audio to the rabbit connection
			err = rabbitConnection.WriteMessage(2, audioData)
			if err != nil {
				log.Println("error writing audio to rabbit:", err)
				continue
			}
		case "register":
			if rabbitConnection == nil {
				log.Println("register received before logon")
				continue
			}

			// Unmarshal the message
			registerData := interfaces.RegisterRequest{}
			err = json.Unmarshal([]byte(message), &registerData)
			if err != nil {
				log.Println("error unmarshalling register data:", err)
				continue
			}

			// registerData has a "Data" field which is a base64 encoded string, decode it into a png
			imageData, err := base64.StdEncoding.DecodeString(registerData.Data)
			if err != nil {
				log.Println("error decoding base64 image:", err)
				continue
			}

			url, err := rabbit.DecodeQRAndValidateURL(imageData)
			if err != nil {
				log.Println(err)
				continue
			}

			IMEI := rabbit.GenerateIMEI()
			url += "&deviceId=" + IMEI

			// Perform a GET request to the URL
			body := rabbit.Register(url)

			if strings.Contains(string(body), "\"error\":") {
				log.Println("error registering device:", string(body))
				continue
			}

			// Unmarshal the response
			var registerResponse interfaces.RabbitRegisterResponse
			err = json.Unmarshal(body, &registerResponse)
			if err != nil {
				log.Println("error unmarshalling register response:", err)
				continue
			}

			// Create a new response
			response := interfaces.RegisterResponse{
				Type: "register",
				Data: interfaces.RegisterResponseData{
					ActualUserID: registerResponse.ActualUserID,
					UserID:       registerResponse.UserID,
					AccountKey:   registerResponse.AccountKey,
					UserName:     registerResponse.UserName,
					IMEI:         IMEI,
				},
			}

			// Marshal the response
			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshalling register response:", err)
				continue
			}

			// Write the response
			err = ws.WriteMessage(1, responseBytes)
			if err != nil {
				log.Println("error writing register response:", err)
				continue
			}
		default:
			log.Println("unknown message type:", msg.Type)
		}
	}

}
