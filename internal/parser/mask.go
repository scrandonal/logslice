package parser

import (
	"regexp"
	"strings"
)

// MaskMode controls how sensitive data is replaced.
type MaskMode int

const (
	MaskRedact  MaskMode = iota // replace with [REDACTED]
	MaskPartial                 // keep first/last char, mask middle with ***
)

// Masker redacts sensitive patterns from log lines.
type Masker struct {
	patterns []*regexp.Regexp
	mode     MaskMode
}

// NewMasker compiles the given regex patterns and returns a Masker.
// Returns an error if any pattern fails to compile.
func NewMasker(patterns []string, mode MaskMode) (*Masker, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Masker{patterns: compiled, mode: mode}, nil
}

// Mask applies all patterns to line and returns the redacted string.
func (m *Masker) Mask(line string) string {
	for _, re := range m.patterns {
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			return m.replace(match)
		})
	}
	return line
}

func (m *Masker) replace(s string) string {
	switch m.mode {
	case MaskPartial:
		if len(s) <= 2 {
			return strings.Repeat("*", len(s))
		}
		return string(s[0]) + strings.Repeat("*", len(s)-2) + string(s[len(s)-1])
	default:
		return "[REDACTED]"
	}
}

// MaskLines applies the masker to a slice of raw log lines.
func MaskLines(lines []Line, patterns []string, mode MaskMode) ([]Line, error) {
	m, err := NewMasker(patterns, mode)
	if err != nil {
		return nil, err
	}
	out := make([]Line, len(lines))
	for i, l := range lines {
		out[i] = Line{
			Raw:       m.Mask(l.Raw),
			Timestamp: l.Timestamp,
		}
	}
	return out, nil
}
