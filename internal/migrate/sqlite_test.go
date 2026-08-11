package migrate

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/db/migrations"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestApplySQLite(t *testing.T) {
	ctx := context.Background()
	database := openTestDatabase(t, ctx)

	if err := ApplySQLite(ctx, database); err != nil {
		t.Fatalf("ApplySQLite() unexpected error: %v", err)
	}

	for _, table := range []string{
		"categories",
		"phrase_templates",
		"phrase_parts",
		"favorite_phrases",
		"generation_history",
		"schema_migrations",
	} {
		if !tableExists(ctx, database, table) {
			t.Errorf("table %q does not exist", table)
		}
	}

	var version int64
	if err := database.QueryRowContext(ctx, "SELECT MAX(version) FROM schema_migrations").Scan(&version); err != nil {
		t.Fatalf("query latest migration version: %v", err)
	}
	if version != 2 {
		t.Errorf("latest migration version = %d, want 2", version)
	}
}

func TestApplySQLiteIsIdempotent(t *testing.T) {
	ctx := context.Background()
	database := openTestDatabase(t, ctx)

	if err := ApplySQLite(ctx, database); err != nil {
		t.Fatalf("first ApplySQLite() unexpected error: %v", err)
	}
	if err := ApplySQLite(ctx, database); err != nil {
		t.Fatalf("second ApplySQLite() unexpected error: %v", err)
	}

	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("count applied migrations: %v", err)
	}
	if count != 2 {
		t.Errorf("applied migrations = %d, want 2", count)
	}
}

func TestFavoritesAndHistoryConstraints(t *testing.T) {
	ctx := context.Background()
	database := openTestDatabase(t, ctx)

	if err := ApplySQLite(ctx, database); err != nil {
		t.Fatalf("ApplySQLite() unexpected error: %v", err)
	}

	categoryID := insertTestCategory(t, ctx, database)
	content := "Codigo simples reduz problemas futuros."

	var createdAt string
	if err := database.QueryRowContext(ctx, `
INSERT INTO favorite_phrases (category_id, content)
VALUES (?, ?)
RETURNING created_at`, categoryID, content).Scan(&createdAt); err != nil {
		t.Fatalf("insert favorite phrase: %v", err)
	}
	if createdAt == "" {
		t.Error("favorite created_at is empty")
	}

	_, err := database.ExecContext(ctx,
		"INSERT INTO favorite_phrases (category_id, content) VALUES (?, ?)",
		categoryID,
		content,
	)
	if err == nil {
		t.Fatal("duplicate favorite insert error = nil, want unique constraint error")
	}

	_, err = database.ExecContext(ctx,
		"INSERT INTO favorite_phrases (category_id, content) VALUES (?, ?)",
		categoryID+1,
		"Outra frase.",
	)
	if err == nil {
		t.Fatal("favorite with unknown category error = nil, want foreign key error")
	}

	_, err = database.ExecContext(ctx,
		"INSERT INTO favorite_phrases (category_id, content) VALUES (?, ?)",
		categoryID,
		"  ",
	)
	if err == nil {
		t.Fatal("favorite with empty content error = nil, want check constraint error")
	}

	var generatedAt string
	if err := database.QueryRowContext(ctx, `
INSERT INTO generation_history (category_id, content)
VALUES (?, ?)
RETURNING generated_at`, categoryID, content).Scan(&generatedAt); err != nil {
		t.Fatalf("insert history entry: %v", err)
	}
	if generatedAt == "" {
		t.Error("history generated_at is empty")
	}

	if _, err := database.ExecContext(ctx,
		"INSERT INTO generation_history (category_id, content) VALUES (?, ?)",
		categoryID,
		content,
	); err != nil {
		t.Fatalf("insert repeated history entry: %v", err)
	}

	var historyCount int
	if err := database.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM generation_history WHERE category_id = ? AND content = ?",
		categoryID,
		content,
	).Scan(&historyCount); err != nil {
		t.Fatalf("count history entries: %v", err)
	}
	if historyCount != 2 {
		t.Errorf("history entries = %d, want 2", historyCount)
	}

	_, err = database.ExecContext(ctx,
		"INSERT INTO generation_history (category_id, content) VALUES (?, ?)",
		categoryID+1,
		"Outra frase.",
	)
	if err == nil {
		t.Fatal("history with unknown category error = nil, want foreign key error")
	}

	_, err = database.ExecContext(ctx,
		"INSERT INTO generation_history (category_id, content) VALUES (?, ?)",
		categoryID,
		"  ",
	)
	if err == nil {
		t.Fatal("history with empty content error = nil, want check constraint error")
	}
}

func TestFavoritesAndHistoryDownMigration(t *testing.T) {
	ctx := context.Background()
	database := openTestDatabase(t, ctx)

	if err := ApplySQLite(ctx, database); err != nil {
		t.Fatalf("ApplySQLite() unexpected error: %v", err)
	}

	downMigration, err := migrations.Files.ReadFile("000002_favorites_history.down.sql")
	if err != nil {
		t.Fatalf("read down migration: %v", err)
	}
	if _, err := database.ExecContext(ctx, string(downMigration)); err != nil {
		t.Fatalf("apply down migration: %v", err)
	}

	for _, table := range []string{"favorite_phrases", "generation_history"} {
		if tableExists(ctx, database, table) {
			t.Errorf("table %q still exists after down migration", table)
		}
	}
	if !tableExists(ctx, database, "categories") {
		t.Error("categories table was removed by down migration")
	}
}

func insertTestCategory(t *testing.T, ctx context.Context, database *sql.DB) int64 {
	t.Helper()

	result, err := database.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", "programming")
	if err != nil {
		t.Fatalf("insert category: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get category ID: %v", err)
	}

	return id
}

func openTestDatabase(t *testing.T, ctx context.Context) *sql.DB {
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

	return database
}

func tableExists(ctx context.Context, database *sql.DB, name string) bool {
	var exists bool
	if err := database.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?)",
		name,
	).Scan(&exists); err != nil {
		return false
	}

	return exists
}
