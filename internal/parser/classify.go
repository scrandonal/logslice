package parser

import (
	"regexp"
	"strings"
)

// Category represents a named log classification rule.
type Category struct {
	Name    string
	Pattern *regexp.Regexp
}

// ClassifyResult holds a log line along with its matched category.
type ClassifyResult struct {
	Line     Line
	Category string // empty string means no category matched
}

// Classifier assigns categories to log lines based on regex rules.
type Classifier struct {
	categories []Category
	defaultCat string
}

// NewClassifier creates a Classifier from a map of name->pattern strings.
// Patterns are compiled case-insensitively. Returns an error if any pattern
// is invalid. defaultCategory is used when no rule matches.
func NewClassifier(rules map[string]string, defaultCategory string) (*Classifier, error) {
	cats := make([]Category, 0, len(rules))
	for name, pat := range rules {
		re, err := regexp.Compile("(?i)" + pat)
		if err != nil {
			return nil, err
		}
		cats = append(cats, Category{Name: name, Pattern: re})
	}
	return &Classifier{categories: cats, defaultCat: defaultCategory}, nil
}

// Classify assigns a category to a single line. Rules are evaluated in
// insertion order; the first match wins.
func (c *Classifier) Classify(l Line) ClassifyResult {
	for _, cat := range c.categories {
		if cat.Pattern.MatchString(l.Raw) {
			return ClassifyResult{Line: l, Category: cat.Name}
		}
	}
	return ClassifyResult{Line: l, Category: c.defaultCat}
}

// ClassifyLines classifies a slice of lines and returns results.
func ClassifyLines(lines []Line, rules map[string]string, defaultCategory string) ([]ClassifyResult, error) {
	c, err := NewClassifier(rules, defaultCategory)
	if err != nil {
		return nil, err
	}
	out := make([]ClassifyResult, len(lines))
	for i, l := range lines {
		out[i] = c.Classify(l)
	}
	return out, nil
}

// GroupByCategory partitions ClassifyResults into a map keyed by category.
func GroupByCategory(results []ClassifyResult) map[string][]ClassifyResult {
	m := make(map[string][]ClassifyResult)
	for _, r := range results {
		key := r.Category
		if key == "" {
			key = "(none)"
		}
		m[key] = append(m[key], r)
	}
	return m
}

// CategorySummary returns a count of lines per category, sorted by category name.
func CategorySummary(results []ClassifyResult) map[string]int {
	counts := make(map[string]int)
	for _, r := range results {
		key := r.Category
		if strings.TrimSpace(key) == "" {
			key = "(none)"
		}
		counts[key]++
	}
	return counts
}
