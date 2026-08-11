package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

// SQLitePhraseRepository loads complete phrase categories from SQLite.
type SQLitePhraseRepository struct {
	database *sql.DB
}

// NewSQLitePhraseRepository creates a repository backed by database.
func NewSQLitePhraseRepository(database *sql.DB) *SQLitePhraseRepository {
	return &SQLitePhraseRepository{database: database}
}

// LoadCategories returns complete, valid categories in alphabetical order.
func (repository *SQLitePhraseRepository) LoadCategories(ctx context.Context) ([]phrase.Category, error) {
	categories, byID, err := loadSQLiteCategories(ctx, repository.database)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return []phrase.Category{}, nil
	}

	if err := loadSQLiteTemplates(ctx, repository.database, byID); err != nil {
		return nil, err
	}
	if err := loadSQLiteParts(ctx, repository.database, byID); err != nil {
		return nil, err
	}

	result := make([]phrase.Category, 0, len(categories))
	for _, category := range categories {
		if category.templateCount != 1 {
			return nil, fmt.Errorf("load SQLite category %q: expected exactly one template, got %d", category.Name, category.templateCount)
		}
		if err := phrase.ValidateCategory(category.Category); err != nil {
			return nil, fmt.Errorf("validate SQLite category %q: %w", category.Name, err)
		}

		result = append(result, category.Category)
	}

	return result, nil
}

type sqlitePhraseCategory struct {
	phrase.Category
	id            int64
	templateCount int
}

func loadSQLiteCategories(ctx context.Context, database *sql.DB) (_ []*sqlitePhraseCategory, _ map[int64]*sqlitePhraseCategory, err error) {
	rows, err := database.QueryContext(ctx, `
SELECT id, name
FROM categories
ORDER BY name ASC`)
	if err != nil {
		return nil, nil, fmt.Errorf("list SQLite phrase categories: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close SQLite phrase category rows: %w", closeErr)
		}
	}()

	categories := make([]*sqlitePhraseCategory, 0)
	byID := make(map[int64]*sqlitePhraseCategory)
	for rows.Next() {
		category := &sqlitePhraseCategory{}
		if err := rows.Scan(&category.id, &category.Name); err != nil {
			return nil, nil, fmt.Errorf("scan SQLite phrase category: %w", err)
		}

		categories = append(categories, category)
		byID[category.id] = category
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate SQLite phrase categories: %w", err)
	}

	return categories, byID, nil
}

func loadSQLiteTemplates(ctx context.Context, database *sql.DB, categories map[int64]*sqlitePhraseCategory) (err error) {
	rows, err := database.QueryContext(ctx, `
SELECT category_id, content
FROM phrase_templates
ORDER BY id ASC`)
	if err != nil {
		return fmt.Errorf("list SQLite phrase templates: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close SQLite phrase template rows: %w", closeErr)
		}
	}()

	for rows.Next() {
		var categoryID int64
		var template string
		if err := rows.Scan(&categoryID, &template); err != nil {
			return fmt.Errorf("scan SQLite phrase template: %w", err)
		}

		category, found := categories[categoryID]
		if !found {
			return fmt.Errorf("load SQLite phrase template: category ID %d not found", categoryID)
		}
		category.Template = template
		category.templateCount++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate SQLite phrase templates: %w", err)
	}

	return nil
}

func loadSQLiteParts(ctx context.Context, database *sql.DB, categories map[int64]*sqlitePhraseCategory) (err error) {
	rows, err := database.QueryContext(ctx, `
SELECT category_id, kind, content
FROM phrase_parts
ORDER BY id ASC`)
	if err != nil {
		return fmt.Errorf("list SQLite phrase parts: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close SQLite phrase part rows: %w", closeErr)
		}
	}()

	for rows.Next() {
		var categoryID int64
		var kind string
		var content string
		if err := rows.Scan(&categoryID, &kind, &content); err != nil {
			return fmt.Errorf("scan SQLite phrase part: %w", err)
		}

		category, found := categories[categoryID]
		if !found {
			return fmt.Errorf("load SQLite phrase part: category ID %d not found", categoryID)
		}
		if err := appendSQLitePart(&category.Parts, kind, content); err != nil {
			return fmt.Errorf("load SQLite category %q: %w", category.Name, err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate SQLite phrase parts: %w", err)
	}

	return nil
}

func appendSQLitePart(parts *phrase.Parts, kind, content string) error {
	switch kind {
	case "introduction":
		parts.Introductions = append(parts.Introductions, content)
	case "subject":
		parts.Subjects = append(parts.Subjects, content)
	case "verb":
		parts.Verbs = append(parts.Verbs, content)
	case "complement":
		parts.Complements = append(parts.Complements, content)
	case "conclusion":
		parts.Conclusions = append(parts.Conclusions, content)
	default:
		return fmt.Errorf("unknown phrase part kind %q", kind)
	}

	return nil
}
