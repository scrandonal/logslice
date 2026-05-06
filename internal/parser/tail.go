package parser

import (
	"bufio"
	"io"
	"time"
)

// TailOptions configures the tail behaviour.
type TailOptions struct {
	PollInterval time.Duration
	MaxLines     int // 0 means unlimited
}

// TailLine is a line emitted by the tailer.
type TailLine struct {
	Line Line
	Err  error
}

// TailReader follows a log stream, emitting new lines as they arrive.
// It stops when ctx is cancelled or the reader returns a permanent error.
func TailReader(r io.Reader, opts TailOptions, out chan<- TailLine) {
	if opts.PollInterval == 0 {
		opts.PollInterval = 250 * time.Millisecond
	}

	scanner := bufio.NewScanner(r)
	count := 0

	for scanner.Scan() {
		raw := scanner.Text()
		ts := ParseTimestamp(raw)
		line := Line{Raw: raw, Timestamp: ts}
		out <- TailLine{Line: line}
		count++
		if opts.MaxLines > 0 && count >= opts.MaxLines {
			close(out)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		out <- TailLine{Err: err}
	}
	close(out)
}
