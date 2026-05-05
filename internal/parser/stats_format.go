package parser

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintStats writes a human-readable summary of s to w.
func PrintStats(w io.Writer, s Stats) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Lines scanned:\t%d\n", s.LinesScanned)
	fmt.Fprintf(tw, "Lines matched:\t%d\n", s.LinesMatched)
	fmt.Fprintf(tw, "Lines skipped:\t%d\n", s.LinesSkipped)
	fmt.Fprintf(tw, "Parse errors:\t%d\n", s.ParseErrors)

	if s.FirstMatch != nil {
		fmt.Fprintf(tw, "First match:\t%s\n", s.FirstMatch.Format("2006-01-02 15:04:05 UTC"))
	} else {
		fmt.Fprintf(tw, "First match:\t—\n")
	}

	if s.LastMatch != nil {
		fmt.Fprintf(tw, "Last match:\t%s\n", s.LastMatch.Format("2006-01-02 15:04:05 UTC"))
	} else {
		fmt.Fprintf(tw, "Last match:\t—\n")
	}

	fmt.Fprintf(tw, "Elapsed:\t%s\n", s.Elapsed.Round(1000).String())

	return tw.Flush()
}

// StatsJSON returns a compact JSON representation of s.
func StatsJSON(s Stats) string {
	first, last := "null", "null"
	if s.FirstMatch != nil {
		first = fmt.Sprintf("%q", s.FirstMatch.UTC().Format("2006-01-02T15:04:05Z"))
	}
	if s.LastMatch != nil {
		last = fmt.Sprintf("%q", s.LastMatch.UTC().Format("2006-01-02T15:04:05Z"))
	}
	return fmt.Sprintf(
		`{"lines_scanned":%d,"lines_matched":%d,"lines_skipped":%d,"parse_errors":%d,"first_match":%s,"last_match":%s,"elapsed_ms":%d}`,
		s.LinesScanned,
		s.LinesMatched,
		s.LinesSkipped,
		s.ParseErrors,
		first,
		last,
		s.Elapsed.Milliseconds(),
	)
}
