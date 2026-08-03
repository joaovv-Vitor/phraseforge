package phrase

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		parts   Parts
		want    string
		wantErr string
	}{
		{
			name: "generates phrase with one option for each part",
			parts: Parts{
				Subjects:    []string{"A pratica constante"},
				Verbs:       []string{"transforma"},
				Complements: []string{"pequenos erros em grandes aprendizados"},
			},
			want: "A pratica constante transforma pequenos erros em grandes aprendizados.",
		},
		{
			name: "returns error when subjects are empty",
			parts: Parts{
				Verbs:       []string{"transforma"},
				Complements: []string{"pequenos erros em grandes aprendizados"},
			},
			wantErr: "subjects cannot be empty",
		},
		{
			name: "returns error when verbs are empty",
			parts: Parts{
				Subjects:    []string{"A pratica constante"},
				Complements: []string{"pequenos erros em grandes aprendizados"},
			},
			wantErr: "verbs cannot be empty",
		},
		{
			name: "returns error when complements are empty",
			parts: Parts{
				Subjects: []string{"A pratica constante"},
				Verbs:    []string{"transforma"},
			},
			wantErr: "complements cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.parts)

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
