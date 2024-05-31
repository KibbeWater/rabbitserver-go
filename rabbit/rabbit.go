package rabbit

import (
	"log"

	"github.com/Noooste/azuretls-client"
)

// Old IMEI generation function
/* func calculateChecksum(imeiWithoutChecksum string) int {
	imeiArray := []int{}
	for _, v := range imeiWithoutChecksum {
		imeiArray = append(imeiArray, int(v))
	}

	sum := 0
	double := false
	for _, digit := range imeiArray {
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	checksum := (10 - (sum % 10)) % 10
	return checksum
}

func generateIMEI() string {
	TAC := "35847631"
	serialNumberPrefix := "00"
	serialNumber := serialNumberPrefix
	for i := 0; i < 4; i++ {
		serialNumber += string(rune(48 + rand.Intn(10)))
	}

	imeiWithoutChecksum := TAC + serialNumber
	checksum := calculateChecksum(imeiWithoutChecksum)
	generatedIMEI := imeiWithoutChecksum + strconv.Itoa(checksum)

	return generatedIMEI
} */

func GenerateMessagePayload(message string) map[string]interface{} {
	return map[string]interface{}{
		"kernel": map[string]interface{}{
			"userText": map[string]string{
				"text": message,
			},
		},
	}
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

	ja3 := "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"
	if err := session.ApplyJa3(ja3, azuretls.SchemeWs); err != nil {
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
