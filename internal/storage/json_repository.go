// Package storage loads PhraseForge data from external sources.
package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

type categoriesFile struct {
	Categories []phrase.Category `json:"categories"`
}

// LoadCategories reads and validates phrase categories from a JSON file.
func LoadCategories(path string) ([]phrase.Category, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open categories file %q: %w", path, err)
	}
	defer file.Close()

	var data categoriesFile
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode categories file %q: %w", path, err)
	}

	if err := validateCategories(data.Categories); err != nil {
		return nil, err
	}

	return data.Categories, nil
}

func validateCategories(categories []phrase.Category) error {
	if len(categories) == 0 {
		return fmt.Errorf("validate categories: categories cannot be empty")
	}

	for index, category := range categories {
		if err := phrase.ValidateCategory(category); err != nil {
			return fmt.Errorf("validate categories: category %d: %w", index+1, err)
		}
	}

	return nil
}
