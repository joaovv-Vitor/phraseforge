package storage

import (
	"fmt"
	"time"
)

func parseSQLiteTimestamp(value string) (time.Time, error) {
	timestamp, err := time.Parse(time.DateTime, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse %q: %w", value, err)
	}

	return timestamp, nil
}
