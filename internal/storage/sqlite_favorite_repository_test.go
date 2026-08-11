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

func TestSQLiteFavoriteRepositoryCreate(t *testing.T) {
	ctx := context.Background()
	database := openFavoriteDatabase(t, ctx)
	insertFavoriteCategory(t, ctx, database, "programming")
	repository := storage.NewSQLiteFavoriteRepository(database)

	favorite, err := repository.Create(ctx, " programming ", " Codigo simples reduz problemas futuros. ")
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if favorite.ID == 0 {
		t.Error("Create() returned an empty ID")
	}
	if favorite.Category != "programming" {
		t.Errorf("Create() category = %q, want %q", favorite.Category, "programming")
	}
	if favorite.Content != "Codigo simples reduz problemas futuros." {
		t.Errorf("Create() content = %q", favorite.Content)
	}
	if favorite.CreatedAt.IsZero() {
		t.Error("Create() returned an empty creation time")
	}
}

func TestSQLiteFavoriteRepositoryCreateRejectsInvalidInput(t *testing.T) {
	ctx := context.Background()
	database := openFavoriteDatabase(t, ctx)
	insertFavoriteCategory(t, ctx, database, "programming")
	repository := storage.NewSQLiteFavoriteRepository(database)

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
			wantErr:  "favorite category not found",
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
			if tt.name == "unknown category" && !errors.Is(err, phrase.ErrFavoriteCategoryNotFound) {
				t.Errorf("Create() error = %v, want ErrFavoriteCategoryNotFound", err)
			}
		})
	}
}

func TestSQLiteFavoriteRepositoryCreateRejectsDuplicate(t *testing.T) {
	ctx := context.Background()
	database := openFavoriteDatabase(t, ctx)
	insertFavoriteCategory(t, ctx, database, "programming")
	repository := storage.NewSQLiteFavoriteRepository(database)

	if _, err := repository.Create(ctx, "programming", "Codigo simples reduz problemas futuros."); err != nil {
		t.Fatalf("first Create() unexpected error: %v", err)
	}
	_, err := repository.Create(ctx, "programming", "Codigo simples reduz problemas futuros.")
	if err == nil {
		t.Fatal("second Create() error = nil, want duplicate error")
	}
	if !errors.Is(err, phrase.ErrFavoriteAlreadyExists) {
		t.Errorf("second Create() error = %v, want ErrFavoriteAlreadyExists", err)
	}
}

func TestSQLiteFavoriteRepositoryList(t *testing.T) {
	ctx := context.Background()
	database := openFavoriteDatabase(t, ctx)
	insertFavoriteCategory(t, ctx, database, "programming")
	insertFavoriteCategory(t, ctx, database, "study")
	repository := storage.NewSQLiteFavoriteRepository(database)

	first, err := repository.Create(ctx, "programming", "Codigo simples reduz problemas futuros.")
	if err != nil {
		t.Fatalf("create first favorite: %v", err)
	}
	second, err := repository.Create(ctx, "study", "A pratica constante fortalece o aprendizado.")
	if err != nil {
		t.Fatalf("create second favorite: %v", err)
	}

	favorites, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
	want := []phrase.Favorite{second, first}
	if !reflect.DeepEqual(favorites, want) {
		t.Errorf("List() = %#v, want %#v", favorites, want)
	}
}

func TestSQLiteFavoriteRepositoryListEmpty(t *testing.T) {
	ctx := context.Background()
	repository := storage.NewSQLiteFavoriteRepository(openFavoriteDatabase(t, ctx))

	favorites, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
	if len(favorites) != 0 {
		t.Errorf("List() returned %d favorites, want 0", len(favorites))
	}
}

func TestSQLiteFavoriteRepositoryListReturnsErrorForClosedDatabase(t *testing.T) {
	ctx := context.Background()
	database := openFavoriteDatabase(t, ctx)
	if err := database.Close(); err != nil {
		t.Fatalf("close SQLite database: %v", err)
	}

	_, err := storage.NewSQLiteFavoriteRepository(database).List(ctx)
	if err == nil {
		t.Fatal("List() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "list SQLite favorites") {
		t.Errorf("List() error = %q, want list context", err)
	}
}

func openFavoriteDatabase(t *testing.T, ctx context.Context) *sql.DB {
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

func insertFavoriteCategory(t *testing.T, ctx context.Context, database *sql.DB, name string) {
	t.Helper()

	if _, err := database.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", name); err != nil {
		t.Fatalf("insert category %q: %v", name, err)
	}
}
