package phrase

import (
	"errors"
	"time"
)

var (
	// ErrFavoriteAlreadyExists indicates that the phrase is already saved for a category.
	ErrFavoriteAlreadyExists = errors.New("favorite already exists")
	// ErrFavoriteCategoryNotFound indicates that a favorite references an unknown category.
	ErrFavoriteCategoryNotFound = errors.New("favorite category not found")
)

// Favorite is a generated phrase saved by the user.
type Favorite struct {
	ID        int64
	Category  string
	Content   string
	CreatedAt time.Time
}
