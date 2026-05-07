package parser

import (
	"strings"
)

// HighlightMode controls how matches are marked in output.
type HighlightMode int

const (
	HighlightNone  HighlightMode = iota
	HighlightANSI                // wrap matches with ANSI color codes
	HighlightBracket             // wrap matches with [[ ]]
)

// Highlighter scans log lines and marks substrings that match any of the
// provided keywords.
type Highlighter struct {
	keywords []string
	mode     HighlightMode
}

// NewHighlighter creates a Highlighter for the given keywords and mode.
// Keywords are matched case-insensitively.
func NewHighlighter(keywords []string, mode HighlightMode) *Highlighter {
	lower := make([]string, len(keywords))
	for i, k := range keywords {
		lower[i] = strings.ToLower(k)
	}
	return &Highlighter{keywords: lower, mode: mode}
}

// Highlight returns a copy of line with all keyword occurrences marked.
// If mode is HighlightNone or no keywords match, the original line is returned.
func (h *Highlighter) Highlight(line string) string {
	if h.mode == HighlightNone || len(h.keywords) == 0 {
		return line
	}
	result := line
	lower := strings.ToLower(line)
	for _, kw := range h.keywords {
		result = replaceInsensitive(result, lower, kw, h.mode)
		lower = strings.ToLower(result)
	}
	return result
}

// replaceInsensitive replaces all occurrences of kw (found via lowerSrc) in
// src with the marked version according to mode.
func replaceInsensitive(src, lowerSrc, kw string, mode HighlightMode) string {
	var b strings.Builder
	offset := 0
	for {
		idx := strings.Index(lowerSrc[offset:], kw)
		if idx < 0 {
			b.WriteString(src[offset:])
			break
		}
		abs := offset + idx
		b.WriteString(src[offset:abs])
		original := src[abs : abs+len(kw)]
		switch mode {
		case HighlightANSI:
			b.WriteString("\x1b[1;33m")
			b.WriteString(original)
			b.WriteString("\x1b[0m")
		case HighlightBracket:
			b.WriteString("[[")
			b.WriteString(original)
			b.WriteString("]]")
		}
		offset = abs + len(kw)
	}
	return b.String()
}
