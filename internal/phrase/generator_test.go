package phrase

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template string
		parts    Parts
		want     string
		wantErr  string
	}{
		{
			name:     "generates phrase with one option for each part",
			template: "{subject} {verb} {complement}",
			parts: Parts{
				Subjects:    []string{"A pratica constante"},
				Verbs:       []string{"transforma"},
				Complements: []string{"pequenos erros em grandes aprendizados"},
			},
			want: "A pratica constante transforma pequenos erros em grandes aprendizados.",
		},
		{
			name:     "generates phrase with reordered template",
			template: "{complement} e o que {subject} {verb}",
			parts: Parts{
				Subjects:    []string{"codigo simples"},
				Verbs:       []string{"reduz"},
				Complements: []string{"Problemas futuros"},
			},
			want: "Problemas futuros e o que codigo simples reduz.",
		},
		{
			name:     "returns error when subjects are empty",
			template: "{subject} {verb} {complement}",
			parts: Parts{
				Verbs:       []string{"transforma"},
				Complements: []string{"pequenos erros em grandes aprendizados"},
			},
			wantErr: "subjects cannot be empty",
		},
		{
			name:     "returns error when verbs are empty",
			template: "{subject} {verb} {complement}",
			parts: Parts{
				Subjects:    []string{"A pratica constante"},
				Complements: []string{"pequenos erros em grandes aprendizados"},
			},
			wantErr: "verbs cannot be empty",
		},
		{
			name:     "returns error when complements are empty",
			template: "{subject} {verb} {complement}",
			parts: Parts{
				Subjects: []string{"A pratica constante"},
				Verbs:    []string{"transforma"},
			},
			wantErr: "complements cannot be empty",
		},
		{
			name:    "returns error for empty template",
			parts:   validParts(),
			wantErr: "template cannot be empty",
		},
		{
			name:     "returns error for unknown placeholder",
			template: "{subject} {verb} {complement} {conclusion}",
			parts:    validParts(),
			wantErr:  "unknown placeholder",
		},
		{
			name:     "returns error for missing placeholder",
			template: "{subject} {verb}",
			parts:    validParts(),
			wantErr:  "{complement}",
		},
		{
			name:     "returns error for repeated placeholder",
			template: "{subject} {subject} {verb} {complement}",
			parts:    validParts(),
			wantErr:  "{subject}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.template, tt.parts)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("Generate() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Generate() error = %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Generate() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Generate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func validParts() Parts {
	return Parts{
		Subjects:    []string{"Codigo simples"},
		Verbs:       []string{"reduz"},
		Complements: []string{"problemas futuros"},
	}
}
