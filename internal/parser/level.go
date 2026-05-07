package parser

import (
	"strings"
)

// LogLevel represents a severity level parsed from a log line.
type LogLevel int

const (
	LevelUnknown LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LevelFilter filters log lines by minimum severity level.
type LevelFilter struct {
	min LogLevel
}

// NewLevelFilter creates a LevelFilter that keeps lines at or above min.
func NewLevelFilter(min LogLevel) *LevelFilter {
	return &LevelFilter{min: min}
}

// ParseLevel extracts the log level from a raw log line.
// It looks for bracketed tokens like [INFO] or bare tokens like ERROR.
func ParseLevel(line string) LogLevel {
	upper := strings.ToUpper(line)
	for _, candidate := range []struct {
		token string
		level LogLevel
	}{
		{"FATAL", LevelFatal},
		{"ERROR", LevelError},
		{"WARN", LevelWarn},
		{"INFO", LevelInfo},
		{"DEBUG", LevelDebug},
	} {
		if strings.Contains(upper, candidate.token) {
			return candidate.level
		}
	}
	return LevelUnknown
}

// Match returns true if the line's level meets the minimum threshold.
func (f *LevelFilter) Match(line string) bool {
	lvl := ParseLevel(line)
	if lvl == LevelUnknown {
		return true // pass through lines with no detectable level
	}
	return lvl >= f.min
}

// FilterByLevel returns only lines that meet the minimum level.
func FilterByLevel(lines []Line, min LogLevel) []Line {
	f := NewLevelFilter(min)
	out := make([]Line, 0, len(lines))
	for _, l := range lines {
		if f.Match(l.Raw) {
			out = append(out, l)
		}
	}
	return out
}
