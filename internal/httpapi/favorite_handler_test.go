package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func TestFavoriteEndpoints(t *testing.T) {
	createdAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	favorite := phrase.Favorite{
		ID:        1,
		Category:  "programming",
		Content:   "Codigo simples reduz problemas futuros.",
		CreatedAt: createdAt,
	}

	tests := []struct {
		name       string
		method     string
		body       string
		store      fakeFavoriteStore
		wantStatus int
		wantError  string
		wantAllow  string
		wantOne    *favoriteResponse
		wantMany   []favoriteResponse
	}{
		{
			name:   "creates favorite",
			method: http.MethodPost,
			body:   `{"category":"programming","content":"Codigo simples reduz problemas futuros."}`,
			store: fakeFavoriteStore{create: func(_ context.Context, category, content string) (phrase.Favorite, error) {
				if category != favorite.Category || content != favorite.Content {
					t.Errorf("Create() received category = %q, content = %q", category, content)
				}
				return favorite, nil
			}},
			wantStatus: http.StatusCreated,
			wantOne:    &favoriteResponse{ID: 1, Category: favorite.Category, Content: favorite.Content, CreatedAt: createdAt},
		},
		{
			name:       "lists favorites",
			method:     http.MethodGet,
			store:      fakeFavoriteStore{favorites: []phrase.Favorite{favorite}},
			wantStatus: http.StatusOK,
			wantMany:   []favoriteResponse{{ID: 1, Category: favorite.Category, Content: favorite.Content, CreatedAt: createdAt}},
		},
		{
			name:       "rejects invalid JSON",
			method:     http.MethodPost,
			body:       `{"category":`,
			store:      fakeFavoriteStore{},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid favorite request",
		},
		{
			name:       "rejects empty category",
			method:     http.MethodPost,
			body:       `{"category":" ","content":"Codigo simples reduz problemas futuros."}`,
			store:      fakeFavoriteStore{},
			wantStatus: http.StatusBadRequest,
			wantError:  "category cannot be empty",
		},
		{
			name:       "rejects empty content",
			method:     http.MethodPost,
			body:       `{"category":"programming","content":" "}`,
			store:      fakeFavoriteStore{},
			wantStatus: http.StatusBadRequest,
			wantError:  "content cannot be empty",
		},
		{
			name:       "returns not found for unknown category",
			method:     http.MethodPost,
			body:       `{"category":"unknown","content":"Codigo simples reduz problemas futuros."}`,
			store:      fakeFavoriteStore{createErr: phrase.ErrFavoriteCategoryNotFound},
			wantStatus: http.StatusNotFound,
			wantError:  "category not found",
		},
		{
			name:       "returns conflict for duplicate favorite",
			method:     http.MethodPost,
			body:       `{"category":"programming","content":"Codigo simples reduz problemas futuros."}`,
			store:      fakeFavoriteStore{createErr: phrase.ErrFavoriteAlreadyExists},
			wantStatus: http.StatusConflict,
			wantError:  "favorite already exists",
		},
		{
			name:       "returns internal error for repository failure",
			method:     http.MethodGet,
			store:      fakeFavoriteStore{listErr: errors.New("database unavailable")},
			wantStatus: http.StatusInternalServerError,
			wantError:  "failed to list favorites",
		},
		{
			name:       "rejects unsupported method",
			method:     http.MethodDelete,
			store:      fakeFavoriteStore{},
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  "GET, POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(testCategories(), &tt.store)
			request := httptest.NewRequest(tt.method, "/favorites", strings.NewReader(tt.body))
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if tt.wantAllow != "" && recorder.Header().Get("Allow") != tt.wantAllow {
				t.Errorf("Allow = %q, want %q", recorder.Header().Get("Allow"), tt.wantAllow)
			}
			if tt.wantError != "" {
				var response errorResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if response.Error != tt.wantError {
					t.Errorf("error body = %q, want %q", response.Error, tt.wantError)
				}
				return
			}

			if tt.wantOne != nil {
				var response favoriteResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode favorite response: %v", err)
				}
				if response != *tt.wantOne {
					t.Errorf("favorite response = %#v, want %#v", response, *tt.wantOne)
				}
				return
			}

			var response favoritesResponse
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decode favorites response: %v", err)
			}
			if len(response.Favorites) != len(tt.wantMany) {
				t.Fatalf("favorites count = %d, want %d", len(response.Favorites), len(tt.wantMany))
			}
			for index := range tt.wantMany {
				if response.Favorites[index] != tt.wantMany[index] {
					t.Errorf("favorite %d = %#v, want %#v", index, response.Favorites[index], tt.wantMany[index])
				}
			}
		})
	}
}

type fakeFavoriteStore struct {
	create    func(context.Context, string, string) (phrase.Favorite, error)
	createErr error
	favorites []phrase.Favorite
	listErr   error
}

func (store *fakeFavoriteStore) Create(ctx context.Context, category, content string) (phrase.Favorite, error) {
	if store.create != nil {
		return store.create(ctx, category, content)
	}

	return phrase.Favorite{}, store.createErr
}

func (store *fakeFavoriteStore) List(context.Context) ([]phrase.Favorite, error) {
	if store.listErr != nil {
		return nil, store.listErr
	}

	return store.favorites, nil
}
