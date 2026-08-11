package storage_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/migrate"
	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestSQLitePhraseRepositoryLoadCategories(t *testing.T) {
	ctx := context.Background()
	database := openPhraseDatabase(t, ctx)
	want := []phrase.Category{
		{
			Name:     "programming",
			Template: "{introduction} {subject} {verb} {complement}{conclusion}",
			Parts: phrase.Parts{
				Introductions: []string{"Com foco,"},
				Subjects:      []string{"Codigo simples", "Um bom desenvolvedor"},
				Verbs:         []string{"reduz"},
				Complements:   []string{"problemas futuros"},
				Conclusions:   []string{", passo a passo"},
			},
		},
		phraseCategory("study", "A pratica constante"),
	}
	if err := storage.ImportCategories(ctx, database, []phrase.Category{want[1], want[0]}); err != nil {
		t.Fatalf("ImportCategories() unexpected error: %v", err)
	}

	categories, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err != nil {
		t.Fatalf("LoadCategories() unexpected error: %v", err)
	}
	if !reflect.DeepEqual(categories, want) {
		t.Errorf("LoadCategories() = %#v, want %#v", categories, want)
	}
}

func TestSQLitePhraseRepositoryLoadCategoriesEmpty(t *testing.T) {
	ctx := context.Background()
	database := openPhraseDatabase(t, ctx)

	categories, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err != nil {
		t.Fatalf("LoadCategories() unexpected error: %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("LoadCategories() returned %d categories, want 0", len(categories))
	}
}

func TestSQLitePhraseRepositoryLoadCategoriesRejectsMissingTemplate(t *testing.T) {
	ctx := context.Background()
	database := openPhraseDatabase(t, ctx)
	insertPhraseCategory(t, ctx, database, "programming")

	_, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err == nil {
		t.Fatal("LoadCategories() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "expected exactly one template, got 0") {
		t.Fatalf("LoadCategories() error = %q, want missing template error", err)
	}
}

func TestSQLitePhraseRepositoryLoadCategoriesRejectsMultipleTemplates(t *testing.T) {
	ctx := context.Background()
	database := openPhraseDatabase(t, ctx)
	categoryID := insertPhraseCategory(t, ctx, database, "programming")
	insertPhraseTemplate(t, ctx, database, categoryID, "{subject} {verb} {complement}")
	insertPhraseTemplate(t, ctx, database, categoryID, "{subject} {verb} {complement}")

	_, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err == nil {
		t.Fatal("LoadCategories() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "expected exactly one template, got 2") {
		t.Fatalf("LoadCategories() error = %q, want multiple template error", err)
	}
}

func TestSQLitePhraseRepositoryLoadCategoriesRejectsMissingSubjects(t *testing.T) {
	ctx := context.Background()
	database := openPhraseDatabase(t, ctx)
	categoryID := insertPhraseCategory(t, ctx, database, "programming")
	insertPhraseTemplate(t, ctx, database, categoryID, "{subject} {verb} {complement}")
	insertPhrasePart(t, ctx, database, categoryID, "verb", "reduz")
	insertPhrasePart(t, ctx, database, categoryID, "complement", "problemas futuros")

	_, err := storage.NewSQLitePhraseRepository(database).LoadCategories(ctx)
	if err == nil {
		t.Fatal("LoadCategories() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "subjects cannot be empty") {
		t.Fatalf("LoadCategories() error = %q, want validation error", err)
	}
}

func openPhraseDatabase(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	database, err := storage.OpenSQLite(ctx, filepath.Join(t.TempDir(), "phraseforge.db"))
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Errorf("close SQLite database: %v", err)
		}
	})

	if err := migrate.ApplySQLite(ctx, database); err != nil {
		t.Fatalf("apply SQLite migrations: %v", err)
	}

	return database
}

func phraseCategory(name, subject string) phrase.Category {
	return phrase.Category{
		Name:     name,
		Template: "{subject} {verb} {complement}",
		Parts: phrase.Parts{
			Subjects:    []string{subject},
			Verbs:       []string{"fortalece"},
			Complements: []string{"o aprendizado"},
		},
	}
}

func insertPhraseCategory(t *testing.T, ctx context.Context, database *sql.DB, name string) int64 {
	t.Helper()

	result, err := database.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		t.Fatalf("insert category %q: %v", name, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get category ID for %q: %v", name, err)
	}

	return id
}

func insertPhraseTemplate(t *testing.T, ctx context.Context, database *sql.DB, categoryID int64, content string) {
	t.Helper()

	if _, err := database.ExecContext(ctx, "INSERT INTO phrase_templates (category_id, content) VALUES (?, ?)", categoryID, content); err != nil {
		t.Fatalf("insert template: %v", err)
	}
}

func insertPhrasePart(t *testing.T, ctx context.Context, database *sql.DB, categoryID int64, kind, content string) {
	t.Helper()

	if _, err := database.ExecContext(ctx, "INSERT INTO phrase_parts (category_id, kind, content) VALUES (?, ?, ?)", categoryID, kind, content); err != nil {
		t.Fatalf("insert %s part: %v", kind, err)
	}
}
