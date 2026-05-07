package parser

import "strings"

// TruncateMode controls how long lines are shortened.
type TruncateMode int

const (
	TruncateNone   TruncateMode = iota
	TruncateRight               // cut from the right (default)
	TruncateMiddle              // preserve start and end, replace centre with ellipsis
)

// Truncator shortens log lines that exceed a maximum byte length.
type Truncator struct {
	maxLen  int
	mode    TruncateMode
	ellipsis string
}

// NewTruncator creates a Truncator.
// maxLen <= 0 disables truncation.
func NewTruncator(maxLen int, mode TruncateMode) *Truncator {
	return &Truncator{
		maxLen:   maxLen,
		mode:     mode,
		ellipsis: "...",
	}
}

// Truncate shortens s if it exceeds the configured maximum length.
// If maxLen <= 0 the original string is returned unchanged.
func (t *Truncator) Truncate(s string) string {
	if t.maxLen <= 0 || len(s) <= t.maxLen {
		return s
	}

	elLen := len(t.ellipsis)

	switch t.mode {
	case TruncateMiddle:
		if t.maxLen <= elLen {
			return t.ellipsis[:t.maxLen]
		}
		avail := t.maxLen - elLen
		left := avail / 2
		right := avail - left
		return s[:left] + t.ellipsis + s[len(s)-right:]
	default: // TruncateRight
		if t.maxLen <= elLen {
			return t.ellipsis[:t.maxLen]
		}
		return s[:t.maxLen-elLen] + t.ellipsis
	}
}

// TruncateLines applies Truncate to the Raw field of each Line.
func TruncateLines(lines []Line, maxLen int, mode TruncateMode) []Line {
	if maxLen <= 0 {
		return lines
	}
	t := NewTruncator(maxLen, mode)
	out := make([]Line, len(lines))
	for i, l := range lines {
		l.Raw = t.Truncate(strings.TrimRight(l.Raw, "\n"))
		out[i] = l
	}
	return out
}
