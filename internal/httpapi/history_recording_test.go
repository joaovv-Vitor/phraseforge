package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRandomPhraseRecordsHistory(t *testing.T) {
	store := &fakeHistoryStore{}
	handler := NewHandler(testCategories(), nil, store)
	request := httptest.NewRequest(http.MethodGet, "/phrases/random?category=programming&count=3", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if store.category != "programming" {
		t.Errorf("Record() category = %q, want %q", store.category, "programming")
	}
	want := []string{
		"Codigo simples reduz problemas futuros.",
		"Codigo simples reduz problemas futuros.",
		"Codigo simples reduz problemas futuros.",
	}
	if !reflect.DeepEqual(store.contents, want) {
		t.Errorf("Record() contents = %#v, want %#v", store.contents, want)
	}
}

func TestRandomPhraseReturnsErrorWhenHistoryRecordingFails(t *testing.T) {
	store := &fakeHistoryStore{recordErr: errors.New("database unavailable")}
	handler := NewHandler(testCategories(), nil, store)
	request := httptest.NewRequest(http.MethodGet, "/phrases/random", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	var response errorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error != "failed to record generation history" {
		t.Errorf("error body = %q, want %q", response.Error, "failed to record generation history")
	}
}

type noopHistoryStore struct{}

func (noopHistoryStore) Record(context.Context, string, []string) error {
	return nil
}

type fakeHistoryStore struct {
	category  string
	contents  []string
	recordErr error
}

func (store *fakeHistoryStore) Record(_ context.Context, category string, contents []string) error {
	if store.recordErr != nil {
		return store.recordErr
	}

	store.category = category
	store.contents = append([]string(nil), contents...)
	return nil
}
