package parser

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// PrintFields writes extracted fields for each line as human-readable output.
func PrintFields(w io.Writer, lines []Line, fields []string) {
	fe := NewFieldExtractor(fields)
	for _, l := range lines {
		extracted := fe.Extract(l.Text)
		if len(extracted) == 0 {
			continue
		}
		parts := make([]string, 0, len(extracted))
		for _, f := range fields {
			if v, ok := extracted[f]; ok {
				parts = append(parts, f+"="+v)
			}
		}
		fmt.Fprintln(w, strings.Join(parts, " "))
	}
}

// FieldsJSON returns a JSON array of objects, one per line with extracted fields.
func FieldsJSON(lines []Line, fields []string) string {
	fe := NewFieldExtractor(fields)
	var sb strings.Builder
	sb.WriteString("[")
	first := true
	for _, l := range lines {
		extracted := fe.Extract(l.Text)
		if len(extracted) == 0 {
			continue
		}
		if !first {
			sb.WriteString(",")
		}
		first = false
		sb.WriteString("{")
		keys := make([]string, 0, len(extracted))
		for k := range extracted {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(`"`)
			sb.WriteString(jsonEscapeString(k))
			sb.WriteString(`":"`)
			sb.WriteString(jsonEscapeString(extracted[k]))
			sb.WriteString(`"`)
		}
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}
