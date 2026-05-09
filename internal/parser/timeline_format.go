package parser

import (
	"fmt"
	"io"
	"strings"
)

const timelineBarWidth = 40

// PrintTimeline writes an ASCII bar-chart of the timeline to w.
func PrintTimeline(w io.Writer, result TimelineResult) {
	if len(result.Buckets) == 0 {
		fmt.Fprintln(w, "(no timestamped lines)")
		return
	}

	max := 0
	for _, b := range result.Buckets {
		if b.Count > max {
			max = b.Count
		}
	}

	for _, b := range result.Buckets {
		barLen := 0
		if max > 0 {
			barLen = b.Count * timelineBarWidth / max
		}
		bar := strings.Repeat("█", barLen)
		fmt.Fprintf(w, "%-21s │ %-*s %d\n", b.Label, timelineBarWidth, bar, b.Count)
	}
	fmt.Fprintf(w, "total: %d lines, interval: %s\n", result.Total, result.Interval)
}

// TimelineJSON returns a JSON representation of the timeline result.
func TimelineJSON(result TimelineResult) string {
	if len(result.Buckets) == 0 {
		return `{"interval_ns":` + fmt.Sprintf("%d", result.Interval.Nanoseconds()) + `,"total":0,"buckets":[]}`
	}
	var sb strings.Builder
	sb.WriteString(`{"interval_ns":`)
	sb.WriteString(fmt.Sprintf("%d", result.Interval.Nanoseconds()))
	sb.WriteString(`,"total":`)
	sb.WriteString(fmt.Sprintf("%d", result.Total))
	sb.WriteString(`,"buckets":[`)
	for i, b := range result.Buckets {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf(`{"label":%q,"start_ns":%d,"count":%d}`,
			b.Label, b.Start.UnixNano(), b.Count))
	}
	sb.WriteString("]}")  
	return sb.String()
}
