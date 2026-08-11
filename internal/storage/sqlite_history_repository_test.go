package storage_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/migrate"
	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestSQLiteHistoryRepositoryCreate(t *testing.T) {
	ctx := context.Background()
	database := openHistoryDatabase(t, ctx)
	insertHistoryCategory(t, ctx, database, "programming")
	repository := storage.NewSQLiteHistoryRepository(database)

	entry, err := repository.Create(ctx, " programming ", " Codigo simples reduz problemas futuros. ")
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if entry.ID == 0 {
		t.Error("Create() returned an empty ID")
	}
	if entry.Category != "programming" {
		t.Errorf("Create() category = %q, want %q", entry.Category, "programming")
	}
	if entry.Content != "Codigo simples reduz problemas futuros." {
		t.Errorf("Create() content = %q", entry.Content)
	}
	if entry.GeneratedAt.IsZero() {
		t.Error("Create() returned an empty generation time")
	}
}

func TestSQLiteHistoryRepositoryCreateRejectsInvalidInput(t *testing.T) {
	ctx := context.Background()
	database := openHistoryDatabase(t, ctx)
	insertHistoryCategory(t, ctx, database, "programming")
	repository := storage.NewSQLiteHistoryRepository(database)

	tests := []struct {
		name     string
		category string
		content  string
		wantErr  string
	}{
		{
			name:    "empty category",
			content: "Codigo simples reduz problemas futuros.",
			wantErr: "category name cannot be empty",
		},
		{
			name:     "empty content",
			category: "programming",
			content:  "   ",
			wantErr:  "content cannot be empty",
		},
		{
			name:     "unknown category",
			category: "study",
			content:  "A pratica constante fortalece o aprendizado.",
			wantErr:  "history category not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repository.Create(ctx, tt.category, tt.content)
			if err == nil {
				t.Fatal("Create() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Create() error = %q, want %q", err, tt.wantErr)
			}
			if tt.name == "unknown category" && !errors.Is(err, phrase.ErrHistoryCategoryNotFound) {
				t.Errorf("Create() error = %v, want ErrHistoryCategoryNotFound", err)
			}
		})
	}
}

func TestSQLiteHistoryRepositoryAllowsRepeatedEntries(t *testing.T) {
	ctx := context.Background()
	database := openHistoryDatabase(t, ctx)
	insertHistoryCategory(t, ctx, database, "programming")
	repository := storage.NewSQLiteHistoryRepository(database)

	first, err := repository.Create(ctx, "programming", "Codigo simples reduz problemas futuros.")
	if err != nil {
		t.Fatalf("first Create() unexpected error: %v", err)
	}
	second, err := repository.Create(ctx, "programming", "Codigo simples reduz problemas futuros.")
	if err != nil {
		t.Fatalf("second Create() unexpected error: %v", err)
	}
	if second.ID <= first.ID {
		t.Errorf("second ID = %d, want greater than first ID %d", second.ID, first.ID)
	}
}

func TestSQLiteHistoryRepositoryList(t *testing.T) {
	ctx := context.Background()
	database := openHistoryDatabase(t, ctx)
	insertHistoryCategory(t, ctx, database, "programming")
	insertHistoryCategory(t, ctx, database, "study")
	repository := storage.NewSQLiteHistoryRepository(database)

	first, err := repository.Create(ctx, "programming", "Codigo simples reduz problemas futuros.")
	if err != nil {
		t.Fatalf("create first history entry: %v", err)
	}
	second, err := repository.Create(ctx, "study", "A pratica constante fortalece o aprendizado.")
	if err != nil {
		t.Fatalf("create second history entry: %v", err)
	}

	entries, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
	want := []phrase.HistoryEntry{second, first}
	if !reflect.DeepEqual(entries, want) {
		t.Errorf("List() = %#v, want %#v", entries, want)
	}
}

func TestSQLiteHistoryRepositoryListEmpty(t *testing.T) {
	ctx := context.Background()
	repository := storage.NewSQLiteHistoryRepository(openHistoryDatabase(t, ctx))

	entries, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("List() returned %d entries, want 0", len(entries))
	}
}

func TestSQLiteHistoryRepositoryListReturnsErrorForClosedDatabase(t *testing.T) {
	ctx := context.Background()
	database := openHistoryDatabase(t, ctx)
	if err := database.Close(); err != nil {
		t.Fatalf("close SQLite database: %v", err)
	}

	_, err := storage.NewSQLiteHistoryRepository(database).List(ctx)
	if err == nil {
		t.Fatal("List() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "list SQLite history entries") {
		t.Errorf("List() error = %q, want list context", err)
	}
}

func openHistoryDatabase(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	database, err := storage.OpenSQLite(ctx, filepath.Join(t.TempDir(), "phraseforge.db"))
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Errorf("close SQLite database: %v", err)
		}
	})

	if err := migrate.ApplySQLite(ctx, database); err != nil {
		t.Fatalf("apply SQLite migrations: %v", err)
	}

	return database
}

func insertHistoryCategory(t *testing.T, ctx context.Context, database *sql.DB, name string) {
	t.Helper()

	if _, err := database.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", name); err != nil {
		t.Fatalf("insert category %q: %v", name, err)
	}
}
