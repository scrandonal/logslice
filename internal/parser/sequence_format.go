package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PrintSequences writes a human-readable summary of sequence matches to stdout.
func PrintSequences(matches []SequenceMatch) {
	printSequencesTo(os.Stdout, matches)
}

func printSequencesTo(w io.Writer, matches []SequenceMatch) {
	if len(matches) == 0 {
		fmt.Fprintln(w, "no sequences found")
		return
	}
	for i, m := range matches {
		fmt.Fprintf(w, "--- sequence %d ---\n", i+1)
		if m.Start != nil && m.End != nil {
			fmt.Fprintf(w, "  start:   %s\n", m.Start.Format("2006-01-02T15:04:05"))
			fmt.Fprintf(w, "  end:     %s\n", m.End.Format("2006-01-02T15:04:05"))
		}
		if m.Elapsed != nil {
			fmt.Fprintf(w, "  elapsed: %s\n", m.Elapsed.String())
		}
		for s, l := range m.Steps {
			fmt.Fprintf(w, "  step %d: %s\n", s+1, l.Raw)
		}
	}
}

type sequenceMatchJSON struct {
	Start   *string  `json:"start"`
	End     *string  `json:"end"`
	Elapsed *string  `json:"elapsed"`
	Steps   []string `json:"steps"`
}

// SequencesJSON serialises sequence matches as a JSON array.
func SequencesJSON(matches []SequenceMatch) (string, error) {
	out := make([]sequenceMatchJSON, len(matches))
	for i, m := range matches {
		var start, end, elapsed *string
		if m.Start != nil {
			s := m.Start.Format("2006-01-02T15:04:05Z07:00")
			start = &s
		}
		if m.End != nil {
			e := m.End.Format("2006-01-02T15:04:05Z07:00")
			end = &e
		}
		if m.Elapsed != nil {
			e := m.Elapsed.String()
			elapsed = &e
		}
		steps := make([]string, len(m.Steps))
		for j, l := range m.Steps {
			steps[j] = l.Raw
		}
		out[i] = sequenceMatchJSON{Start: start, End: end, Elapsed: elapsed, Steps: steps}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
