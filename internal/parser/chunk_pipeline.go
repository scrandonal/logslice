package parser

import (
	"bufio"
	"io"
	"time"
)

// ChunkPipelineOptions configures the chunk extraction pipeline.
type ChunkPipelineOptions struct {
	// BucketSize controls how lines are grouped into time chunks.
	BucketSize time.Duration
	// Filter, if non-nil, restricts lines to a time range before chunking.
	Filter *Filter
}

// RunChunkPipeline reads lines from r, optionally filters them by time range,
// and splits them into time-bucketed chunks. It returns all resulting chunks.
func RunChunkPipeline(r io.Reader, opts ChunkPipelineOptions) ([]*Chunk, error) {
	bucketSize := opts.BucketSize
	if bucketSize <= 0 {
		bucketSize = time.Minute
	}

	splitter := NewChunkSplitter(bucketSize)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if opts.Filter != nil {
			ts := ParseTimestamp(line)
			if ts != nil && !opts.Filter.Match(ts) {
				continue
			}
		}
		splitter.Add(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return splitter.Flush(), nil
}
