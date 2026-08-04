package main

import (
	"log"
	"net/http"

	"github.com/joaovv-Vitor/phraseforge/internal/httpapi"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func main() {
	const address = ":8080"

	categories, err := storage.LoadCategories("data/phrases.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PhraseForge API listening on %s", address)
	if err := http.ListenAndServe(address, httpapi.NewHandler(categories)); err != nil {
		log.Fatal(err)
	}
}
