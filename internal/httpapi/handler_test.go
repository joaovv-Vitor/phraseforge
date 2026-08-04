package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantJSON     bool
		wantNames    []string
		wantCategory string
		wantPhrases  []string
		wantError    string
		wantAllow    bool
	}{
		{
			name:         "returns a random programming phrase",
			method:       http.MethodGet,
			path:         "/phrases/random",
			wantStatus:   http.StatusOK,
			wantJSON:     true,
			wantCategory: "programming",
			wantPhrases:  []string{"Codigo simples reduz problemas futuros."},
		},
		{
			name:         "returns multiple random programming phrases",
			method:       http.MethodGet,
			path:         "/phrases/random?count=3",
			wantStatus:   http.StatusOK,
			wantJSON:     true,
			wantCategory: "programming",
			wantPhrases: []string{
				"Codigo simples reduz problemas futuros.",
				"Codigo simples reduz problemas futuros.",
				"Codigo simples reduz problemas futuros.",
			},
		},
		{
			name:         "returns a random phrase from selected category",
			method:       http.MethodGet,
			path:         "/phrases/random?category=study&count=3",
			wantStatus:   http.StatusOK,
			wantJSON:     true,
			wantCategory: "study",
			wantPhrases: []string{
				"A pratica constante fortalece o aprendizado.",
				"A pratica constante fortalece o aprendizado.",
				"A pratica constante fortalece o aprendizado.",
			},
		},
		{
			name:       "returns bad request for empty count",
			method:     http.MethodGet,
			path:       "/phrases/random?count=",
			wantStatus: http.StatusBadRequest,
			wantError:  "count cannot be empty",
		},
		{
			name:       "returns bad request for invalid count",
			method:     http.MethodGet,
			path:       "/phrases/random?count=invalid",
			wantStatus: http.StatusBadRequest,
			wantError:  "count must be a number between 1 and 10",
		},
		{
			name:       "returns bad request for zero count",
			method:     http.MethodGet,
			path:       "/phrases/random?count=0",
			wantStatus: http.StatusBadRequest,
			wantError:  "count must be a number between 1 and 10",
		},
		{
			name:       "returns bad request for negative count",
			method:     http.MethodGet,
			path:       "/phrases/random?count=-1",
			wantStatus: http.StatusBadRequest,
			wantError:  "count must be a number between 1 and 10",
		},
		{
			name:       "returns bad request for count above limit",
			method:     http.MethodGet,
			path:       "/phrases/random?count=11",
			wantStatus: http.StatusBadRequest,
			wantError:  "count must be a number between 1 and 10",
		},
		{
			name:       "returns bad request for repeated count",
			method:     http.MethodGet,
			path:       "/phrases/random?count=2&count=3",
			wantStatus: http.StatusBadRequest,
			wantError:  "count must be specified exactly once",
		},
		{
			name:       "returns bad request for empty category",
			method:     http.MethodGet,
			path:       "/phrases/random?category=",
			wantStatus: http.StatusBadRequest,
			wantError:  "category cannot be empty",
		},
		{
			name:       "returns not found for unknown category",
			method:     http.MethodGet,
			path:       "/phrases/random?category=unknown",
			wantStatus: http.StatusNotFound,
			wantError:  "category not found",
		},
		{
			name:       "rejects unsupported random phrase method",
			method:     http.MethodPost,
			path:       "/phrases/random",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  true,
		},
		{
			name:       "returns categories in configured order",
			method:     http.MethodGet,
			path:       "/categories",
			wantStatus: http.StatusOK,
			wantJSON:   true,
			wantNames:  []string{"programming", "study"},
		},
		{
			name:       "rejects unsupported categories method",
			method:     http.MethodPost,
			path:       "/categories",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  true,
		},
		{
			name:       "returns health status",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
			wantJSON:   true,
		},
		{
			name:       "rejects unsupported health method",
			method:     http.MethodPost,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  true,
		},
		{
			name:       "returns not found for unknown route",
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
	}

	handler := NewHandler(testCategories())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if tt.wantError != "" {
				if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
					t.Fatalf("Content-Type = %q, want application/json", contentType)
				}
				if tt.wantAllow && recorder.Header().Get("Allow") != http.MethodGet {
					t.Errorf("Allow = %q, want %q", recorder.Header().Get("Allow"), http.MethodGet)
				}

				var response errorResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if response.Error != tt.wantError {
					t.Errorf("error body = %q, want %q", response.Error, tt.wantError)
				}
				return
			}
			if !tt.wantJSON {
				return
			}
			if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
				t.Fatalf("Content-Type = %q, want application/json", contentType)
			}

			switch request.URL.Path {
			case "/health":
				var response healthResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if response.Status != "ok" {
					t.Errorf("status body = %q, want %q", response.Status, "ok")
				}
				return
			case "/phrases/random":
				var response randomPhrasesResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if response.Category != tt.wantCategory {
					t.Errorf("category body = %q, want %q", response.Category, tt.wantCategory)
				}
				if !reflect.DeepEqual(response.Phrases, tt.wantPhrases) {
					t.Errorf("phrases body = %#v, want %#v", response.Phrases, tt.wantPhrases)
				}
				return
			}

			var response categoriesResponse
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if !reflect.DeepEqual(response.Categories, tt.wantNames) {
				t.Errorf("categories body = %#v, want %#v", response.Categories, tt.wantNames)
			}
		})
	}
}

func TestRandomPhraseWithoutProgrammingCategory(t *testing.T) {
	handler := NewHandler([]phrase.Category{{Name: "study"}})
	request := httptest.NewRequest(http.MethodGet, "/phrases/random", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}

	var response errorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error != "failed to generate phrase" {
		t.Errorf("error body = %q, want %q", response.Error, "failed to generate phrase")
	}
}

func testCategories() []phrase.Category {
	return []phrase.Category{
		{
			Name:     "programming",
			Template: "{subject} {verb} {complement}",
			Parts: phrase.Parts{
				Subjects:    []string{"Codigo simples"},
				Verbs:       []string{"reduz"},
				Complements: []string{"problemas futuros"},
			},
		},
		{
			Name:     "study",
			Template: "{subject} {verb} {complement}",
			Parts: phrase.Parts{
				Subjects:    []string{"A pratica constante"},
				Verbs:       []string{"fortalece"},
				Complements: []string{"o aprendizado"},
			},
		},
	}
}
