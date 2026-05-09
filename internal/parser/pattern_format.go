package parser

import (
	"fmt"
	"io"
	"strings"
)

// PrintPatternCounts writes a human-readable pattern count summary to w.
func PrintPatternCounts(w io.Writer, counts map[string]int) {
	if len(counts) == 0 {
		fmt.Fprintln(w, "no patterns matched")
		return
	}
	names := patternNamesSorted(counts)
	for _, name := range names {
		fmt.Fprintf(w, "%-30s %d\n", name, counts[name])
	}
}

// PrintPatternMatches writes each matched line prefixed by its pattern name.
func PrintPatternMatches(w io.Writer, matches []PatternMatch) {
	for _, m := range matches {
		fmt.Fprintf(w, "[%s] %s\n", m.Pattern, m.Line.Raw)
	}
}

// PatternCountsJSON returns a JSON object mapping pattern names to counts.
func PatternCountsJSON(counts map[string]int) string {
	if len(counts) == 0 {
		return "{}"
	}
	names := patternNamesSorted(counts)
	var sb strings.Builder
	sb.WriteString("{")
	for i, name := range names {
		if i > 0 {
			sb.WriteString(",")
		}
		fmt.Fprintf(&sb, "%q:%d", name, counts[name])
	}
	sb.WriteString("}")
	return sb.String()
}

// PatternMatchesJSON returns a JSON array of matched line objects.
func PatternMatchesJSON(matches []PatternMatch) string {
	if len(matches) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, m := range matches {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("{")
		fmt.Fprintf(&sb, "%q:%q,", "pattern", m.Pattern)
		fmt.Fprintf(&sb, "%q:%q,", "raw", jsonEscapeString(m.Line.Raw))
		fmt.Fprintf(&sb, "%q:%d", "count", m.Count)
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}
