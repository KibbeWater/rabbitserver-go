package rabbit

import (
	"log"

	"github.com/Noooste/azuretls-client"
)

const JA3 = "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"

func GenerateMessagePayload(message string) map[string]interface{} {
	return map[string]interface{}{
		"kernel": map[string]interface{}{
			"userText": map[string]string{
				"text": message,
			},
		},
	}
}

func Register(URL string) []byte {
	session := azuretls.NewSession()

	if err := session.ApplyJa3(JA3, azuretls.SchemeHttps); err != nil {
		panic(err)
	}

	resp, err := session.Get(URL, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return resp.Body
}

func SpawnRabbitConnection(IMEI string, accountKey string, osVer string, appVer string) *azuretls.Websocket {
	authPayload := map[string]interface{}{
		"global": map[string]interface{}{
			"initialize": map[string]interface{}{
				"deviceId":  IMEI,
				"evaluate":  false,
				"greet":     true,
				"language":  "en",
				"listening": true,
				"location": map[string]float64{
					"latitude":  0.0,
					"longitude": 0.0,
				},
				"mimeType": "wav",
				"timeZone": "GMT",
				"token":    "rabbit-account-key+" + accountKey,
			},
		},
	}

	session := azuretls.NewSession()

	if err := session.ApplyJa3(JA3, azuretls.SchemeWss); err != nil {
		panic(err)
	}

	ws, err := session.NewWebsocket("wss://r1-api.rabbit.tech/session", 1024, 1024,
		azuretls.OrderedHeaders{
			{"deviceId", IMEI},
			{"App-Version", appVer},
			{"OS-Version", osVer},
			{"User-Agent", "okhttp/4.9.1"},
		},
	)

	if err != nil {
		log.Println(err)
		return nil
	}

	// Send the auth payload
	if err := ws.WriteJSON(authPayload); err != nil {
		log.Println(err)
		return nil
	}

	return ws
}
