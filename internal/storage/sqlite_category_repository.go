package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

// SQLiteCategoryRepository reads categories stored in SQLite.
type SQLiteCategoryRepository struct {
	database *sql.DB
}

// NewSQLiteCategoryRepository creates a repository backed by database.
func NewSQLiteCategoryRepository(database *sql.DB) *SQLiteCategoryRepository {
	return &SQLiteCategoryRepository{database: database}
}

// List returns category names in alphabetical order.
func (repository *SQLiteCategoryRepository) List(ctx context.Context) (_ []phrase.Category, err error) {
	rows, err := repository.database.QueryContext(ctx, `
SELECT name
FROM categories
ORDER BY name ASC`)
	if err != nil {
		return nil, fmt.Errorf("list SQLite categories: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close SQLite category rows: %w", closeErr)
		}
	}()

	categories := make([]phrase.Category, 0)
	for rows.Next() {
		var category phrase.Category
		if err := rows.Scan(&category.Name); err != nil {
			return nil, fmt.Errorf("scan SQLite category: %w", err)
		}

		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate SQLite categories: %w", err)
	}

	return categories, nil
}
