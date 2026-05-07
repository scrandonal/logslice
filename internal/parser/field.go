package parser

import (
	"strings"
)

// FieldExtractor extracts named fields from structured log lines.
// It supports key=value and key="value" formats.
type FieldExtractor struct {
	fields []string
}

// NewFieldExtractor creates a FieldExtractor that extracts the given field names.
func NewFieldExtractor(fields []string) *FieldExtractor {
	return &FieldExtractor{fields: fields}
}

// Extract returns a map of field->value for a given log line.
// Fields not found in the line are omitted from the result.
func (fe *FieldExtractor) Extract(line string) map[string]string {
	result := make(map[string]string, len(fe.fields))
	for _, field := range fe.fields {
		if v, ok := extractField(line, field); ok {
			result[field] = v
		}
	}
	return result
}

// ExtractAll returns all key=value pairs found in the line.
func ExtractAll(line string) map[string]string {
	result := make(map[string]string)
	remaining := line
	for {
		eqIdx := strings.IndexByte(remaining, '=')
		if eqIdx < 0 {
			break
		}
		key := lastToken(remaining[:eqIdx])
		if key == "" {
			remaining = remaining[eqIdx+1:]
			continue
		}
		val, advance := readValue(remaining[eqIdx+1:])
		result[key] = val
		remaining = remaining[eqIdx+1+advance:]
	}
	return result
}

func extractField(line, field string) (string, bool) {
	prefix := field + "="
	idx := strings.Index(line, prefix)
	if idx < 0 {
		return "", false
	}
	val, _ := readValue(line[idx+len(prefix):])
	return val, true
}

func readValue(s string) (string, int) {
	if len(s) == 0 {
		return "", 0
	}
	if s[0] == '"' {
		end := strings.IndexByte(s[1:], '"')
		if end < 0 {
			return s[1:], len(s)
		}
		return s[1 : end+1], end + 2
	}
	end := strings.IndexAny(s, " \t\n")
	if end < 0 {
		return s, len(s)
	}
	return s[:end], end
}

func lastToken(s string) string {
	s = strings.TrimRight(s, " \t")
	idx := strings.LastIndexAny(s, " \t[]{},;")
	if idx < 0 {
		return s
	}
	return s[idx+1:]
}
