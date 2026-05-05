package parser

import "time"

// Stats holds counters collected during a log slice operation.
type Stats struct {
	LinesScanned  int64
	LinesMatched  int64
	LinesSkipped  int64
	ParseErrors   int64
	FirstMatch    *time.Time
	LastMatch     *time.Time
	Elapsed       time.Duration
}

// Collector accumulates statistics while processing log lines.
type Collector struct {
	stats   Stats
	started time.Time
}

// NewCollector creates a new Collector and starts the elapsed timer.
func NewCollector() *Collector {
	return &Collector{started: time.Now()}
}

// RecordScanned increments the scanned line counter.
func (c *Collector) RecordScanned() {
	c.stats.LinesScanned++
}

// RecordMatch increments the matched counter and tracks first/last timestamps.
func (c *Collector) RecordMatch(ts *time.Time) {
	c.stats.LinesMatched++
	if ts != nil {
		if c.stats.FirstMatch == nil {
			t := *ts
			c.stats.FirstMatch = &t
		}
		t := *ts
		c.stats.LastMatch = &t
	}
}

// RecordSkipped increments the skipped counter.
func (c *Collector) RecordSkipped() {
	c.stats.LinesSkipped++
}

// RecordParseError increments the parse-error counter.
func (c *Collector) RecordParseError() {
	c.stats.ParseErrors++
}

// Finalize stops the elapsed timer and returns the final Stats snapshot.
func (c *Collector) Finalize() Stats {
	c.stats.Elapsed = time.Since(c.started)
	return c.stats
}
