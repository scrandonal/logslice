package parser

import (
	"regexp"
	"strings"
)

// PatternMatch holds a single matched pattern result.
type PatternMatch struct {
	Line    Line
	Pattern string
	Count   int
}

// PatternCounter counts occurrences of named patterns across log lines.
type PatternCounter struct {
	patterns map[string]*regexp.Regexp
	counts   map[string]int
	matches  []PatternMatch
}

// NewPatternCounter creates a PatternCounter from a map of name->regex strings.
// Returns an error if any pattern fails to compile.
func NewPatternCounter(patterns map[string]string) (*PatternCounter, error) {
	compiled := make(map[string]*regexp.Regexp, len(patterns))
	for name, expr := range patterns {
		re, err := regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
		compiled[name] = re
	}
	return &PatternCounter{
		patterns: compiled,
		counts:   make(map[string]int, len(patterns)),
	}, nil
}

// Feed processes a single line against all patterns.
func (pc *PatternCounter) Feed(line Line) {
	for name, re := range pc.patterns {
		if re.MatchString(line.Raw) {
			pc.counts[name]++
			pc.matches = append(pc.matches, PatternMatch{
				Line:    line,
				Pattern: name,
				Count:   pc.counts[name],
			})
		}
	}
}

// Counts returns a copy of the pattern hit counts.
func (pc *PatternCounter) Counts() map[string]int {
	out := make(map[string]int, len(pc.counts))
	for k, v := range pc.counts {
		out[k] = v
	}
	return out
}

// Matches returns all matched lines.
func (pc *PatternCounter) Matches() []PatternMatch {
	return pc.matches
}

// CountPatterns is a convenience helper that runs all lines through a PatternCounter.
func CountPatterns(lines []Line, patterns map[string]string) ([]PatternMatch, map[string]int, error) {
	pc, err := NewPatternCounter(patterns)
	if err != nil {
		return nil, nil, err
	}
	for _, l := range lines {
		pc.Feed(l)
	}
	return pc.Matches(), pc.Counts(), nil
}

// patternNamesSorted returns pattern names sorted for deterministic output.
func patternNamesSorted(counts map[string]int) []string {
	names := make([]string, 0, len(counts))
	for k := range counts {
		names = append(names, k)
	}
	// simple insertion sort — pattern counts are small
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && strings.ToLower(names[j-1]) > strings.ToLower(names[j]); j-- {
			names[j-1], names[j] = names[j], names[j-1]
		}
	}
	return names
}
