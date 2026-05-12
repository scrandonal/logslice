package parser

import (
	"regexp"
	"time"
)

// SequenceMatch holds a group of lines that matched a multi-step sequence.
type SequenceMatch struct {
	Steps  []Line
	Start  *time.Time
	End    *time.Time
	Elapsed *time.Duration
}

// SequenceDetector finds ordered multi-pattern sequences within a log stream.
type SequenceDetector struct {
	patterns []*regexp.Regexp
	window   time.Duration
}

// NewSequenceDetector creates a detector for the given ordered patterns and
// optional time window (0 means no window constraint).
func NewSequenceDetector(patterns []string, window time.Duration) (*SequenceDetector, error) {
	regs := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		regs[i] = r
	}
	return &SequenceDetector{patterns: regs, window: window}, nil
}

// Detect scans lines for all non-overlapping occurrences of the sequence.
func (sd *SequenceDetector) Detect(lines []Line) []SequenceMatch {
	if len(sd.patterns) == 0 || len(lines) == 0 {
		return nil
	}
	var matches []SequenceMatch
	i := 0
	for i < len(lines) {
		if !sd.patterns[0].MatchString(lines[i].Raw) {
			i++
			continue
		}
		steps := make([]Line, 0, len(sd.patterns))
		steps = append(steps, lines[i])
		step := 1
		for j := i + 1; j < len(lines) && step < len(sd.patterns); j++ {
			if sd.window > 0 && lines[i].Timestamp != nil && lines[j].Timestamp != nil {
				if lines[j].Timestamp.Sub(*lines[i].Timestamp) > sd.window {
					break
				}
			}
			if sd.patterns[step].MatchString(lines[j].Raw) {
				steps = append(steps, lines[j])
				step++
			}
		}
		if step == len(sd.patterns) {
			m := SequenceMatch{Steps: steps}
			m.Start = steps[0].Timestamp
			m.End = steps[len(steps)-1].Timestamp
			if m.Start != nil && m.End != nil {
				e := m.End.Sub(*m.Start)
				m.Elapsed = &e
			}
			matches = append(matches, m)
			i = i + 1
		} else {
			i++
		}
	}
	return matches
}

// DetectSequences is a convenience wrapper.
func DetectSequences(lines []Line, patterns []string, window time.Duration) ([]SequenceMatch, error) {
	sd, err := NewSequenceDetector(patterns, window)
	if err != nil {
		return nil, err
	}
	return sd.Detect(lines), nil
}
