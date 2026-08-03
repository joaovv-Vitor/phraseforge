package phrase

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Generate builds a phrase by randomly choosing one value from each required part.
func Generate(parts Parts) (string, error) {
	if len(parts.Subjects) == 0 {
		return "", fmt.Errorf("generate phrase: subjects cannot be empty")
	}
	if len(parts.Verbs) == 0 {
		return "", fmt.Errorf("generate phrase: verbs cannot be empty")
	}
	if len(parts.Complements) == 0 {
		return "", fmt.Errorf("generate phrase: complements cannot be empty")
	}

	subject := parts.Subjects[rand.IntN(len(parts.Subjects))]
	verb := parts.Verbs[rand.IntN(len(parts.Verbs))]
	complement := parts.Complements[rand.IntN(len(parts.Complements))]

	return fmt.Sprintf("%s.", strings.Join([]string{subject, verb, complement}, " ")), nil
}
