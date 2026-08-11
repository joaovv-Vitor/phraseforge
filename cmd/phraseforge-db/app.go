package main

import (
	"context"
	"fmt"
	"io"

	"github.com/joaovv-Vitor/phraseforge/internal/migrate"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func run(ctx context.Context, out io.Writer, configuration config) error {
	count, err := SetupDatabase(ctx, configuration.databaseFile, configuration.dataFile)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "Database prepared successfully: imported %d categories.\n", count)
	return err
}

// SetupDatabase applies SQLite migrations and imports categories from a JSON file.
func SetupDatabase(ctx context.Context, databaseFile, dataFile string) (_ int, err error) {
	categories, err := storage.LoadCategories(dataFile)
	if err != nil {
		return 0, fmt.Errorf("load setup data: %w", err)
	}

	database, err := storage.OpenSQLite(ctx, databaseFile)
	if err != nil {
		return 0, fmt.Errorf("open setup database: %w", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close setup database: %w", closeErr)
		}
	}()

	if err := migrate.ApplySQLite(ctx, database); err != nil {
		return 0, fmt.Errorf("apply setup migrations: %w", err)
	}
	if err := storage.ImportCategories(ctx, database, categories); err != nil {
		return 0, fmt.Errorf("import setup categories: %w", err)
	}

	return len(categories), nil
}
