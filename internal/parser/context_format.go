package parser

import (
	"fmt"
	"io"
	"strings"
)

// PrintContext writes context results to w in a human-readable format.
func PrintContext(w io.Writer, results []ContextLine) {
	for i, r := range results {
		if i > 0 {
			fmt.Fprintln(w, "--")
		}
		for _, b := range r.Before {
			fmt.Fprintf(w, "  %s\n", b.Raw)
		}
		fmt.Fprintf(w, "> %s\n", r.Line.Raw)
		for _, a := range r.After {
			fmt.Fprintf(w, "  %s\n", a.Raw)
		}
	}
}

// ContextJSON returns a JSON array of context match objects.
func ContextJSON(results []ContextLine) string {
	if len(results) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, r := range results {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("{")
		sb.WriteString(`"matched":` + jsonEscapeString(r.Line.Raw) + ",")
		sb.WriteString(`"before":[`)
		for j, b := range r.Before {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(jsonEscapeString(b.Raw))
		}
		sb.WriteString(`],"after":[`)
		for j, a := range r.After {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(jsonEscapeString(a.Raw))
		}
		sb.WriteString(`]}`)
	}
	sb.WriteString("]")
	return sb.String()
}
