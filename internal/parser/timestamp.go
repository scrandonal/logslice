package parser

import (
	"errors"
	"time"
)

// Common log timestamp formats to attempt parsing
var timestampFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// ErrNoTimestamp is returned when no recognizable timestamp is found in a line.
var ErrNoTimestamp = errors.New("no timestamp found in line")

// ParseTimestamp attempts to extract and parse a timestamp from a log line.
// It tries each known format against the beginning of the line and common
// bracket/space-delimited prefixes.
func ParseTimestamp(line string) (time.Time, error) {
	if len(line) == 0 {
		return time.Time{}, ErrNoTimestamp
	}

	// Strip leading bracket if present (e.g. [2024-01-02T15:04:05Z])
	candidate := line
	if line[0] == '[' {
		end := indexByte(line, ']')
		if end > 1 {
			candidate = line[1:end]
		}
	}

	for _, format := range timestampFormats {
		prefixLen := len(format)
		if prefixLen > len(candidate) {
			prefixLen = len(candidate)
		}
		// Try progressively shorter prefixes down to format length
		for l := len(candidate); l >= prefixLen; l-- {
			t, err := time.Parse(format, candidate[:l])
			if err == nil {
				return t, nil
			}
		}
	}

	return time.Time{}, ErrNoTimestamp
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
