package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

// loadAPICategories loads the categories used by the HTTP API from SQLite.
func loadAPICategories(ctx context.Context, databaseFile string) (_ []phrase.Category, err error) {
	categories, database, err := openAPIData(ctx, databaseFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close API database: %w", closeErr)
		}
	}()

	return categories, nil
}

// openAPIData opens the SQLite database and loads the categories used by the HTTP API.
func openAPIData(ctx context.Context, databaseFile string) (_ []phrase.Category, _ *sql.DB, err error) {
	if _, err := os.Stat(databaseFile); err != nil {
		return nil, nil, fmt.Errorf("access API database %q: %w", databaseFile, err)
	}

	database, err := storage.OpenSQLite(ctx, databaseFile)
	if err != nil {
		return nil, nil, fmt.Errorf("open API database: %w", err)
	}
	defer func() {
		if err != nil {
			if closeErr := database.Close(); closeErr != nil {
				err = fmt.Errorf("%w; close API database: %v", err, closeErr)
			}
		}
	}()

	categories, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load API categories: %w", err)
	}
	if len(categories) == 0 {
		return nil, nil, fmt.Errorf("load API categories: database %q contains no categories", databaseFile)
	}

	return categories, database, nil
}
