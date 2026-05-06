package parser

import (
	"fmt"
	"io"
	"strings"
)

// PrintSampled writes sampled lines to w in a human-readable format.
// Each line is prefixed with its 1-based index.
func PrintSampled(w io.Writer, lines []Line) {
	for i, l := range lines {
		ts := "-"
		if l.Timestamp != nil {
			ts = l.Timestamp.UTC().Format("2006-01-02T15:04:05Z")
		}
		fmt.Fprintf(w, "%4d  [%s]  %s\n", i+1, ts, l.Raw)
	}
}

// SampledJSON returns a JSON array of sampled lines.
func SampledJSON(lines []Line) string {
	if len(lines) == 0 {
		return "[]\n"
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, l := range lines {
		ts := "null"
		if l.Timestamp != nil {
			ts = `"` + l.Timestamp.UTC().Format("2006-01-02T15:04:05Z") + `"`
		}
		sb.WriteString(fmt.Sprintf(
			"  {\"index\":%d,\"timestamp\":%s,\"raw\":\"%s\"}",
			i+1, ts, jsonEscapeString(l.Raw),
		))
		if i < len(lines)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	return sb.String()
}
