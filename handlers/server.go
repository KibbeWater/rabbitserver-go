package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"main/config"
	"main/interfaces"
	"main/rabbit"

	"github.com/go-audio/wav"
	"github.com/gorilla/websocket"
)

func ServerHandler(ws *websocket.Conn, OSVersion string, AppVersion string, imei string) {
	var isLoggedIn bool = false

	rabbitConnection, err := rabbit.SpawnRabbitConnection(OSVersion, AppVersion, imei)

	if err != nil {
		panic(fmt.Sprintf("Error connecting to the RabbitHole, please check and verify that the OSVersion and AppVersion headers are correct. Error: %v", err))
	}

	if *config.Debug {
		log.Println("Connected to RabbitHole at " + config.URL)
	}

	go HandleRabbit(rabbitConnection, ws, &isLoggedIn)

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

		if *config.Debug {
			println("received message from client:", msg.Type)
			println("message:", string(message))
		}

		switch msg.Type {
		case "logon":
			if isLoggedIn {
				log.Println("logon received after already logged in")
				continue
			}

			request := interfaces.LogonRequest{}
			err = json.Unmarshal([]byte(message), &request)
			if err != nil {
				log.Println("error unmarshalling logon data:", err)
				continue
			}

			// Send the logon request to the rabbit connection
			err = rabbitConnection.WriteJSON(rabbit.GenerateAuthPayload(request.Data.IMEI, request.Data.AccountKey))
			if err != nil {
				log.Println("error writing logon to rabbit:", err)
				continue
			}

		case "message":
			if !isLoggedIn {
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
			if !isLoggedIn {
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
			var responseJSON map[string]interface{}
			if data.Data.Active {
				responseJSON = map[string]interface{}{
					"kernel": map[string]interface{}{
						"voiceActivity": map[string]interface{}{
							"imageBase64": data.Data.Image,
							"state":       "pttButtonPressed",
						},
					},
				}
			} else {
				responseJSON = map[string]interface{}{
					"kernel": map[string]interface{}{
						"voiceActivity": map[string]interface{}{
							"imageBase64": data.Data.Image,
							"state":       "pttButtonReleased",
						},
					},
				}
			}

			// Marshal the response
			responseBytes, err := json.Marshal(responseJSON)
			if err != nil {
				log.Println("error marshalling ptt response:", err)
				continue
			}

			// Write the response
			err = rabbitConnection.WriteMessage(1, responseBytes)
			if err != nil {
				log.Println("error writing ptt response:", err)
				continue
			}
		case "audio":
			if !isLoggedIn {
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

			// Get the base64 data from the data url
			b64 := strings.Split(data.Data, ",")[1]

			// Convert the base64 audio to a byte array
			audioData, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				log.Println("error decoding base64 audio:", err)
				continue
			}

			// audioData is a data URL, if the mimeType is audio/webm, convert it to a WAV, otherwise send it as is
			/* if strings.Contains(data.Data, "audio/webm") {
				audioData, err = rabbit.WebMToWav(b64)
				if err != nil {
					log.Println("error converting audio to WAV:", err)
					continue
				}
			} else {
				continue
			} */

			// Check if audio data is valid WAV
			reader := bytes.NewReader(audioData)
			decoder := wav.NewDecoder(reader)

			_, err = decoder.FullPCMBuffer()
			if err != nil {
				log.Println("error decoding WAV audio:", err)
				continue
			}

			// Send the audio to the rabbit connection
			if *config.Debug {
				println("sending audio to rabbit")
			}
			err = rabbitConnection.WriteMessage(2, audioData)
			if err != nil {
				log.Println("error writing audio to rabbit:", err)
				continue
			}
		case "meeting":
			if !isLoggedIn {
				log.Println("meeting received before logon")
				continue
			}

			// Unmarshal the message
			data := interfaces.MeetingRequest{}
			err = json.Unmarshal([]byte(message), &data)
			if err != nil {
				log.Println("error unmarshalling meeting data:", err)
				continue
			}

			// Send the message to the rabbit connection
			err = rabbitConnection.WriteMessage(1, []byte(fmt.Sprintf("{\"meetingAssistant\":{\"endRecording\":{\"cancel\":%t}}}", data.Data)))
			if err != nil {
				log.Println("error writing meeting to rabbit:", err)
				continue
			}

			// Send back a meeting message
			response := interfaces.MeetingMessageData{
				Type: "meeting",
				Data: false,
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshalling meeting response:", err)
				continue
			}

			err = ws.WriteMessage(1, responseBytes)
			if err != nil {
				log.Println("error writing meeting response:", err)
				continue
			}
		case "raw":
			// Unmarshal the message
			data := interfaces.RAWRequest{}
			err = json.Unmarshal([]byte(message), &data)
			if err != nil {
				log.Println("error unmarshalling ptt data:", err)
				continue
			}

			// Send the message to the rabbit connection
			err = rabbitConnection.WriteMessage(1, []byte(data.Data))
			if err != nil {
				log.Println("error writing raw to rabbit:", err)
				continue
			}
		case "register":
			if isLoggedIn {
				log.Println("register received after already logged in")
				continue
			}

			// Unmarshal the message
			registerData := interfaces.RegisterRequest{}
			err = json.Unmarshal([]byte(message), &registerData)
			if err != nil {
				log.Println("error unmarshalling register data:", err)
				continue
			}

			urlStr := ""
			if strings.Contains(registerData.Data, "https://") {
				qrUrl, err := url.ParseRequestURI(registerData.Data)
				if err != nil {
					continue
				}

				if qrUrl.Scheme != "https" || qrUrl.Host != "hole.rabbit.tech" || qrUrl.Path != "/apis/linkDevice" {
					continue
				}

				urlStr = registerData.Data
			}

			// registerData has a "Data" field which is a base64 encoded string, decode it into a png
			if urlStr == "" {
				imageData, err := base64.StdEncoding.DecodeString(registerData.Data)
				if err != nil {
					log.Println("error decoding base64 image:", err)
					continue
				}

				urlStr, err = rabbit.DecodeQRAndValidateURL(imageData)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			if urlStr == "" {
				log.Println("error decoding QR code")
				continue
			}

			IMEI := rabbit.GenerateIMEI()
			urlStr += "&deviceId=" + IMEI

			// Perform a GET request to the URL
			body := rabbit.Register(urlStr)

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
			if *config.Debug {
				log.Println("unknown message type:", msg.Type)
			}
		}
	}

}
