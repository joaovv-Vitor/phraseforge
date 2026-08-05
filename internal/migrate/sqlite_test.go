package migrate

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestApplySQLite(t *testing.T) {
	ctx := context.Background()
	database := openTestDatabase(t, ctx)

	if err := ApplySQLite(ctx, database); err != nil {
		t.Fatalf("ApplySQLite() unexpected error: %v", err)
	}

	for _, table := range []string{"categories", "phrase_templates", "phrase_parts", "schema_migrations"} {
		if !tableExists(ctx, database, table) {
			t.Errorf("table %q does not exist", table)
		}
	}

	var version int
	if err := database.QueryRowContext(ctx, "SELECT version FROM schema_migrations").Scan(&version); err != nil {
		t.Fatalf("query migration version: %v", err)
	}
	if version != 1 {
		t.Errorf("migration version = %d, want 1", version)
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
	if count != 1 {
		t.Errorf("applied migrations = %d, want 1", count)
	}
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
