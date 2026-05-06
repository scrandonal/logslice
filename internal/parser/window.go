package parser

import "time"

// Window holds a rolling buffer of the most recent N lines.
type Window struct {
	buf  []Line
	cap  int
	head int
	size int
}

// NewWindow creates a Window that retains at most n lines.
func NewWindow(n int) *Window {
	if n <= 0 {
		n = 1
	}
	return &Window{buf: make([]Line, n), cap: n}
}

// Push adds a line to the window, evicting the oldest if full.
func (w *Window) Push(l Line) {
	w.buf[w.head] = l
	w.head = (w.head + 1) % w.cap
	if w.size < w.cap {
		w.size++
	}
}

// Lines returns the buffered lines in chronological order.
func (w *Window) Lines() []Line {
	out := make([]Line, w.size)
	start := (w.head - w.size + w.cap) % w.cap
	for i := 0; i < w.size; i++ {
		out[i] = w.buf[(start+i)%w.cap]
	}
	return out
}

// Len returns the current number of buffered lines.
func (w *Window) Len() int { return w.size }

// InRange returns lines whose timestamp falls within [from, to].
// Lines without a timestamp are always included.
func (w *Window) InRange(from, to time.Time) []Line {
	all := w.Lines()
	out := make([]Line, 0, len(all))
	for _, l := range all {
		if l.Timestamp == nil || (!l.Timestamp.Before(from) && !l.Timestamp.After(to)) {
			out = append(out, l)
		}
	}
	return out
}
