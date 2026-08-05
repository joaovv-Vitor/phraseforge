package storage

import (
	"context"
	"path/filepath"
	"testing"
)

func TestOpenSQLite(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "phraseforge.db")

	database, err := OpenSQLite(ctx, path)
	if err != nil {
		t.Fatalf("OpenSQLite() unexpected error: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Errorf("close SQLite database: %v", err)
		}
	})

	if err := database.PingContext(ctx); err != nil {
		t.Fatalf("PingContext() unexpected error: %v", err)
	}

	var foreignKeys int
	if err := database.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("query foreign_keys pragma: %v", err)
	}
	if foreignKeys != 1 {
		t.Errorf("foreign_keys = %d, want 1", foreignKeys)
	}
}

func TestOpenSQLiteReturnsErrorForInvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "phraseforge.db")

	database, err := OpenSQLite(context.Background(), path)
	if err == nil {
		database.Close()
		t.Fatal("OpenSQLite() error = nil, want an error")
	}
}
