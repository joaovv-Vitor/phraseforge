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

type errorResponse struct {
	Error string `json:"error"`
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
	mux.HandleFunc("/", notFound)

	return mux
}

func randomPhrase(w http.ResponseWriter, r *http.Request, categories []phrase.Category) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	categoryName := "programming"
	categoryRequested := false
	if values, provided := r.URL.Query()["category"]; provided {
		if len(values) != 1 {
			writeJSONError(w, http.StatusBadRequest, "category must be specified exactly once")
			return
		}
		if strings.TrimSpace(values[0]) == "" {
			writeJSONError(w, http.StatusBadRequest, "category cannot be empty")
			return
		}
		categoryName = values[0]
		categoryRequested = true
	}

	count := 1
	if values, provided := r.URL.Query()["count"]; provided {
		if len(values) != 1 {
			writeJSONError(w, http.StatusBadRequest, "count must be specified exactly once")
			return
		}
		if strings.TrimSpace(values[0]) == "" {
			writeJSONError(w, http.StatusBadRequest, "count cannot be empty")
			return
		}

		parsedCount, err := strconv.Atoi(values[0])
		if err != nil || parsedCount < 1 || parsedCount > maxRandomPhraseCount {
			writeJSONError(w, http.StatusBadRequest, "count must be a number between 1 and 10")
			return
		}
		count = parsedCount
	}

	category, err := phrase.FindCategory(categories, categoryName)
	if err != nil {
		if !categoryRequested {
			writeJSONError(w, http.StatusInternalServerError, "failed to generate phrase")
			return
		}
		writeJSONError(w, http.StatusNotFound, "category not found")
		return
	}

	phrases := make([]string, 0, count)
	for range count {
		generated, err := phrase.Generate(category.Template, category.Parts)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to generate phrase")
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
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		return
	}
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	writeJSONError(w, http.StatusNotFound, "not found")
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse{Error: message}); err != nil {
		return
	}
}
