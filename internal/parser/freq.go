package parser

import (
	"sort"
	"strings"
)

// FreqEntry holds a token and its occurrence count.
type FreqEntry struct {
	Token string
	Count int
}

// FreqResult holds the top-N frequency analysis results.
type FreqResult struct {
	Field   string
	Entries []FreqEntry
	Total   int
}

// FreqOptions configures frequency analysis.
type FreqOptions struct {
	// Field name to extract (e.g. "level", "msg"). Empty means full raw line.
	Field string
	// TopN limits results to the N most frequent tokens. 0 means all.
	TopN int
	// CaseInsensitive normalises tokens to lowercase before counting.
	CaseInsensitive bool
}

// CalcFreq counts token occurrences across lines and returns a FreqResult.
func CalcFreq(lines []Line, opts FreqOptions) FreqResult {
	counts := make(map[string]int, 64)

	for i := range lines {
		var token string
		if opts.Field == "" {
			token = lines[i].Raw
		} else {
			v, ok := extractField(lines[i].Raw, opts.Field)
			if !ok {
				continue
			}
			token = v
		}

		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		if opts.CaseInsensitive {
			token = strings.ToLower(token)
		}
		counts[token]++
	}

	entries := make([]FreqEntry, 0, len(counts))
	total := 0
	for tok, cnt := range counts {
		entries = append(entries, FreqEntry{Token: tok, Count: cnt})
		total += cnt
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Token < entries[j].Token
	})

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}

	return FreqResult{
		Field:   opts.Field,
		Entries: entries,
		Total:   total,
	}
}
