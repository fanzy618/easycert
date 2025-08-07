package cert

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration converts a human-friendly duration string into a time.Duration.
// It supports the suffixes "y" for years (365 days) and "d" for days in addition
// to the standard time.ParseDuration formats like "720h" or "15m".
func ParseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "y") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "y"))
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}
		return time.Duration(n) * 365 * 24 * time.Hour, nil
	}
	if strings.HasSuffix(s, "d") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", s, err)
	}
	return d, nil
}
