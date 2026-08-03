package phrase

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Generate builds a phrase by rendering a template with randomly selected parts.
func Generate(template string, parts Parts) (string, error) {
	if err := validateTemplate(template); err != nil {
		return "", fmt.Errorf("generate phrase: %w", err)
	}
	if err := validateParts(parts); err != nil {
		return "", fmt.Errorf("generate phrase: %w", err)
	}

	subject := parts.Subjects[rand.IntN(len(parts.Subjects))]
	verb := parts.Verbs[rand.IntN(len(parts.Verbs))]
	complement := parts.Complements[rand.IntN(len(parts.Complements))]

	return strings.NewReplacer(
		"{subject}", subject,
		"{verb}", verb,
		"{complement}", complement,
	).Replace(template) + ".", nil
}
