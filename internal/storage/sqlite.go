package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// OpenSQLite opens a SQLite database and configures connection-level settings.
func OpenSQLite(ctx context.Context, path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open SQLite database %q: %w", path, err)
	}

	database.SetMaxOpenConns(1)

	if err := database.PingContext(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping SQLite database %q: %w", path, err)
	}
	if _, err := database.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		database.Close()
		return nil, fmt.Errorf("enable SQLite foreign keys for %q: %w", path, err)
	}

	return database, nil
}
