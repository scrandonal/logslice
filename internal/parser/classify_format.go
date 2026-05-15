package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// PrintClassifySummary writes a human-readable category summary to stdout.
func PrintClassifySummary(results []ClassifyResult) {
	printClassifySummaryTo(os.Stdout, results)
}

func printClassifySummaryTo(w io.Writer, results []ClassifyResult) {
	counts := CategorySummary(results)
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintln(w, "Category Summary:")
	for _, k := range keys {
		fmt.Fprintf(w, "  %-20s %d\n", k, counts[k])
	}
}

// PrintClassifyResults writes each line prefixed with its category.
func PrintClassifyResults(results []ClassifyResult) {
	for _, r := range results {
		cat := r.Category
		if cat == "" {
			cat = "(none)"
		}
		fmt.Printf("[%s] %s\n", cat, r.Line.Raw)
	}
}

type classifyJSONEntry struct {
	Category  string `json:"category"`
	Raw       string `json:"raw"`
	Timestamp string `json:"timestamp,omitempty"`
}

// ClassifyJSON serialises classify results to a JSON array string.
func ClassifyJSON(results []ClassifyResult) string {
	entries := make([]classifyJSONEntry, len(results))
	for i, r := range results {
		cat := r.Category
		if cat == "" {
			cat = "(none)"
		}
		e := classifyJSONEntry{Category: cat, Raw: r.Line.Raw}
		if r.Line.Timestamp != nil {
			e.Timestamp = r.Line.Timestamp.Format("2006-01-02T15:04:05Z07:00")
		}
		entries[i] = e
	}
	b, err := json.Marshal(entries)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// ClassifySummaryJSON serialises the category counts to a JSON object string.
func ClassifySummaryJSON(results []ClassifyResult) string {
	counts := CategorySummary(results)
	b, err := json.Marshal(counts)
	if err != nil {
		return "{}"
	}
	return string(b)
}
