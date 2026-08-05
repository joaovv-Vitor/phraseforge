package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joaovv-Vitor/phraseforge/internal/httpapi"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func main() {
	categories, err := storage.LoadCategories("data/phrases.json")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              apiAddress(),
		Handler:           httpapi.NewHandler(categories),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("PhraseForge API listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
