// Package phrase contains the core types and operations for phrase generation.
package phrase

// Parts contains the components used to build a phrase.
type Parts struct {
	Subjects    []string
	Verbs       []string
	Complements []string
}
