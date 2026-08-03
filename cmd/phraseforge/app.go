package main

import (
	"fmt"
	"io"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required: use generate or categories")
	}

	categories := defaultCategories()

	switch args[0] {
	case "generate":
		category, err := phrase.FindCategory(categories, "programming")
		if err != nil {
			return fmt.Errorf("run generate command: %w", err)
		}

		generated, err := phrase.Generate(category.Parts)
		if err != nil {
			return fmt.Errorf("run generate command: %w", err)
		}

		_, err = fmt.Fprintln(out, generated)
		return err
	case "categories":
		for _, name := range phrase.CategoryNames(categories) {
			if _, err := fmt.Fprintln(out, name); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown command %q: use generate or categories", args[0])
	}
}

func defaultCategories() []phrase.Category {
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
