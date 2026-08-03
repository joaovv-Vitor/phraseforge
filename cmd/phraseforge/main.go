package main

import (
	"fmt"
	"os"

	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func main() {
	categories, err := storage.LoadCategories("data/phrases.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(os.Args[1:], os.Stdout, categories); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
