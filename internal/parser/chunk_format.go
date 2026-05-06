package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ChunkSummary is a JSON-serialisable summary of a single Chunk.
type ChunkSummary struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	LineCount int   `json:"line_count"`
}

const tsLayout = "2006-01-02T15:04:05Z"

// PrintChunkSummaries writes a human-readable table of chunk summaries to w.
func PrintChunkSummaries(w io.Writer, chunks []*Chunk) {
	fmt.Fprintf(w, "%-25s %-25s %s\n", "start", "end", "lines")
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, c := range chunks {
		start, end := formatChunkTimes(c)
		fmt.Fprintf(w, "%-25s %-25s %d\n", start, end, len(c.Lines))
	}
}

// ChunksJSON serialises chunk summaries as a JSON array string.
func ChunksJSON(chunks []*Chunk) (string, error) {
	summaries := make([]ChunkSummary, 0, len(chunks))
	for _, c := range chunks {
		start, end := formatChunkTimes(c)
		summaries = append(summaries, ChunkSummary{
			Start:    start,
			End:      end,
			LineCount: len(c.Lines),
		})
	}
	b, err := json.Marshal(summaries)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func formatChunkTimes(c *Chunk) (start, end string) {
	if c.Start != nil {
		start = c.Start.UTC().Format(tsLayout)
	}
	if c.End != nil {
		end = c.End.UTC().Format(tsLayout)
	}
	return
}
