package parser

import (
	"fmt"
	"io"
	"time"
)

// TailPipelineOptions configures a live tail + window pipeline.
type TailPipelineOptions struct {
	WindowSize   int
	From         *time.Time
	To           *time.Time
	Format       string // "raw" or "json"
	MaxLines     int
	PollInterval time.Duration
}

// RunTailPipeline reads lines from r, maintains a rolling window, and writes
// matching lines to w using the specified format.
func RunTailPipeline(r io.Reader, w io.Writer, opts TailPipelineOptions) error {
	if opts.WindowSize <= 0 {
		opts.WindowSize = 100
	}

	ch := make(chan TailLine, 64)
	go TailReader(r, TailOptions{
		PollInterval: opts.PollInterval,
		MaxLines:     opts.MaxLines,
	}, ch)

	win := NewWindow(opts.WindowSize)
	fmt_ := NewFormatter(opts.Format)
	var count int

	for tl := range ch {
		if tl.Err != nil {
			return fmt.Errorf("tail error: %w", tl.Err)
		}
		win.Push(tl.Line)

		if opts.From != nil && opts.To != nil {
			// Only emit lines that fall within the requested range.
			if tl.Line.Timestamp != nil {
				if tl.Line.Timestamp.Before(*opts.From) || tl.Line.Timestamp.After(*opts.To) {
					continue
				}
			}
		}

		if _, err := fmt.Fprintln(w, fmt_.Format(tl.Line)); err != nil {
			return err
		}
		count++
	}
	return nil
}
