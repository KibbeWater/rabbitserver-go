package config

import (
	"flag"
	"log"
	"os"
)

var Debug = flag.Bool("D", false, "Enable debug mode")

var (
	AppVersion = "undefined"
	OSVersion  = "undefined"
)

func Init() {
	flag.Parse()

	validateEnv()
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
