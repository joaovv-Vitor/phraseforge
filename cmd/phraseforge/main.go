package main

import (
	"fmt"
	"log"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func main() {
	categories := []phrase.Category{
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

	category, err := phrase.FindCategory(categories, "programming")
	if err != nil {
		log.Fatal(err)
	}

	generated, err := phrase.Generate(category.Parts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(generated)
}
