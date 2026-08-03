package phrase

import (
	"fmt"
	"strings"
)

// ValidateCategory checks whether a category can generate complete phrases.
func ValidateCategory(category Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	if err := validateTemplate(category.Template); err != nil {
		return fmt.Errorf("category %q: %w", category.Name, err)
	}
	if err := validateParts(category.Parts); err != nil {
		return fmt.Errorf("category %q: %w", category.Name, err)
	}

	return nil
}

func validateTemplate(template string) error {
	if strings.TrimSpace(template) == "" {
		return fmt.Errorf("template cannot be empty")
	}

	remaining := template
	for _, placeholder := range []string{"{subject}", "{verb}", "{complement}"} {
		if strings.Count(template, placeholder) != 1 {
			return fmt.Errorf("template must contain %q exactly once", placeholder)
		}
		remaining = strings.ReplaceAll(remaining, placeholder, "")
	}
	if strings.ContainsAny(remaining, "{}") {
		return fmt.Errorf("template contains an unknown placeholder")
	}

	return nil
}

func validateParts(parts Parts) error {
	if len(parts.Subjects) == 0 {
		return fmt.Errorf("subjects cannot be empty")
	}
	if len(parts.Verbs) == 0 {
		return fmt.Errorf("verbs cannot be empty")
	}
	if len(parts.Complements) == 0 {
		return fmt.Errorf("complements cannot be empty")
	}

	return nil
}
