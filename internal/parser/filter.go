package parser

import (
	"time"
)

// Filter holds the time range used to select log lines.
type Filter struct {
	From time.Time
	To   time.Time
}

// NewFilter creates a Filter from two RFC3339 (or compatible) timestamp strings.
func NewFilter(from, to string) (Filter, error) {
	f, err := ParseTimestamp(from)
	if err != nil {
		return Filter{}, err
	}
	t, err := ParseTimestamp(to)
	if err != nil {
		return Filter{}, err
	}
	return Filter{From: f, To: t}, nil
}

// Match reports whether ts falls within the filter's inclusive time range.
func (f Filter) Match(ts time.Time) bool {
	if ts.IsZero() {
		return false
	}
	return !ts.Before(f.From) && !ts.After(f.To)
}

// Slice reads lines from the scanner and returns only those whose timestamps
// fall within the filter range. It stops early once lines pass the upper bound.
func (f Filter) Slice(s *Scanner) []string {
	var result []string
	for s.Scan() {
		line := s.Line()
		ts := s.Timestamp()
		if ts.IsZero() {
			continue
		}
		if ts.After(f.To) {
			break
		}
		if f.Match(ts) {
			result = append(result, line)
		}
	}
	return result
}
