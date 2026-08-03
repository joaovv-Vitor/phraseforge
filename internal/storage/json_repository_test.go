package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents string
		wantName string
		wantErr  string
	}{
		{
			name: "loads valid categories",
			contents: `{
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
}`,
			wantName: "programming",
		},
		{
			name:     "returns error for invalid JSON",
			contents: `{"categories":`,
			wantErr:  "decode categories file",
		},
		{
			name:     "returns error for empty categories",
			contents: `{"categories": []}`,
			wantErr:  "categories cannot be empty",
		},
		{
			name: "returns error for category without name",
			contents: `{
  "categories": [{
    "template": "{subject} {verb} {complement}",
    "subjects": ["Codigo simples"],
    "verbs": ["reduz"],
    "complements": ["problemas futuros"]
  }]
}`,
			wantErr: "category name cannot be empty",
		},
		{
			name: "returns error for category without subjects",
			contents: `{
  "categories": [{
    "name": "programming",
    "template": "{subject} {verb} {complement}",
    "subjects": [],
    "verbs": ["reduz"],
    "complements": ["problemas futuros"]
  }]
}`,
			wantErr: "category \"programming\": subjects cannot be empty",
		},
		{
			name: "returns error for category without verbs",
			contents: `{
  "categories": [{
    "name": "programming",
    "template": "{subject} {verb} {complement}",
    "subjects": ["Codigo simples"],
    "verbs": [],
    "complements": ["problemas futuros"]
  }]
}`,
			wantErr: "category \"programming\": verbs cannot be empty",
		},
		{
			name: "returns error for category without complements",
			contents: `{
  "categories": [{
    "name": "programming",
    "template": "{subject} {verb} {complement}",
    "subjects": ["Codigo simples"],
    "verbs": ["reduz"],
    "complements": []
  }]
}`,
			wantErr: "category \"programming\": complements cannot be empty",
		},
		{
			name: "returns error for category without template",
			contents: `{
  "categories": [{
    "name": "programming",
    "subjects": ["Codigo simples"],
    "verbs": ["reduz"],
    "complements": ["problemas futuros"]
  }]
}`,
			wantErr: "template cannot be empty",
		},
		{
			name: "returns error for category with unknown template placeholder",
			contents: `{
  "categories": [{
    "name": "programming",
    "template": "{subject} {verb} {complement} {conclusion}",
    "subjects": ["Codigo simples"],
    "verbs": ["reduz"],
    "complements": ["problemas futuros"]
  }]
}`,
			wantErr: "unknown placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories, err := LoadCategories(writeTestFile(t, tt.contents))

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("LoadCategories() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("LoadCategories() error = %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("LoadCategories() unexpected error: %v", err)
			}
			if len(categories) != 2 {
				t.Fatalf("LoadCategories() returned %d categories, want 2", len(categories))
			}
			if categories[0].Name != tt.wantName {
				t.Errorf("LoadCategories() first category = %q, want %q", categories[0].Name, tt.wantName)
			}
		})
	}
}

func TestLoadCategoriesMissingFile(t *testing.T) {
	_, err := LoadCategories(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("LoadCategories() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "open categories file") {
		t.Fatalf("LoadCategories() error = %q, want an open error", err)
	}
}

func writeTestFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "phrases.json")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	return path
}
