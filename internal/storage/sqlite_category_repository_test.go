package storage_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/migrate"
	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestSQLiteCategoryRepositoryList(t *testing.T) {
	ctx := context.Background()
	database := openMigratedDatabase(t, ctx)
	insertCategory(t, ctx, database, "study")
	insertCategory(t, ctx, database, "humor")
	insertCategory(t, ctx, database, "programming")

	repository := storage.NewSQLiteCategoryRepository(database)
	categories, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	want := []phrase.Category{{Name: "humor"}, {Name: "programming"}, {Name: "study"}}
	if !reflect.DeepEqual(categories, want) {
		t.Errorf("List() = %#v, want %#v", categories, want)
	}
}

func TestSQLiteCategoryRepositoryListEmpty(t *testing.T) {
	ctx := context.Background()
	repository := storage.NewSQLiteCategoryRepository(openMigratedDatabase(t, ctx))

	categories, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("List() returned %d categories, want 0", len(categories))
	}
}

func TestSQLiteCategoryRepositoryListReturnsErrorForClosedDatabase(t *testing.T) {
	ctx := context.Background()
	database := openMigratedDatabase(t, ctx)
	if err := database.Close(); err != nil {
		t.Fatalf("close SQLite database: %v", err)
	}

	_, err := storage.NewSQLiteCategoryRepository(database).List(ctx)
	if err == nil {
		t.Fatal("List() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "list SQLite categories") {
		t.Fatalf("List() error = %q, want list context", err)
	}
}

func openMigratedDatabase(t *testing.T, ctx context.Context) *sql.DB {
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

func insertCategory(t *testing.T, ctx context.Context, database *sql.DB, name string) {
	t.Helper()

	if _, err := database.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", name); err != nil {
		t.Fatalf("insert category %q: %v", name, err)
	}
}
