package main

import (
	"os"
	"strings"
)

const defaultAPIAddress = ":8080"

func apiAddress() string {
	address := strings.TrimSpace(os.Getenv("PHRASEFORGE_API_ADDR"))
	if address == "" {
		return defaultAPIAddress
	}

	return address
}
