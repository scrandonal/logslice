package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PrintLatency writes a human-readable latency summary to stdout.
func PrintLatency(stats LatencyStats) {
	printLatencyTo(os.Stdout, stats)
}

func printLatencyTo(w io.Writer, stats LatencyStats) {
	if stats.Count == 0 {
		fmt.Fprintln(w, "no latency data")
		return
	}
	fmt.Fprintf(w, "count : %d\n", stats.Count)
	fmt.Fprintf(w, "min   : %s\n", stats.Min)
	fmt.Fprintf(w, "max   : %s\n", stats.Max)
	fmt.Fprintf(w, "mean  : %s\n", stats.Mean)
	fmt.Fprintf(w, "p50   : %s\n", stats.P50)
	fmt.Fprintf(w, "p90   : %s\n", stats.P90)
	fmt.Fprintf(w, "p99   : %s\n", stats.P99)
	fmt.Fprintf(w, "stddev: %s\n", stats.Stddev)
}

type latencyJSON struct {
	Count  int    `json:"count"`
	MinMs  int64  `json:"min_ms"`
	MaxMs  int64  `json:"max_ms"`
	MeanMs int64  `json:"mean_ms"`
	P50Ms  int64  `json:"p50_ms"`
	P90Ms  int64  `json:"p90_ms"`
	P99Ms  int64  `json:"p99_ms"`
	Stddev int64  `json:"stddev_ms"`
}

// LatencyJSON serialises LatencyStats as a JSON string.
func LatencyJSON(stats LatencyStats) (string, error) {
	v := latencyJSON{
		Count:  stats.Count,
		MinMs:  stats.Min.Milliseconds(),
		MaxMs:  stats.Max.Milliseconds(),
		MeanMs: stats.Mean.Milliseconds(),
		P50Ms:  stats.P50.Milliseconds(),
		P90Ms:  stats.P90.Milliseconds(),
		P99Ms:  stats.P99.Milliseconds(),
		Stddev: stats.Stddev.Milliseconds(),
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
