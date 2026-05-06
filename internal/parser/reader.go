package parser

import (
	"bufio"
	"io"
	"time"
)

// ReadOptions configures the behavior of ReadLines.
type ReadOptions struct {
	From   *time.Time
	To     *time.Time
	Format string // "raw" or "json"
}

// ReadResult holds the output of a ReadLines call.
type ReadResult struct {
	Lines     []string
	Collector *Collector
}

// ReadLines reads log lines from r, applies timestamp filtering based on opts,
// and returns the matching lines along with collected statistics.
func ReadLines(r io.Reader, opts ReadOptions) (*ReadResult, error) {
	scanner := NewScanner(bufio.NewReader(r))
	filter, err := NewFilter(opts.From, opts.To)
	if err != nil {
		return nil, err
	}

	format := opts.Format
	if format == "" {
		format = "raw"
	}
	fmt := NewFormatter(format)
	collector := NewCollector()

	var lines []string
	for scanner.Scan() {
		line := scanner.Line()
		collector.Count(line)

		if !filter.Match(line) {
			continue
		}

		collector.CountMatch(line)
		formatted, fErr := fmt.Format(line)
		if fErr != nil {
			return nil, fErr
		}
		lines = append(lines, formatted)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ReadResult{
		Lines:     lines,
		Collector: collector,
	}, nil
}
