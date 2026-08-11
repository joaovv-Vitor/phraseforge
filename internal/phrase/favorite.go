package phrase

import "time"

// Favorite is a generated phrase saved by the user.
type Favorite struct {
	ID        int64
	Category  string
	Content   string
	CreatedAt time.Time
}
