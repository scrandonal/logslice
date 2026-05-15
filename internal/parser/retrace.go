package parser

import (
	"regexp"
	"strings"
)

// RetraceEntry holds a single stack-trace group: the trigger line plus
// any following continuation lines (e.g. Java/Python/Go stack frames).
type RetraceEntry struct {
	Trigger  Line
	Frames   []Line
}

// RetraceOptions configures the retrace extractor.
type RetraceOptions struct {
	// FramePattern is a regex that identifies a stack-frame continuation line.
	// Defaults to a pattern covering Java, Python, and Go style frames.
	FramePattern string
	// MaxFrames caps how many frames are collected per entry (0 = unlimited).
	MaxFrames int
}

var defaultFrameRe = regexp.MustCompile(
	`^\s+(at |File "|goroutine |\.\.\.|\w+\.go:|panic)`,
)

// RetraceLines scans lines and groups stack traces together.
// A new entry is started whenever a non-frame line is encountered after
// frames have already been collected, or when a new trigger is seen.
func RetraceLines(lines []Line, opts RetraceOptions) []RetraceEntry {
	frameRe := defaultFrameRe
	if opts.FramePattern != "" {
		var err error
		frameRe, err = regexp.Compile(opts.FramePattern)
		if err != nil {
			frameRe = defaultFrameRe
		}
	}

	var entries []RetraceEntry
	var current *RetraceEntry

	for _, l := range lines {
		isFrame := frameRe.MatchString(l.Raw) || strings.HasPrefix(l.Raw, "\t")

		if isFrame && current != nil {
			if opts.MaxFrames == 0 || len(current.Frames) < opts.MaxFrames {
				current.Frames = append(current.Frames, l)
			}
			continue
		}

		// Non-frame line: save current entry if it has frames, start new.
		if current != nil {
			entries = append(entries, *current)
		}
		current = &RetraceEntry{Trigger: l}
	}

	if current != nil {
		entries = append(entries, *current)
	}

	return entries
}
