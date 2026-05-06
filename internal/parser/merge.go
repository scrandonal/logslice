package parser

import (
	"sort"
	"time"
)

// MergedLine represents a log line with its source index for multi-file merging.
type MergedLine struct {
	Line
	Source int
}

// MergeLines merges multiple sorted slices of Lines into a single
// chronologically ordered slice. Lines without timestamps are appended last.
func MergeLines(sources [][]Line) []MergedLine {
	var withTS []MergedLine
	var noTS []MergedLine

	for srcIdx, lines := range sources {
		for _, l := range lines {
			ml := MergedLine{Line: l, Source: srcIdx}
			if l.Timestamp != nil {
				withTS = append(withTS, ml)
			} else {
				noTS = append(noTS, ml)
			}
		}
	}

	sort.SliceStable(withTS, func(i, j int) bool {
		return withTS[i].Timestamp.Before(*withTS[j].Timestamp)
	})

	return append(withTS, noTS...)
}

// MergeTimeRange returns only merged lines whose timestamps fall within
// [from, to]. Lines with nil timestamps are always excluded.
func MergeTimeRange(sources [][]Line, from, to time.Time) []MergedLine {
	all := MergeLines(sources)
	var out []MergedLine
	for _, ml := range all {
		if ml.Timestamp == nil {
			continue
		}
		if !ml.Timestamp.Before(from) && !ml.Timestamp.After(to) {
			out = append(out, ml)
		}
	}
	return out
}
