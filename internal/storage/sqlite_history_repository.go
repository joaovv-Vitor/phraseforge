package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

// SQLiteHistoryRepository stores generated phrase history in SQLite.
type SQLiteHistoryRepository struct {
	database *sql.DB
}

// NewSQLiteHistoryRepository creates a repository backed by database.
func NewSQLiteHistoryRepository(database *sql.DB) *SQLiteHistoryRepository {
	return &SQLiteHistoryRepository{database: database}
}

// Create records a generated phrase for the named category.
func (repository *SQLiteHistoryRepository) Create(ctx context.Context, categoryName, content string) (phrase.HistoryEntry, error) {
	categoryName = strings.TrimSpace(categoryName)
	content = strings.TrimSpace(content)
	if categoryName == "" {
		return phrase.HistoryEntry{}, fmt.Errorf("create SQLite history entry: category name cannot be empty")
	}
	if content == "" {
		return phrase.HistoryEntry{}, fmt.Errorf("create SQLite history entry: content cannot be empty")
	}

	categoryID, err := repository.categoryID(ctx, categoryName)
	if err != nil {
		return phrase.HistoryEntry{}, err
	}

	var entry phrase.HistoryEntry
	var generatedAt string
	err = repository.database.QueryRowContext(ctx, `
INSERT INTO generation_history (category_id, content)
VALUES (?, ?)
RETURNING id, generated_at`, categoryID, content).Scan(&entry.ID, &generatedAt)
	if err != nil {
		return phrase.HistoryEntry{}, fmt.Errorf("insert SQLite history entry for category %q: %w", categoryName, err)
	}

	entry.Category = categoryName
	entry.Content = content
	entry.GeneratedAt, err = parseSQLiteTimestamp(generatedAt)
	if err != nil {
		return phrase.HistoryEntry{}, fmt.Errorf("parse SQLite history generation time: %w", err)
	}

	return entry, nil
}

// Record stores all generated phrases in a single transaction.
func (repository *SQLiteHistoryRepository) Record(ctx context.Context, categoryName string, contents []string) error {
	categoryName = strings.TrimSpace(categoryName)
	if categoryName == "" {
		return fmt.Errorf("record SQLite history entries: category name cannot be empty")
	}
	if len(contents) == 0 {
		return fmt.Errorf("record SQLite history entries: contents cannot be empty")
	}

	normalizedContents := make([]string, 0, len(contents))
	for index, content := range contents {
		content = strings.TrimSpace(content)
		if content == "" {
			return fmt.Errorf("record SQLite history entries: content %d cannot be empty", index+1)
		}
		normalizedContents = append(normalizedContents, content)
	}

	transaction, err := repository.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin SQLite history recording: %w", err)
	}
	defer transaction.Rollback()

	var categoryID int64
	err = transaction.QueryRowContext(ctx, "SELECT id FROM categories WHERE name = ?", categoryName).Scan(&categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("record SQLite history entries for category %q: %w", categoryName, phrase.ErrHistoryCategoryNotFound)
	}
	if err != nil {
		return fmt.Errorf("find SQLite history category %q: %w", categoryName, err)
	}

	for _, content := range normalizedContents {
		if _, err := transaction.ExecContext(
			ctx,
			"INSERT INTO generation_history (category_id, content) VALUES (?, ?)",
			categoryID,
			content,
		); err != nil {
			return fmt.Errorf("insert SQLite history entry for category %q: %w", categoryName, err)
		}
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit SQLite history recording: %w", err)
	}

	return nil
}

// List returns history entries from most recent to oldest.
func (repository *SQLiteHistoryRepository) List(ctx context.Context) (_ []phrase.HistoryEntry, err error) {
	rows, err := repository.database.QueryContext(ctx, `
SELECT generation_history.id, categories.name, generation_history.content, generation_history.generated_at
FROM generation_history
JOIN categories ON categories.id = generation_history.category_id
ORDER BY generation_history.generated_at DESC, generation_history.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list SQLite history entries: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close SQLite history rows: %w", closeErr)
		}
	}()

	entries := make([]phrase.HistoryEntry, 0)
	for rows.Next() {
		var entry phrase.HistoryEntry
		var generatedAt string
		if err := rows.Scan(&entry.ID, &entry.Category, &entry.Content, &generatedAt); err != nil {
			return nil, fmt.Errorf("scan SQLite history entry: %w", err)
		}

		entry.GeneratedAt, err = parseSQLiteTimestamp(generatedAt)
		if err != nil {
			return nil, fmt.Errorf("parse SQLite history generation time: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate SQLite history entries: %w", err)
	}

	return entries, nil
}

func (repository *SQLiteHistoryRepository) categoryID(ctx context.Context, name string) (int64, error) {
	var id int64
	err := repository.database.QueryRowContext(ctx, "SELECT id FROM categories WHERE name = ?", name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("create SQLite history entry for category %q: %w", name, phrase.ErrHistoryCategoryNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("find SQLite history category %q: %w", name, err)
	}

	return id, nil
}
