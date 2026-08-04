// Package httpapi exposes PhraseForge through HTTP handlers.
package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

type healthResponse struct {
	Status string `json:"status"`
}

type categoriesResponse struct {
	Categories []string `json:"categories"`
}

// NewHandler returns the HTTP handler for the PhraseForge API.
func NewHandler(categories []phrase.Category) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		listCategories(w, r, categories)
	})

	return mux
}

func listCategories(w http.ResponseWriter, r *http.Request, categories []phrase.Category) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categoriesResponse{Categories: phrase.CategoryNames(categories)}); err != nil {
		return
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		return
	}
}
