package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

// ImportCategories stores categories, templates, and phrase parts in SQLite.
func ImportCategories(ctx context.Context, database *sql.DB, categories []phrase.Category) error {
	transaction, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin SQLite category import: %w", err)
	}
	defer transaction.Rollback()

	for _, category := range categories {
		if err := importCategory(ctx, transaction, category); err != nil {
			return fmt.Errorf("import category %q: %w", category.Name, err)
		}
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit SQLite category import: %w", err)
	}

	return nil
}

func importCategory(ctx context.Context, transaction *sql.Tx, category phrase.Category) error {
	result, err := transaction.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", category.Name)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get category ID: %w", err)
	}

	if _, err := transaction.ExecContext(
		ctx,
		"INSERT INTO phrase_templates (category_id, content) VALUES (?, ?)",
		categoryID,
		category.Template,
	); err != nil {
		return fmt.Errorf("insert phrase template: %w", err)
	}

	for _, part := range []struct {
		kind     string
		contents []string
	}{
		{kind: "introduction", contents: category.Introductions},
		{kind: "subject", contents: category.Subjects},
		{kind: "verb", contents: category.Verbs},
		{kind: "complement", contents: category.Complements},
		{kind: "conclusion", contents: category.Conclusions},
	} {
		if err := importParts(ctx, transaction, categoryID, part.kind, part.contents); err != nil {
			return err
		}
	}

	return nil
}

func importParts(ctx context.Context, transaction *sql.Tx, categoryID int64, kind string, contents []string) error {
	for _, content := range contents {
		if _, err := transaction.ExecContext(
			ctx,
			"INSERT INTO phrase_parts (category_id, kind, content) VALUES (?, ?, ?)",
			categoryID,
			kind,
			content,
		); err != nil {
			return fmt.Errorf("insert %s phrase part: %w", kind, err)
		}
	}

	return nil
}
