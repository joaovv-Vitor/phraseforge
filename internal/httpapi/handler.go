// Package httpapi exposes PhraseForge through HTTP handlers.
package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

type healthResponse struct {
	Status string `json:"status"`
}

type categoriesResponse struct {
	Categories []string `json:"categories"`
}

const maxRandomPhraseCount = 10

type randomPhrasesResponse struct {
	Category string   `json:"category"`
	Phrases  []string `json:"phrases"`
}

// NewHandler returns the HTTP handler for the PhraseForge API.
func NewHandler(categories []phrase.Category) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		listCategories(w, r, categories)
	})
	mux.HandleFunc("/phrases/random", func(w http.ResponseWriter, r *http.Request) {
		randomPhrase(w, r, categories)
	})

	return mux
}

func randomPhrase(w http.ResponseWriter, r *http.Request, categories []phrase.Category) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryName := "programming"
	categoryRequested := false
	if values, provided := r.URL.Query()["category"]; provided {
		if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
			http.Error(w, "category must be specified once and cannot be empty", http.StatusBadRequest)
			return
		}
		categoryName = values[0]
		categoryRequested = true
	}

	count := 1
	if values, provided := r.URL.Query()["count"]; provided {
		if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
			http.Error(w, "count must be specified once and cannot be empty", http.StatusBadRequest)
			return
		}

		parsedCount, err := strconv.Atoi(values[0])
		if err != nil || parsedCount < 1 || parsedCount > maxRandomPhraseCount {
			http.Error(w, "count must be a number between 1 and 10", http.StatusBadRequest)
			return
		}
		count = parsedCount
	}

	category, err := phrase.FindCategory(categories, categoryName)
	if err != nil {
		if !categoryRequested {
			http.Error(w, "failed to generate phrase", http.StatusInternalServerError)
			return
		}
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	phrases := make([]string, 0, count)
	for range count {
		generated, err := phrase.Generate(category.Template, category.Parts)
		if err != nil {
			http.Error(w, "failed to generate phrase", http.StatusInternalServerError)
			return
		}
		phrases = append(phrases, generated)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(randomPhrasesResponse{
		Category: category.Name,
		Phrases:  phrases,
	}); err != nil {
		return
	}
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
