package rabbit

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"main/config"
	"net/url"

	"github.com/Noooste/azuretls-client"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

const JA3 = "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"

func DecodeQRAndValidateURL(imageData []byte) (string, error) {
	// Create a reader with the QRCodeReader
	reader := qrcode.NewQRCodeReader()

	// Create a BinaryBitmap from your image data
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("error decoding image data: %w", err)
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("error creating binary bitmap: %w", err)
	}

	// Decode the QR code
	result, err := reader.Decode(bitmap, nil)
	if err != nil {
		return "", fmt.Errorf("error decoding QR code: %w", err)
	}

	// Check if the decoded information is a valid URL
	qrUrl, err := url.ParseRequestURI(result.GetText())
	if err != nil {
		return "", fmt.Errorf("decoded QR code is not a valid URL: %w", err)
	}

	// Does URL point to https://hole.rabbit.tech/apis/linkDevice
	if qrUrl.Scheme != "https" || qrUrl.Host != "hole.rabbit.tech" || qrUrl.Path != "/apis/linkDevice" {
		return "", fmt.Errorf("decoded QR code does not point to https://hole.rabbit.tech/apis/linkDevice")
	}

	// If we reach here, the QR code contains a valid URL
	return result.GetText(), nil
}

func GenerateMessagePayload(message string) map[string]interface{} {
	return map[string]interface{}{
		"kernel": map[string]interface{}{
			"userText": map[string]string{
				"text": message,
			},
		},
	}
}

func GenerateAuthPayload(IMEI string, accountKey string) map[string]interface{} {
	return map[string]interface{}{
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

func SpawnRabbitConnection(osVer string, appVer string) (*azuretls.Websocket, error) {
	session := azuretls.NewSession()

	if err := session.ApplyJa3(JA3, azuretls.SchemeWss); err != nil {
		panic(err)
	}

	ws, err := session.NewWebsocket(config.URL, 1024, 1024,
		azuretls.OrderedHeaders{
			{"App-Version", appVer},
			{"OS-Version", osVer},
			{"Device-Health", GetHealth()},
			{"User-Agent", "okhttp/4.11.0"},
		},
	)

	if err != nil {
		return nil, err
	}

	return ws, nil
}
