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
		name       string
		method     string
		path       string
		wantStatus int
		wantJSON   bool
		wantNames  []string
	}{
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
		},
		{
			name:       "returns not found for unknown route",
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
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
			if !tt.wantJSON {
				return
			}
			if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
				t.Fatalf("Content-Type = %q, want application/json", contentType)
			}

			if tt.path == "/health" {
				var response healthResponse
				if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if response.Status != "ok" {
					t.Errorf("status body = %q, want %q", response.Status, "ok")
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

func testCategories() []phrase.Category {
	return []phrase.Category{
		{Name: "programming"},
		{Name: "study"},
	}
}
