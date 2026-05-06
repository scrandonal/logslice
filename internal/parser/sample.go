package parser

import "time"

// SampleConfig controls how lines are sampled from a log stream.
type SampleConfig struct {
	// Every N lines, keep one (1 = keep all).
	Rate int
	// If non-zero, only keep lines within this time bucket width.
	Bucket time.Duration
}

// Sampler filters a slice of Lines according to a SampleConfig.
type Sampler struct {
	cfg     SampleConfig
	counter int
}

// NewSampler creates a Sampler with the given config.
// Rate < 1 is treated as 1 (keep every line).
func NewSampler(cfg SampleConfig) *Sampler {
	if cfg.Rate < 1 {
		cfg.Rate = 1
	}
	return &Sampler{cfg: cfg}
}

// Sample returns the subset of lines that pass the sampling policy.
func (s *Sampler) Sample(lines []Line) []Line {
	out := make([]Line, 0, len(lines))
	var bucketStart *time.Time

	for i := range lines {
		l := lines[i]

		// Rate-based sampling.
		s.counter++
		if s.counter%s.cfg.Rate != 0 {
			continue
		}

		// Bucket-based sampling: keep first line of each time bucket.
		if s.cfg.Bucket > 0 && l.Timestamp != nil {
			bucketKey := l.Timestamp.Truncate(s.cfg.Bucket)
			if bucketStart != nil && bucketKey.Equal(*bucketStart) {
				continue
			}
			bucketStart = &bucketKey
		}

		out = append(out, l)
	}
	return out
}

// Reset resets internal counters so the sampler can be reused.
func (s *Sampler) Reset() {
	s.counter = 0
}
