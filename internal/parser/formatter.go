package parser

import (
	"fmt"
	"io"
	"time"
)

// OutputFormat controls how matched log lines are written.
type OutputFormat int

const (
	// FormatRaw writes lines exactly as they appear in the source.
	FormatRaw OutputFormat = iota
	// FormatJSON wraps each line in a simple JSON envelope.
	FormatJSON
)

// Formatter writes log lines to an io.Writer in the requested format.
type Formatter struct {
	w      io.Writer
	format OutputFormat
}

// NewFormatter creates a Formatter that writes to w using the given format.
func NewFormatter(w io.Writer, format OutputFormat) *Formatter {
	return &Formatter{w: w, format: format}
}

// WriteLine writes a single log line, optionally with its parsed timestamp.
// ts may be zero if the timestamp is unavailable or irrelevant.
func (f *Formatter) WriteLine(line []byte, ts time.Time) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(line, ts)
	default:
		return f.writeRaw(line)
	}
}

func (f *Formatter) writeRaw(line []byte) error {
	_, err := fmt.Fprintf(f.w, "%s\n", line)
	return err
}

func (f *Formatter) writeJSON(line []byte, ts time.Time) error {
	var tsField string
	if !ts.IsZero() {
		tsField = fmt.Sprintf(`"timestamp":%q,`, ts.UTC().Format(time.RFC3339Nano))
	} else {
		tsField = `"timestamp":null,`
	}
	// Escape the raw line as a JSON string value.
	escaped := jsonEscapeString(line)
	_, err := fmt.Fprintf(f.w, `{%s"line":%s}`+"\n", tsField, escaped)
	return err
}

// jsonEscapeString returns a JSON-quoted representation of b.
func jsonEscapeString(b []byte) string {
	out := make([]byte, 0, len(b)+2)
	out = append(out, '"')
	for _, c := range b {
		switch c {
		case '\\', '"':
			out = append(out, '\\', c)
		case '\n':
			out = append(out, '\\', 'n')
		case '\r':
			out = append(out, '\\', 'r')
		case '\t':
			out = append(out, '\\', 't')
		default:
			out = append(out, c)
		}
	}
	out = append(out, '"')
	return string(out)
}
