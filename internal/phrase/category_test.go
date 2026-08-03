package phrase

import (
	"strings"
	"testing"
)

func TestFindCategory(t *testing.T) {
	t.Parallel()

	categories := []Category{
		{
			Name: "programming",
			Parts: Parts{
				Subjects:    []string{"Codigo simples"},
				Verbs:       []string{"reduz"},
				Complements: []string{"problemas futuros"},
			},
		},
		{
			Name: "study",
			Parts: Parts{
				Subjects:    []string{"A revisao"},
				Verbs:       []string{"fortalece"},
				Complements: []string{"o aprendizado"},
			},
		},
	}

	tests := []struct {
		name    string
		input   []Category
		query   string
		want    string
		wantErr string
	}{
		{
			name:  "finds first category",
			input: categories,
			query: "programming",
			want:  "programming",
		},
		{
			name:  "finds category after first position",
			input: categories,
			query: "study",
			want:  "study",
		},
		{
			name:    "returns error for unknown category",
			input:   categories,
			query:   "motivation",
			wantErr: "category \"motivation\" not found",
		},
		{
			name:    "returns error for empty categories",
			input:   []Category{},
			query:   "programming",
			wantErr: "category \"programming\" not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindCategory(tt.input, tt.query)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("FindCategory() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("FindCategory() error = %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("FindCategory() unexpected error: %v", err)
			}
			if got.Name != tt.want {
				t.Errorf("FindCategory() name = %q, want %q", got.Name, tt.want)
			}
			if len(got.Parts.Subjects) == 0 {
				t.Error("FindCategory() returned category without subjects")
			}
		})
	}
}
