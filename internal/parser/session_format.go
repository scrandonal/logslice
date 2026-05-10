package parser

import (
	"fmt"
	"io"
	"strings"
)

// PrintSessions writes a human-readable session summary to w.
func PrintSessions(w io.Writer, sessions []Session) {
	if len(sessions) == 0 {
		fmt.Fprintln(w, "no sessions found")
		return
	}
	for i, s := range sessions {
		start := "n/a"
		end := "n/a"
		if s.Start != nil {
			start = s.Start.Format("2006-01-02T15:04:05")
		}
		if s.End != nil {
			end = s.End.Format("2006-01-02T15:04:05")
		}
		fmt.Fprintf(w, "session %d: lines=%d start=%s end=%s duration=%s\n",
			i+1, len(s.Lines), start, end, s.Duration)
	}
}

// SessionsJSON returns a JSON array of session summaries.
func SessionsJSON(sessions []Session) string {
	if len(sessions) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, s := range sessions {
		if i > 0 {
			sb.WriteString(",")
		}
		start := "null"
		end := "null"
		if s.Start != nil {
			start = fmt.Sprintf("%q", s.Start.Format("2006-01-02T15:04:05Z07:00"))
		}
		if s.End != nil {
			end = fmt.Sprintf("%q", s.End.Format("2006-01-02T15:04:05Z07:00"))
		}
		fmt.Fprintf(&sb,
			`{"session":%d,"lines":%d,"start":%s,"end":%s,"duration_ms":%d}`,
			i+1, len(s.Lines), start, end, s.Duration.Milliseconds(),
		)
	}
	sb.WriteString("]")
	return sb.String()
}
