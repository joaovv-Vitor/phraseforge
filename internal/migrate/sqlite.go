// Package migrate applies database schema migrations.
package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"

	"github.com/joaovv-Vitor/phraseforge/db/migrations"
)

type migration struct {
	version int64
	name    string
	sql     string
}

// ApplySQLite applies all embedded SQLite migrations that have not run yet.
func ApplySQLite(ctx context.Context, database *sql.DB) error {
	if _, err := database.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
)`); err != nil {
		return fmt.Errorf("create schema migrations table: %w", err)
	}

	migrations, err := sqliteMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := migrationApplied(ctx, database, migration.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		if err := applyMigration(ctx, database, migration); err != nil {
			return err
		}
	}

	return nil
}

func sqliteMigrations() ([]migration, error) {
	entries, err := fs.ReadDir(migrations.Files, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}

	filenames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			filenames = append(filenames, entry.Name())
		}
	}
	sort.Strings(filenames)

	result := make([]migration, 0, len(filenames))
	for _, filename := range filenames {
		version, err := migrationVersion(filename)
		if err != nil {
			return nil, err
		}

		sql, err := migrations.Files.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", filename, err)
		}
		result = append(result, migration{
			version: version,
			name:    filename,
			sql:     string(sql),
		})
	}

	return result, nil
}

func migrationVersion(filename string) (int64, error) {
	versionText, _, found := strings.Cut(filename, "_")
	if !found {
		return 0, fmt.Errorf("invalid migration filename %q", filename)
	}

	version, err := strconv.ParseInt(versionText, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse migration version from %q: %w", filename, err)
	}

	return version, nil
}

func migrationApplied(ctx context.Context, database *sql.DB, version int64) (bool, error) {
	var applied bool
	err := database.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)",
		version,
	).Scan(&applied)
	if err != nil {
		return false, fmt.Errorf("check migration %d: %w", version, err)
	}

	return applied, nil
}

func applyMigration(ctx context.Context, database *sql.DB, migration migration) error {
	transaction, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %d: %w", migration.version, err)
	}
	defer transaction.Rollback()

	if _, err := transaction.ExecContext(ctx, migration.sql); err != nil {
		return fmt.Errorf("execute migration %d: %w", migration.version, err)
	}
	if _, err := transaction.ExecContext(
		ctx,
		"INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
		migration.version,
		migration.name,
	); err != nil {
		return fmt.Errorf("record migration %d: %w", migration.version, err)
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit migration %d: %w", migration.version, err)
	}

	return nil
}
