package parser

import (
	"time"
)

// Session represents a group of consecutive log lines within a gap threshold.
type Session struct {
	Lines     []Line
	Start     *time.Time
	End       *time.Time
	Duration  time.Duration
}

// SessionSplitter groups lines into sessions based on inactivity gaps.
type SessionSplitter struct {
	gapThreshold time.Duration
	sessions     []Session
	current      []Line
}

// NewSessionSplitter creates a SessionSplitter that starts a new session
// whenever consecutive lines are separated by more than gapThreshold.
func NewSessionSplitter(gap time.Duration) *SessionSplitter {
	return &SessionSplitter{gapThreshold: gap}
}

// Add pushes a line into the splitter, starting a new session if needed.
func (s *SessionSplitter) Add(l Line) {
	if len(s.current) == 0 {
		s.current = append(s.current, l)
		return
	}
	last := s.current[len(s.current)-1]
	if last.Timestamp != nil && l.Timestamp != nil {
		if l.Timestamp.Sub(*last.Timestamp) > s.gapThreshold {
			s.flush()
		}
	}
	s.current = append(s.current, l)
}

// Flush finalises any open session and returns all sessions.
func (s *SessionSplitter) Flush() []Session {
	s.flush()
	return s.sessions
}

func (s *SessionSplitter) flush() {
	if len(s.current) == 0 {
		return
	}
	sess := Session{Lines: s.current}
	for _, l := range s.current {
		if l.Timestamp != nil {
			if sess.Start == nil || l.Timestamp.Before(*sess.Start) {
				t := *l.Timestamp
				sess.Start = &t
			}
			if sess.End == nil || l.Timestamp.After(*sess.End) {
				t := *l.Timestamp
				sess.End = &t
			}
		}
	}
	if sess.Start != nil && sess.End != nil {
		sess.Duration = sess.End.Sub(*sess.Start)
	}
	s.sessions = append(s.sessions, sess)
	s.current = nil
}

// SplitSessions is a convenience wrapper that splits a slice of lines into sessions.
func SplitSessions(lines []Line, gap time.Duration) []Session {
	sp := NewSessionSplitter(gap)
	for _, l := range lines {
		sp.Add(l)
	}
	return sp.Flush()
}
