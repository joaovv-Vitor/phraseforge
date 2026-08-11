package phrase

import (
	"errors"
	"time"
)

var (
	// ErrHistoryCategoryNotFound indicates that a history entry references an unknown category.
	ErrHistoryCategoryNotFound = errors.New("history category not found")
)

// HistoryEntry records a generated phrase.
type HistoryEntry struct {
	ID          int64
	Category    string
	Content     string
	GeneratedAt time.Time
}
