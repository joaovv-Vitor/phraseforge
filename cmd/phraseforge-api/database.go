package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

// loadAPICategories loads the categories used by the HTTP API from SQLite.
func loadAPICategories(ctx context.Context, databaseFile string) (_ []phrase.Category, err error) {
	if _, err := os.Stat(databaseFile); err != nil {
		return nil, fmt.Errorf("access API database %q: %w", databaseFile, err)
	}

	database, err := storage.OpenSQLite(ctx, databaseFile)
	if err != nil {
		return nil, fmt.Errorf("open API database: %w", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close API database: %w", closeErr)
		}
	}()

	categories, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("load API categories: %w", err)
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("load API categories: database %q contains no categories", databaseFile)
	}

	return categories, nil
}
