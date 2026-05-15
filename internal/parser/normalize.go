package parser

import (
	"regexp"
	"strings"
)

// NormalizeRule defines a substitution applied to log lines.
type NormalizeRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Normalizer applies a set of normalization rules to log lines,
// replacing dynamic tokens (IDs, IPs, numbers) with stable placeholders.
type Normalizer struct {
	rules []NormalizeRule
}

// NewNormalizer constructs a Normalizer from a slice of pattern/replacement pairs.
// Each pattern is a regular expression string. Returns an error if any pattern
// fails to compile.
func NewNormalizer(rules [][2]string) (*Normalizer, error) {
	compiled := make([]NormalizeRule, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile(r[0])
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, NormalizeRule{Pattern: re, Replacement: r[1]})
	}
	return &Normalizer{rules: compiled}, nil
}

// Apply runs all normalization rules against raw and returns the result.
func (n *Normalizer) Apply(raw string) string {
	for _, r := range n.rules {
		raw = r.Pattern.ReplaceAllString(raw, r.Replacement)
	}
	return raw
}

// NormalizeLine applies the normalizer to a single Line's Raw field.
func (n *Normalizer) NormalizeLine(l Line) Line {
	l.Raw = n.Apply(l.Raw)
	return l
}

// NormalizeLines applies the normalizer to every line in the slice.
func NormalizeLines(lines []Line, rules [][2]string) ([]Line, error) {
	norm, err := NewNormalizer(rules)
	if err != nil {
		return nil, err
	}
	out := make([]Line, len(lines))
	for i, l := range lines {
		out[i] = norm.NormalizeLine(l)
	}
	return out, nil
}

// DefaultNormalizer returns a Normalizer with common rules for UUIDs, IPv4
// addresses, integers, and hex tokens.
func DefaultNormalizer() *Normalizer {
	n, _ := NewNormalizer([][2]string{
		{`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`, "<UUID>"},
		{`\b(?:\d{1,3}\.){3}\d{1,3}\b`, "<IP>"},
		{`\b0x[0-9a-fA-F]+\b`, "<HEX>"},
		{`\b\d{4,}\b`, "<NUM>"},
	})
	return n
}

// NormalizedKey returns a lowercase, whitespace-collapsed version of s,
// useful as a map key after normalization.
func NormalizedKey(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}
