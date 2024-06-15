package rabbit

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"main/config"
	"math"
	"strings"
	"time"
)

func getOsVer() string {
	input := config.OSVersion
	parts := strings.Split(input, "_")

	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "_")
	}
	return input
}

func getHealthRand(seed int64) int {
	rand := Srand(uint32(seed))
	for {
		x := rand.Rand() % 100
		if x >= 10 {
			return int(x &^ 1) // round down to nearest even int
		}
	}
}

// Credits to @retr0id for the original py code
// https://gist.github.com/DavidBuchanan314/b6c9102c327f2ba42a3ed374e6ede90f#file-get_health-py
func GetHealth() string {
	OS_VERSION := getOsVer()

	now := time.Now()
	gmt := now.UTC()
	millis := int(math.Mod(float64(now.UnixNano()/1e6), 1000))
	timestamp := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d", gmt.Year(), gmt.Month(), gmt.Day(), gmt.Hour(), gmt.Minute(), gmt.Second(), millis)
	msg := fmt.Sprintf("%s,%s,%d", OS_VERSION, timestamp, getHealthRand(now.Unix()))

	if *config.Debug {
		fmt.Println("Device-Health RAW:", msg)
	}

	block, _ := pem.Decode([]byte(config.HealthPK))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pubkey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pubkey, []byte(msg), nil)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext)
}
