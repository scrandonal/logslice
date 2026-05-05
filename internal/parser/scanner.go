package parser

import (
	"bufio"
	"io"
	"time"
)

// Entry represents a single parsed log entry with its timestamp and raw line.
type Entry struct {
	Timestamp time.Time
	Line      []byte
}

// Scanner reads log lines from a reader and parses timestamps from each line.
type Scanner struct {
	scanner *bufio.Scanner
	current Entry
	err     error
}

// NewScanner creates a new Scanner wrapping the given reader.
func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &Scanner{scanner: s}
}

// Scan advances to the next log entry. Returns true if an entry is available.
func (s *Scanner) Scan() bool {
	for s.scanner.Scan() {
		line := s.scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		ts, err := ParseTimestamp(line)
		if err != nil {
			// Skip lines that don't have a parseable timestamp.
			continue
		}

		// Copy the line bytes since scanner reuses the buffer.
		copied := make([]byte, len(line))
		copy(copied, line)

		s.current = Entry{
			Timestamp: ts,
			Line:      copied,
		}
		return true
	}

	s.err = s.scanner.Err()
	return false
}

// Entry returns the current log entry.
func (s *Scanner) Entry() Entry {
	return s.current
}

// Err returns any scanning error, excluding io.EOF.
func (s *Scanner) Err() error {
	return s.err
}
