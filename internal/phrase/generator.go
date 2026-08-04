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
	if err := validateParts(template, parts); err != nil {
		return "", fmt.Errorf("generate phrase: %w", err)
	}

	subject := parts.Subjects[rand.IntN(len(parts.Subjects))]
	verb := parts.Verbs[rand.IntN(len(parts.Verbs))]
	complement := parts.Complements[rand.IntN(len(parts.Complements))]

	replacements := []string{
		"{subject}", subject,
		"{verb}", verb,
		"{complement}", complement,
	}
	if strings.Contains(template, "{introduction}") {
		introduction := parts.Introductions[rand.IntN(len(parts.Introductions))]
		replacements = append(replacements, "{introduction}", introduction)
	}
	if strings.Contains(template, "{conclusion}") {
		conclusion := parts.Conclusions[rand.IntN(len(parts.Conclusions))]
		replacements = append(replacements, "{conclusion}", conclusion)
	}

	return strings.NewReplacer(replacements...).Replace(template) + ".", nil
}
