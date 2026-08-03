// Package phrase contains the core types and operations for phrase generation.
package phrase

// Parts contains the components used to build a phrase.
type Parts struct {
	Subjects    []string `json:"subjects"`
	Verbs       []string `json:"verbs"`
	Complements []string `json:"complements"`
}

// Category groups phrase parts under a name.
type Category struct {
	Name     string `json:"name"`
	Template string `json:"template"`
	Parts
}
