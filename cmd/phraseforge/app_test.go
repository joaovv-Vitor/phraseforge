package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		want     string
		wantErr  string
		subjects []string
	}{
		{
			name: "generates a programming phrase",
			args: []string{"generate"},
		},
		{
			name:     "generates a study phrase",
			args:     []string{"generate", "--category", "study"},
			subjects: []string{"A pratica constante", "A revisao diaria"},
		},
		{
			name: "generates multiple numbered phrases",
			args: []string{"generate", "--count", "3"},
		},
		{
			name:     "generates multiple phrases for selected category",
			args:     []string{"generate", "--category", "study", "--count", "3"},
			subjects: []string{"A pratica constante", "A revisao diaria"},
		},
		{
			name: "lists categories in configured order",
			args: []string{"categories"},
			want: "programming\nstudy\n",
		},
		{
			name:    "returns error without command",
			wantErr: "command is required",
		},
		{
			name:    "returns error for unknown command",
			args:    []string{"invalid"},
			wantErr: "unknown command \"invalid\"",
		},
		{
			name:    "returns error for unknown category",
			args:    []string{"generate", "--category", "motivation"},
			wantErr: "category \"motivation\" not found",
		},
		{
			name:    "returns error for zero count",
			args:    []string{"generate", "--count", "0"},
			wantErr: "count must be greater than zero",
		},
		{
			name:    "returns error for negative count",
			args:    []string{"generate", "--count", "-1"},
			wantErr: "count must be greater than zero",
		},
		{
			name:    "returns error for unknown generate option",
			args:    []string{"generate", "--invalid"},
			wantErr: "flag provided but not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := run(tt.args, &output, testCategories())

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("run() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("run() error = %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("run() unexpected error: %v", err)
			}

			if strings.Contains(tt.name, "multiple") {
				assertNumberedPhrases(t, output.String(), 3)
				assertPhrasesUseSubjects(t, output.String(), tt.subjects)
				return
			}

			if tt.name == "generates a programming phrase" || tt.name == "generates a study phrase" {
				assertSinglePhrase(t, output.String())
				assertPhrasesUseSubjects(t, output.String(), tt.subjects)
				return
			}

			if got := output.String(); got != tt.want {
				t.Errorf("run() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func testCategories() []phrase.Category {
	return []phrase.Category{
		{
			Name: "programming",
			Parts: phrase.Parts{
				Subjects:    []string{"Codigo simples", "Um bom desenvolvedor"},
				Verbs:       []string{"reduz", "simplifica"},
				Complements: []string{"problemas futuros", "o trabalho da equipe"},
			},
		},
		{
			Name: "study",
			Parts: phrase.Parts{
				Subjects:    []string{"A pratica constante", "A revisao diaria"},
				Verbs:       []string{"fortalece", "melhora"},
				Complements: []string{"o aprendizado", "o raciocinio logico"},
			},
		},
	}
}

func assertSinglePhrase(t *testing.T, output string) {
	t.Helper()
	if output == "" || !strings.HasSuffix(output, ".\n") {
		t.Errorf("run() output = %q, want a non-empty phrase ending in a period", output)
	}
}

func assertNumberedPhrases(t *testing.T, output string, wantCount int) {
	t.Helper()

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != wantCount {
		t.Fatalf("run() produced %d lines, want %d", len(lines), wantCount)
	}

	for index, line := range lines {
		prefix := fmt.Sprintf("%d. ", index+1)
		if !strings.HasPrefix(line, prefix) || !strings.HasSuffix(line, ".") {
			t.Errorf("run() line %d = %q, want a numbered phrase", index+1, line)
		}
	}
}

func assertPhrasesUseSubjects(t *testing.T, output string, subjects []string) {
	t.Helper()
	if len(subjects) == 0 {
		return
	}

	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		found := false
		for _, subject := range subjects {
			if strings.Contains(line, subject) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("run() line = %q, want a phrase with one of %q", line, subjects)
		}
	}
}
