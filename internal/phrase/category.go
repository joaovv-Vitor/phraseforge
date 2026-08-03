package phrase

import "fmt"

// FindCategory returns the category with the requested name.
func FindCategory(categories []Category, name string) (Category, error) {
	for _, category := range categories {
		if category.Name == name {
			return category, nil
		}
	}

	return Category{}, fmt.Errorf("find category: category %q not found", name)
}
