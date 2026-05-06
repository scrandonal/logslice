package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// MergedLineJSON is the JSON representation of a MergedLine.
type MergedLineJSON struct {
	Source    int    `json:"source"`
	Timestamp string `json:"timestamp,omitempty"`
	Raw       string `json:"raw"`
}

// PrintMerged writes merged lines to w in raw format, prefixed by source index.
func PrintMerged(w io.Writer, lines []MergedLine) {
	for _, ml := range lines {
		ts := "-"
		if ml.Timestamp != nil {
			ts = ml.Timestamp.UTC().Format("2006-01-02T15:04:05Z")
		}
		fmt.Fprintf(w, "[src:%d] [%s] %s\n", ml.Source, ts, ml.Raw)
	}
}

// MergedJSON serialises merged lines as a JSON array.
func MergedJSON(lines []MergedLine) (string, error) {
	out := make([]MergedLineJSON, 0, len(lines))
	for _, ml := range lines {
		entry := MergedLineJSON{
			Source: ml.Source,
			Raw:    ml.Raw,
		}
		if ml.Timestamp != nil {
			entry.Timestamp = ml.Timestamp.UTC().Format("2006-01-02T15:04:05Z")
		}
		out = append(out, entry)
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
