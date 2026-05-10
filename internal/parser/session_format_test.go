package parser

import (
	"strings"
	"testing"
	"time"
)

func buildSessions() []Session {
	t1 := time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 3, 10, 9, 0, 45, 0, time.UTC)
	return []Session{
		{
			Lines:    []Line{{Raw: "foo"}, {Raw: "bar"}},
			Start:    &t1,
			End:      &t2,
			Duration: t2.Sub(t1),
		},
	}
}

func TestPrintSessionsEmpty(t *testing.T) {
	var sb strings.Builder
	PrintSessions(&sb, nil)
	if !strings.Contains(sb.String(), "no sessions") {
		t.Errorf("expected 'no sessions' message, got: %s", sb.String())
	}
}

func TestPrintSessionsSingle(t *testing.T) {
	var sb strings.Builder
	PrintSessions(&sb, buildSessions())
	out := sb.String()
	if !strings.Contains(out, "session 1") {
		t.Errorf("missing session label: %s", out)
	}
	if !strings.Contains(out, "lines=2") {
		t.Errorf("missing line count: %s", out)
	}
	if !strings.Contains(out, "2024-03-10T09:00:00") {
		t.Errorf("missing start timestamp: %s", out)
	}
	if !strings.Contains(out, "45s") {
		t.Errorf("missing duration: %s", out)
	}
}

func TestSessionsJSONFields(t *testing.T) {
	out := SessionsJSON(buildSessions())
	for _, want := range []string{`"session":1`, `"lines":2`, `"duration_ms":45000`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in: %s", want, out)
		}
	}
}

func TestSessionsJSONNullTimestamps(t *testing.T) {
	sessions := []Session{
		{Lines: []Line{{Raw: "only"}}, Start: nil, End: nil},
	}
	out := SessionsJSON(sessions)
	if !strings.Contains(out, `"start":null`) {
		t.Errorf("expected null start: %s", out)
	}
	if !strings.Contains(out, `"end":null`) {
		t.Errorf("expected null end: %s", out)
	}
}
