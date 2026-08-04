package main

import (
	"fmt"
	"os"

	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func main() {
	config, args, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	categories, err := storage.LoadCategories(config.dataFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(args, os.Stdout, categories); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
