package parser

import (
	"regexp"
	"strings"
)

// GrepMode controls how pattern matching is applied.
type GrepMode int

const (
	GrepLiteral GrepMode = iota
	GrepRegex
	GrepInvert
)

// Grepper filters log lines by pattern match against the raw text.
type Grepper struct {
	mode    GrepMode
	pattern string
	re      *regexp.Regexp
}

// NewGrepper creates a Grepper for the given pattern and mode.
// For GrepRegex mode the pattern is compiled; returns error on bad regex.
func NewGrepper(pattern string, mode GrepMode) (*Grepper, error) {
	g := &Grepper{mode: mode, pattern: pattern}
	if mode == GrepRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		g.re = re
	}
	return g, nil
}

// Match reports whether the raw line satisfies the grep condition.
func (g *Grepper) Match(raw string) bool {
	switch g.mode {
	case GrepRegex:
		return g.re.MatchString(raw)
	case GrepInvert:
		return !strings.Contains(raw, g.pattern)
	default: // GrepLiteral
		return strings.Contains(raw, g.pattern)
	}
}

// GrepLines filters a slice of Lines, returning only those whose Raw field
// satisfies the Grepper's condition.
func GrepLines(lines []Line, g *Grepper) []Line {
	out := make([]Line, 0, len(lines))
	for _, l := range lines {
		if g.Match(l.Raw) {
			out = append(out, l)
		}
	}
	return out
}
