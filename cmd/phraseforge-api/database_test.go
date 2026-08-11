package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/migrate"
	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestLoadAPICategories(t *testing.T) {
	ctx := context.Background()
	databaseFile := filepath.Join(t.TempDir(), "phraseforge.db")
	want := []phrase.Category{{
		Name:     "programming",
		Template: "{subject} {verb} {complement}",
		Parts: phrase.Parts{
			Subjects:    []string{"Codigo simples"},
			Verbs:       []string{"reduz"},
			Complements: []string{"problemas futuros"},
		},
	}}
	prepareAPIDatabase(t, ctx, databaseFile, want)

	got, err := loadAPICategories(ctx, databaseFile)
	if err != nil {
		t.Fatalf("loadAPICategories() unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("loadAPICategories() = %#v, want %#v", got, want)
	}
}

func TestLoadAPICategoriesMissingDatabase(t *testing.T) {
	databaseFile := filepath.Join(t.TempDir(), "missing.db")

	_, err := loadAPICategories(context.Background(), databaseFile)
	if err == nil {
		t.Fatal("loadAPICategories() error = nil, want missing database error")
	}
	if !strings.Contains(err.Error(), "access API database") {
		t.Errorf("loadAPICategories() error = %q, want database access context", err)
	}
	if _, statErr := os.Stat(databaseFile); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("database file was created, stat error = %v, want not exist", statErr)
	}
}

func TestLoadAPICategoriesEmptyDatabase(t *testing.T) {
	ctx := context.Background()
	databaseFile := filepath.Join(t.TempDir(), "phraseforge.db")
	prepareAPIDatabase(t, ctx, databaseFile, nil)

	_, err := loadAPICategories(ctx, databaseFile)
	if err == nil {
		t.Fatal("loadAPICategories() error = nil, want empty database error")
	}
	if !strings.Contains(err.Error(), "contains no categories") {
		t.Errorf("loadAPICategories() error = %q, want empty database context", err)
	}
}

func prepareAPIDatabase(t *testing.T, ctx context.Context, databaseFile string, categories []phrase.Category) {
	t.Helper()

	database, err := storage.OpenSQLite(ctx, databaseFile)
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			t.Errorf("close SQLite database: %v", err)
		}
	}()

	if err := migrate.ApplySQLite(ctx, database); err != nil {
		t.Fatalf("apply SQLite migrations: %v", err)
	}
	if err := storage.ImportCategories(ctx, database, categories); err != nil {
		t.Fatalf("import SQLite categories: %v", err)
	}
}
