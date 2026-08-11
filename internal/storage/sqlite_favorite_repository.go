package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

// SQLiteFavoriteRepository stores favorite phrases in SQLite.
type SQLiteFavoriteRepository struct {
	database *sql.DB
}

// NewSQLiteFavoriteRepository creates a repository backed by database.
func NewSQLiteFavoriteRepository(database *sql.DB) *SQLiteFavoriteRepository {
	return &SQLiteFavoriteRepository{database: database}
}

// Create saves a favorite phrase for the named category.
func (repository *SQLiteFavoriteRepository) Create(ctx context.Context, categoryName, content string) (phrase.Favorite, error) {
	categoryName = strings.TrimSpace(categoryName)
	content = strings.TrimSpace(content)
	if categoryName == "" {
		return phrase.Favorite{}, fmt.Errorf("create SQLite favorite: category name cannot be empty")
	}
	if content == "" {
		return phrase.Favorite{}, fmt.Errorf("create SQLite favorite: content cannot be empty")
	}

	categoryID, err := repository.categoryID(ctx, categoryName)
	if err != nil {
		return phrase.Favorite{}, err
	}

	var favorite phrase.Favorite
	var createdAt string
	err = repository.database.QueryRowContext(ctx, `
INSERT INTO favorite_phrases (category_id, content)
VALUES (?, ?)
ON CONFLICT (category_id, content) DO NOTHING
RETURNING id, created_at`, categoryID, content).Scan(&favorite.ID, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return phrase.Favorite{}, fmt.Errorf("create SQLite favorite for category %q: %w", categoryName, phrase.ErrFavoriteAlreadyExists)
	}
	if err != nil {
		return phrase.Favorite{}, fmt.Errorf("insert SQLite favorite for category %q: %w", categoryName, err)
	}

	favorite.Category = categoryName
	favorite.Content = content
	favorite.CreatedAt, err = parseSQLiteTimestamp(createdAt)
	if err != nil {
		return phrase.Favorite{}, fmt.Errorf("parse SQLite favorite creation time: %w", err)
	}

	return favorite, nil
}

// List returns favorites from most recently created to oldest.
func (repository *SQLiteFavoriteRepository) List(ctx context.Context) (_ []phrase.Favorite, err error) {
	rows, err := repository.database.QueryContext(ctx, `
SELECT favorite_phrases.id, categories.name, favorite_phrases.content, favorite_phrases.created_at
FROM favorite_phrases
JOIN categories ON categories.id = favorite_phrases.category_id
ORDER BY favorite_phrases.created_at DESC, favorite_phrases.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list SQLite favorites: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close SQLite favorite rows: %w", closeErr)
		}
	}()

	favorites := make([]phrase.Favorite, 0)
	for rows.Next() {
		var favorite phrase.Favorite
		var createdAt string
		if err := rows.Scan(&favorite.ID, &favorite.Category, &favorite.Content, &createdAt); err != nil {
			return nil, fmt.Errorf("scan SQLite favorite: %w", err)
		}

		favorite.CreatedAt, err = parseSQLiteTimestamp(createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse SQLite favorite creation time: %w", err)
		}
		favorites = append(favorites, favorite)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate SQLite favorites: %w", err)
	}

	return favorites, nil
}

func (repository *SQLiteFavoriteRepository) categoryID(ctx context.Context, name string) (int64, error) {
	var id int64
	err := repository.database.QueryRowContext(ctx, "SELECT id FROM categories WHERE name = ?", name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("create SQLite favorite for category %q: %w", name, phrase.ErrFavoriteCategoryNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("find SQLite favorite category %q: %w", name, err)
	}

	return id, nil
}
