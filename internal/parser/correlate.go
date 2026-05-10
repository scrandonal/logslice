package parser

import (
	"regexp"
	"time"
)

// CorrelateResult holds a group of lines sharing the same correlation key.
type CorrelateResult struct {
	Key   string
	Lines []Line
	First *time.Time
	Last  *time.Time
	Count int
}

// Correlator groups log lines by a captured regex group.
type Correlator struct {
	re      *regexp.Regexp
	groups  map[string]*CorrelateResult
	ordered []string
}

// NewCorrelator creates a Correlator using the first capture group of pattern.
// Returns nil and an error if the pattern is invalid or has no capture group.
func NewCorrelator(pattern string) (*Correlator, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Correlator{
		re:     re,
		groups: make(map[string]*CorrelateResult),
	}, nil
}

// Add processes a single line and groups it by its correlation key.
func (c *Correlator) Add(l Line) {
	m := c.re.FindStringSubmatch(l.Raw)
	var key string
	if len(m) >= 2 {
		key = m[1]
	} else {
		key = "(unmatched)"
	}
	if _, ok := c.groups[key]; !ok {
		c.groups[key] = &CorrelateResult{Key: key}
		c.ordered = append(c.ordered, key)
	}
	g := c.groups[key]
	g.Lines = append(g.Lines, l)
	g.Count++
	if l.Timestamp != nil {
		if g.First == nil || l.Timestamp.Before(*g.First) {
			t := *l.Timestamp
			g.First = &t
		}
		if g.Last == nil || l.Timestamp.After(*g.Last) {
			t := *l.Timestamp
			g.Last = &t
		}
	}
}

// Results returns correlation groups in insertion order.
func (c *Correlator) Results() []*CorrelateResult {
	out := make([]*CorrelateResult, 0, len(c.ordered))
	for _, k := range c.ordered {
		out = append(out, c.groups[k])
	}
	return out
}

// CorrelateLines is a convenience wrapper that runs all lines through a new Correlator.
func CorrelateLines(lines []Line, pattern string) ([]*CorrelateResult, error) {
	c, err := NewCorrelator(pattern)
	if err != nil {
		return nil, err
	}
	for _, l := range lines {
		c.Add(l)
	}
	return c.Results(), nil
}
