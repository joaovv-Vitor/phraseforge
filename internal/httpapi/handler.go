// Package httpapi exposes PhraseForge through HTTP handlers.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type favoriteRequest struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

type favoriteResponse struct {
	ID        int64     `json:"id"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type favoritesResponse struct {
	Favorites []favoriteResponse `json:"favorites"`
}

// FavoriteStore provides the favorite operations required by the HTTP API.
type FavoriteStore interface {
	Create(context.Context, string, string) (phrase.Favorite, error)
	List(context.Context) ([]phrase.Favorite, error)
}

// HistoryStore records generated phrases required by the HTTP API.
type HistoryStore interface {
	Record(context.Context, string, []string) error
}

// NewHandler returns the HTTP handler for the PhraseForge API.
func NewHandler(categories []phrase.Category, favorites FavoriteStore, history HistoryStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		listCategories(w, r, categories)
	})
	mux.HandleFunc("/phrases/random", func(w http.ResponseWriter, r *http.Request) {
		randomPhrase(w, r, categories, history)
	})
	mux.HandleFunc("/favorites", func(w http.ResponseWriter, r *http.Request) {
		favoritesHandler(w, r, favorites)
	})
	mux.HandleFunc("/", notFound)

	return mux
}

func favoritesHandler(w http.ResponseWriter, r *http.Request, favorites FavoriteStore) {
	switch r.Method {
	case http.MethodGet:
		listFavorites(w, r, favorites)
	case http.MethodPost:
		createFavorite(w, r, favorites)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func createFavorite(w http.ResponseWriter, r *http.Request, favorites FavoriteStore) {
	request, err := decodeFavoriteRequest(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if favorites == nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create favorite")
		return
	}

	favorite, err := favorites.Create(r.Context(), request.Category, request.Content)
	if errors.Is(err, phrase.ErrFavoriteCategoryNotFound) {
		writeJSONError(w, http.StatusNotFound, "category not found")
		return
	}
	if errors.Is(err, phrase.ErrFavoriteAlreadyExists) {
		writeJSONError(w, http.StatusConflict, "favorite already exists")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create favorite")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toFavoriteResponse(favorite)); err != nil {
		return
	}
}

func listFavorites(w http.ResponseWriter, r *http.Request, favorites FavoriteStore) {
	if favorites == nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to list favorites")
		return
	}

	items, err := favorites.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to list favorites")
		return
	}

	response := favoritesResponse{Favorites: make([]favoriteResponse, 0, len(items))}
	for _, favorite := range items {
		response.Favorites = append(response.Favorites, toFavoriteResponse(favorite))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func decodeFavoriteRequest(r *http.Request) (favoriteRequest, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request favoriteRequest
	if err := decoder.Decode(&request); err != nil {
		return favoriteRequest{}, fmt.Errorf("invalid favorite request")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return favoriteRequest{}, fmt.Errorf("invalid favorite request")
	}

	request.Category = strings.TrimSpace(request.Category)
	request.Content = strings.TrimSpace(request.Content)
	if request.Category == "" {
		return favoriteRequest{}, fmt.Errorf("category cannot be empty")
	}
	if request.Content == "" {
		return favoriteRequest{}, fmt.Errorf("content cannot be empty")
	}

	return request, nil
}

func toFavoriteResponse(favorite phrase.Favorite) favoriteResponse {
	return favoriteResponse{
		ID:        favorite.ID,
		Category:  favorite.Category,
		Content:   favorite.Content,
		CreatedAt: favorite.CreatedAt,
	}
}

func randomPhrase(w http.ResponseWriter, r *http.Request, categories []phrase.Category, history HistoryStore) {
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
	if history == nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to record generation history")
		return
	}
	if err := history.Record(r.Context(), category.Name, phrases); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to record generation history")
		return
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
