package main

import (
	"os"
	"strings"
)

const (
	defaultAPIAddress   = ":8080"
	defaultDatabaseFile = "data/phraseforge.db"
)

func apiAddress() string {
	address := strings.TrimSpace(os.Getenv("PHRASEFORGE_API_ADDR"))
	if address == "" {
		return defaultAPIAddress
	}

	return address
}

func apiDatabaseFile() string {
	path := strings.TrimSpace(os.Getenv("PHRASEFORGE_DATABASE_FILE"))
	if path == "" {
		return defaultDatabaseFile
	}

	return path
}
