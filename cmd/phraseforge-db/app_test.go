package main

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func TestRunPreparesDatabase(t *testing.T) {
	ctx := context.Background()
	directory := t.TempDir()
	dataFile := writeCategoriesFile(t, directory)
	databaseFile := filepath.Join(directory, "phraseforge.db")
	var output bytes.Buffer

	err := run(ctx, &output, config{
		dataFile:     dataFile,
		databaseFile: databaseFile,
	})
	if err != nil {
		t.Fatalf("run() unexpected error: %v", err)
	}
	if got, want := output.String(), "Database prepared successfully: imported 2 categories.\n"; got != want {
		t.Errorf("run() output = %q, want %q", got, want)
	}

	database, err := storage.OpenSQLite(ctx, databaseFile)
	if err != nil {
		t.Fatalf("open prepared database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	if got := databaseRowCount(t, ctx, database, "categories"); got != 2 {
		t.Errorf("categories count = %d, want 2", got)
	}
	if got := databaseRowCount(t, ctx, database, "phrase_templates"); got != 2 {
		t.Errorf("phrase templates count = %d, want 2", got)
	}
	if got := databaseRowCount(t, ctx, database, "phrase_parts"); got != 6 {
		t.Errorf("phrase parts count = %d, want 6", got)
	}
}

func TestSetupDatabaseDoesNotCreateDatabaseForInvalidData(t *testing.T) {
	directory := t.TempDir()
	dataFile := filepath.Join(directory, "invalid.json")
	if err := os.WriteFile(dataFile, []byte(`{"categories": []}`), 0o600); err != nil {
		t.Fatalf("write invalid data file: %v", err)
	}
	databaseFile := filepath.Join(directory, "phraseforge.db")

	_, err := SetupDatabase(context.Background(), databaseFile, dataFile)
	if err == nil {
		t.Fatal("SetupDatabase() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "categories cannot be empty") {
		t.Fatalf("SetupDatabase() error = %q, want validation context", err)
	}
	if _, err := os.Stat(databaseFile); !os.IsNotExist(err) {
		t.Errorf("database file exists after invalid data, want it not to exist")
	}
}

func writeCategoriesFile(t *testing.T, directory string) string {
	t.Helper()

	path := filepath.Join(directory, "phrases.json")
	contents := `{
  "categories": [
    {
      "name": "programming",
      "template": "{subject} {verb} {complement}",
      "subjects": ["Codigo simples"],
      "verbs": ["reduz"],
      "complements": ["problemas futuros"]
    },
    {
      "name": "study",
      "template": "{subject} {verb} {complement}",
      "subjects": ["A pratica constante"],
      "verbs": ["fortalece"],
      "complements": ["o aprendizado"]
    }
  ]
}`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write categories file: %v", err)
	}

	return path
}

func databaseRowCount(t *testing.T, ctx context.Context, database *sql.DB, table string) int {
	t.Helper()

	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
		t.Fatalf("count rows in %s: %v", table, err)
	}

	return count
}
