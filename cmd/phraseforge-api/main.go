package main

import (
	"log"
	"net/http"

	"github.com/joaovv-Vitor/phraseforge/internal/httpapi"
)

func main() {
	const address = ":8080"

	log.Printf("PhraseForge API listening on %s", address)
	if err := http.ListenAndServe(address, httpapi.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
