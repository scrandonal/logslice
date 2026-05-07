package parser

import (
	"fmt"
	"io"
	"strings"
)

// PrintLevelSummary writes a breakdown of line counts by log level to w.
func PrintLevelSummary(w io.Writer, lines []Line) {
	counts := countByLevel(lines)
	fmt.Fprintf(w, "Level summary (%d lines):\n", len(lines))
	for _, lvl := range []LogLevel{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelUnknown} {
		if n, ok := counts[lvl]; ok && n > 0 {
			fmt.Fprintf(w, "  %-8s %d\n", lvl.String(), n)
		}
	}
}

// LevelSummaryJSON returns a JSON object mapping level names to counts.
func LevelSummaryJSON(lines []Line) string {
	counts := countByLevel(lines)
	levels := []LogLevel{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelUnknown}
	parts := make([]string, 0, len(levels))
	for _, lvl := range levels {
		if n, ok := counts[lvl]; ok && n > 0 {
			parts = append(parts, fmt.Sprintf(`%q:%d`, lvl.String(), n))
		}
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func countByLevel(lines []Line) map[LogLevel]int {
	counts := make(map[LogLevel]int)
	for _, l := range lines {
		lvl := ParseLevel(l.Raw)
		counts[lvl]++
	}
	return counts
}
