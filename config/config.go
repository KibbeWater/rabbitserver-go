package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var Debug = flag.Bool("D", false, "Enable debug mode")

var (
	AppVersion = "undefined"
	OSVersion  = "undefined"
	HealthPK   = "undefined"
)

func Init() {
	flag.Parse()

	validateEnv()
	getHealthPublicKey()
}

func getHealthPublicKey() {
	// Get executable path excluding the executable name
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath)
	keyDir := filepath.Join(exeDir, "key.pub")

	// Get our health check RSA public key
	// Check if key.pub exists in executable directory
	_, err = os.Stat(keyDir)
	if err != nil {
		log.Fatal("key.pub not found")
	}

	// Read key.pub
	file, err := os.Open(keyDir)
	if err != nil {
		log.Fatal("Failed to open key.pub")
	}
	defer file.Close()

	stat, _ := file.Stat()
	key := make([]byte, stat.Size())
	_, err = file.Read(key)
	if err != nil {
		log.Fatal("Failed to read key.pub")
	}

	HealthPK = string(key)
}

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
