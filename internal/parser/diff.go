package parser

import (
	"time"
)

// DiffEntry represents a line that appears in one log but not the other.
type DiffEntry struct {
	Line    Line
	Source  string // "left" or "right"
}

// DiffResult holds the output of comparing two log streams.
type DiffResult struct {
	OnlyLeft  []Line
	OnlyRight []Line
	Common    []Line
}

// DiffOptions controls how log lines are compared.
type DiffOptions struct {
	// Window is the time window within which lines are considered for matching.
	Window time.Duration
	// IgnoreTimestamp strips timestamps before comparing raw text.
	IgnoreTimestamp bool
}

// DiffLogs compares two slices of Lines and returns a DiffResult.
// Lines are matched by their raw text (optionally ignoring timestamps).
// Lines within the time Window of each other are eligible for matching.
func DiffLogs(left, right []Line, opts DiffOptions) DiffResult {
	if opts.Window == 0 {
		opts.Window = 5 * time.Minute
	}

	keyOf := func(l Line) string {
		if opts.IgnoreTimestamp && l.Timestamp != nil {
			// Strip leading timestamp bracket if present.
			raw := l.Raw
			if len(raw) > 0 && raw[0] == '[' {
				if idx := indexByte(raw, ']'); idx >= 0 && idx+2 < len(raw) {
					return raw[idx+2:]
				}
			}
			return raw
		}
		return l.Raw
	}

	rightUsed := make([]bool, len(right))
	var onlyLeft, common []Line

	for _, l := range left {
		lKey := keyOf(l)
		matched := false
		for i, r := range right {
			if rightUsed[i] {
				continue
			}
			if keyOf(r) != lKey {
				continue
			}
			if l.Timestamp != nil && r.Timestamp != nil {
				diff := l.Timestamp.Sub(*r.Timestamp)
				if diff < 0 {
					diff = -diff
				}
				if diff > opts.Window {
					continue
				}
			}
			rightUsed[i] = true
			matched = true
			common = append(common, l)
			break
		}
		if !matched {
			onlyLeft = append(onlyLeft, l)
		}
	}

	var onlyRight []Line
	for i, r := range right {
		if !rightUsed[i] {
			onlyRight = append(onlyRight, r)
		}
	}

	return DiffResult{
		OnlyLeft:  onlyLeft,
		OnlyRight: onlyRight,
		Common:    common,
	}
}
