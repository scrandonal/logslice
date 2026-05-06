package parser

import "time"

// Chunk represents a contiguous slice of log lines sharing the same time bucket.
type Chunk struct {
	Start *time.Time
	End   *time.Time
	Lines []string
}

// ChunkSplitter groups filtered log lines into time-based chunks of a fixed
// duration (bucket size). Lines without a parseable timestamp are appended to
// the current open chunk.
type ChunkSplitter struct {
	bucketSize time.Duration
	current    *Chunk
	chunks     []*Chunk
}

// NewChunkSplitter creates a ChunkSplitter that buckets lines by bucketSize.
func NewChunkSplitter(bucketSize time.Duration) *ChunkSplitter {
	return &ChunkSplitter{bucketSize: bucketSize}
}

// Add appends a raw log line to the appropriate chunk.
func (cs *ChunkSplitter) Add(line string) {
	ts := ParseTimestamp(line)

	if ts == nil {
		if cs.current != nil {
			cs.current.Lines = append(cs.current.Lines, line)
		}
		return
	}

	if cs.current == nil || cs.bucketStart(ts) != cs.bucketStart(cs.current.Start) {
		if cs.current != nil {
			cs.chunks = append(cs.chunks, cs.current)
		}
		cs.current = &Chunk{Start: ts, End: ts, Lines: []string{line}}
		return
	}

	cs.current.End = ts
	cs.current.Lines = append(cs.current.Lines, line)
}

// Flush finalises and returns all accumulated chunks.
func (cs *ChunkSplitter) Flush() []*Chunk {
	if cs.current != nil {
		cs.chunks = append(cs.chunks, cs.current)
		cs.current = nil
	}
	out := cs.chunks
	cs.chunks = nil
	return out
}

// bucketStart truncates a timestamp to the nearest bucket boundary.
func (cs *ChunkSplitter) bucketStart(ts *time.Time) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.Truncate(cs.bucketSize)
}
