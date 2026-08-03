package phrase

// CategoryNames returns category names in the same order as the input.
func CategoryNames(categories []Category) []string {
	names := make([]string, 0, len(categories))
	for _, category := range categories {
		names = append(names, category.Name)
	}

	return names
}
