package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// PrintRate writes a human-readable rate table to stdout.
func PrintRate(points []RatePoint, interval time.Duration, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if len(points) == 0 {
		fmt.Fprintln(w, "no data")
		return
	}
	fmt.Fprintf(w, "%-30s  %s\n", "bucket", "count")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 42))
	for _, p := range points {
		label := formatBucketLabel(p.Bucket, interval)
		fmt.Fprintf(w, "%-30s  %d\n", label, p.Count)
	}
}

func formatBucketLabel(t time.Time, interval time.Duration) string {
	if t.IsZero() {
		return "(no timestamp)"
	}
	if interval >= 24*time.Hour {
		return t.Format("2006-01-02")
	}
	if interval >= time.Hour {
		return t.Format("2006-01-02 15:00")
	}
	return t.Format("2006-01-02 15:04")
}

type ratePointJSON struct {
	Bucket string `json:"bucket"`
	Count  int    `json:"count"`
}

// RateJSON serialises rate points to a JSON array string.
func RateJSON(points []RatePoint, interval time.Duration) string {
	out := make([]ratePointJSON, len(points))
	for i, p := range points {
		out[i] = ratePointJSON{
			Bucket: formatBucketLabel(p.Bucket, interval),
			Count:  p.Count,
		}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "[]"
	}
	return string(b)
}
