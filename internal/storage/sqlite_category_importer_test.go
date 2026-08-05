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

func TestImportCategories(t *testing.T) {
	ctx := context.Background()
	database := openImportDatabase(t, ctx)
	categories := []phrase.Category{
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
		{
			Name:     "study",
			Template: "{subject} {verb} {complement}",
			Parts: phrase.Parts{
				Subjects:    []string{"A pratica constante"},
				Verbs:       []string{"fortalece"},
				Complements: []string{"o aprendizado"},
			},
		},
	}

	if err := storage.ImportCategories(ctx, database, categories); err != nil {
		t.Fatalf("ImportCategories() unexpected error: %v", err)
	}

	if got := countRows(t, ctx, database, "categories"); got != 2 {
		t.Errorf("categories count = %d, want 2", got)
	}
	if got := countRows(t, ctx, database, "phrase_templates"); got != 2 {
		t.Errorf("phrase templates count = %d, want 2", got)
	}
	if got := countRows(t, ctx, database, "phrase_parts"); got != 9 {
		t.Errorf("phrase parts count = %d, want 9", got)
	}

	assertPartCount(t, ctx, database, "programming", "introduction", 1)
	assertPartCount(t, ctx, database, "programming", "subject", 2)
	assertPartCount(t, ctx, database, "programming", "verb", 1)
	assertPartCount(t, ctx, database, "programming", "complement", 1)
	assertPartCount(t, ctx, database, "programming", "conclusion", 1)
	assertPartCount(t, ctx, database, "study", "introduction", 0)
	assertTemplate(t, ctx, database, "programming", "{introduction} {subject} {verb} {complement}{conclusion}")
	assertPartContents(t, ctx, database, "programming", "subject", []string{"Codigo simples", "Um bom desenvolvedor"})
	assertPartContents(t, ctx, database, "programming", "conclusion", []string{", passo a passo"})
}

func TestImportCategoriesRollsBackOnDuplicateCategory(t *testing.T) {
	ctx := context.Background()
	database := openImportDatabase(t, ctx)
	insertImportCategory(t, ctx, database, "programming")

	err := storage.ImportCategories(ctx, database, []phrase.Category{
		validImportCategory("study"),
		validImportCategory("programming"),
	})
	if err == nil {
		t.Fatal("ImportCategories() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "import category \"programming\"") {
		t.Fatalf("ImportCategories() error = %q, want category context", err)
	}

	if got := countRows(t, ctx, database, "categories"); got != 1 {
		t.Errorf("categories count after rollback = %d, want 1", got)
	}
	if got := countRows(t, ctx, database, "phrase_templates"); got != 0 {
		t.Errorf("phrase templates count after rollback = %d, want 0", got)
	}
	if got := countRows(t, ctx, database, "phrase_parts"); got != 0 {
		t.Errorf("phrase parts count after rollback = %d, want 0", got)
	}
}

func openImportDatabase(t *testing.T, ctx context.Context) *sql.DB {
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

func validImportCategory(name string) phrase.Category {
	return phrase.Category{
		Name:     name,
		Template: "{subject} {verb} {complement}",
		Parts: phrase.Parts{
			Subjects:    []string{"Codigo simples"},
			Verbs:       []string{"reduz"},
			Complements: []string{"problemas futuros"},
		},
	}
}

func insertImportCategory(t *testing.T, ctx context.Context, database *sql.DB, name string) {
	t.Helper()

	if _, err := database.ExecContext(ctx, "INSERT INTO categories (name) VALUES (?)", name); err != nil {
		t.Fatalf("insert category %q: %v", name, err)
	}
}

func countRows(t *testing.T, ctx context.Context, database *sql.DB, table string) int {
	t.Helper()

	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
		t.Fatalf("count rows in %s: %v", table, err)
	}

	return count
}

func assertPartCount(t *testing.T, ctx context.Context, database *sql.DB, categoryName, kind string, want int) {
	t.Helper()

	var got int
	if err := database.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM phrase_parts
JOIN categories ON categories.id = phrase_parts.category_id
WHERE categories.name = ? AND phrase_parts.kind = ?`, categoryName, kind).Scan(&got); err != nil {
		t.Fatalf("count %s parts for %q: %v", kind, categoryName, err)
	}
	if got != want {
		t.Errorf("%s parts for %q = %d, want %d", kind, categoryName, got, want)
	}
}

func assertTemplate(t *testing.T, ctx context.Context, database *sql.DB, categoryName, want string) {
	t.Helper()

	var got string
	if err := database.QueryRowContext(ctx, `
SELECT phrase_templates.content
FROM phrase_templates
JOIN categories ON categories.id = phrase_templates.category_id
WHERE categories.name = ?`, categoryName).Scan(&got); err != nil {
		t.Fatalf("query template for %q: %v", categoryName, err)
	}
	if got != want {
		t.Errorf("template for %q = %q, want %q", categoryName, got, want)
	}
}

func assertPartContents(t *testing.T, ctx context.Context, database *sql.DB, categoryName, kind string, want []string) {
	t.Helper()

	rows, err := database.QueryContext(ctx, `
SELECT phrase_parts.content
FROM phrase_parts
JOIN categories ON categories.id = phrase_parts.category_id
WHERE categories.name = ? AND phrase_parts.kind = ?
ORDER BY phrase_parts.id`, categoryName, kind)
	if err != nil {
		t.Fatalf("query %s parts for %q: %v", kind, categoryName, err)
	}
	defer rows.Close()

	got := make([]string, 0)
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			t.Fatalf("scan %s part for %q: %v", kind, categoryName, err)
		}
		got = append(got, content)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s parts for %q: %v", kind, categoryName, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s parts for %q = %#v, want %#v", kind, categoryName, got, want)
	}
}
