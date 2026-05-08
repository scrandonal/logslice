package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// PrintPivot writes a human-readable pivot table to w.
func PrintPivot(w io.Writer, result *PivotResult) {
	if result == nil || len(result.Rows) == 0 {
		fmt.Fprintln(w, "(no pivot data)")
		return
	}

	// Header
	fmt.Fprintf(w, "pivot: field=%q\n", result.Field)
	header := fmt.Sprintf("%-20s  %6s", "value", "total")
	for _, b := range result.Buckets {
		short := b
		if len(short) > 8 {
			short = short[len(short)-8:]
		}
		header += fmt.Sprintf("  %8s", short)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, row := range result.Rows {
		line := fmt.Sprintf("%-20s  %6d", truncateStr(row.Value, 20), row.Total)
		for _, b := range result.Buckets {
			line += fmt.Sprintf("  %8d", row.Buckets[b])
		}
		fmt.Fprintln(w, line)
	}
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// PivotJSON serialises the pivot result as JSON.
func PivotJSON(result *PivotResult) (string, error) {
	type row struct {
		Value   string         `json:"value"`
		Total   int            `json:"total"`
		Buckets map[string]int `json:"buckets"`
	}
	type out struct {
		Field   string   `json:"field"`
		Buckets []string `json:"buckets"`
		Rows    []row    `json:"rows"`
	}

	o := out{Field: result.Field, Buckets: result.Buckets}
	for _, r := range result.Rows {
		o.Rows = append(o.Rows, row{Value: r.Value, Total: r.Total, Buckets: r.Buckets})
	}
	if o.Rows == nil {
		o.Rows = []row{}
	}
	if o.Buckets == nil {
		o.Buckets = []string{}
	}

	b, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
