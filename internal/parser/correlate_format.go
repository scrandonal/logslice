package parser

import (
	"fmt"
	"io"
	"strings"
)

// PrintCorrelate writes a human-readable summary of correlation results to w.
func PrintCorrelate(w io.Writer, results []*CorrelateResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no correlation groups found")
		return
	}
	for _, r := range results {
		first := "null"
		last := "null"
		if r.First != nil {
			first = r.First.Format("2006-01-02T15:04:05Z")
		}
		if r.Last != nil {
			last = r.Last.Format("2006-01-02T15:04:05Z")
		}
		fmt.Fprintf(w, "[%s] count=%d first=%s last=%s\n", r.Key, r.Count, first, last)
		for _, l := range r.Lines {
			fmt.Fprintf(w, "  %s\n", l.Raw)
		}
	}
}

// CorrelateJSON returns a JSON array of correlation results.
func CorrelateJSON(results []*CorrelateResult) string {
	if len(results) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, r := range results {
		if i > 0 {
			sb.WriteString(",")
		}
		first := "null"
		last := "null"
		if r.First != nil {
			first = fmt.Sprintf("%q", r.First.Format("2006-01-02T15:04:05Z"))
		}
		if r.Last != nil {
			last = fmt.Sprintf("%q", r.Last.Format("2006-01-02T15:04:05Z"))
		}
		sb.WriteString(fmt.Sprintf(
			`{"key":%q,"count":%d,"first":%s,"last":%s,"lines":[`,
			r.Key, r.Count, first, last,
		))
		for j, l := range r.Lines {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("%q", l.Raw))
		}
		sb.WriteString("]}") 
	}
	sb.WriteString("]")
	return sb.String()
}
