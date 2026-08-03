package main

import (
	"fmt"
	"log"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func main() {
	parts := phrase.Parts{
		Subjects:    []string{"A pratica constante"},
		Verbs:       []string{"transforma"},
		Complements: []string{"pequenos erros em grandes aprendizados"},
	}

	generated, err := phrase.Generate(parts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(generated)
}
