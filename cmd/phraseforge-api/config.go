package main

import (
	"os"
	"strings"
)

const (
	defaultAPIAddress = ":8080"
	defaultDataFile   = "data/phrases.json"
)

func apiAddress() string {
	address := strings.TrimSpace(os.Getenv("PHRASEFORGE_API_ADDR"))
	if address == "" {
		return defaultAPIAddress
	}

	return address
}

func apiDataFile() string {
	path := strings.TrimSpace(os.Getenv("PHRASEFORGE_DATA_FILE"))
	if path == "" {
		return defaultDataFile
	}

	return path
}
