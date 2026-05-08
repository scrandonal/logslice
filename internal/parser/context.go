package parser

import "strings"

// ContextLine holds a log line along with its surrounding context lines.
type ContextLine struct {
	Line    Line
	Before  []Line
	After   []Line
	Matched bool
}

// ContextExtractor extracts surrounding lines around matched entries.
type ContextExtractor struct {
	before int
	after  int
	buf    []Line
	queue  []ContextLine
	pend   []pendingCtx
}

type pendingCtx struct {
	cl        ContextLine
	remaining int
}

// NewContextExtractor creates a ContextExtractor that includes `before` lines
// before and `after` lines after each matched line.
func NewContextExtractor(before, after int) *ContextExtractor {
	return &ContextExtractor{
		before: before,
		after:  after,
		buf:    make([]Line, 0, before+1),
	}
}

// Push adds a line and a match flag; returns any ContextLines ready to emit.
func (c *ContextExtractor) Push(l Line, matched bool) []ContextLine {
	// Advance pending after-context
	for i := range c.pend {
		c.pend[i].cl.After = append(c.pend[i].cl.After, l)
		c.pend[i].remaining--
	}

	var ready []ContextLine
	for len(c.pend) > 0 && c.pend[0].remaining <= 0 {
		ready = append(ready, c.pend[0].cl)
		c.pend = c.pend[1:]
	}

	if matched {
		before := make([]Line, len(c.buf))
		copy(before, c.buf)
		cl := ContextLine{
			Line:    l,
			Before:  before,
			Matched: true,
		}
		if c.after == 0 {
			ready = append(ready, cl)
		} else {
			c.pend = append(c.pend, pendingCtx{cl: cl, remaining: c.after})
		}
	}

	// Maintain rolling before-buffer
	c.buf = append(c.buf, l)
	if len(c.buf) > c.before {
		c.buf = c.buf[len(c.buf)-c.before:]
	}

	return ready
}

// Flush returns any remaining pending ContextLines.
func (c *ContextExtractor) Flush() []ContextLine {
	var out []ContextLine
	for _, p := range c.pend {
		out = append(out, p.cl)
	}
	c.pend = nil
	return out
}

// ExtractContext runs lines through the extractor, matching those containing needle.
func ExtractContext(lines []Line, needle string, before, after int, caseInsensitive bool) []ContextLine {
	ex := NewContextExtractor(before, after)
	var out []ContextLine
	for _, l := range lines {
		raw := l.Raw
		if caseInsensitive {
			raw = strings.ToLower(raw)
			needle = strings.ToLower(needle)
		}
		out = append(out, ex.Push(l, strings.Contains(raw, needle))...)
	}
	out = append(out, ex.Flush()...)
	return out
}
