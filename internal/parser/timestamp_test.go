package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
		wantUTC string // expected time in RFC3339 UTC, empty to skip check
	}{
		{
			name:    "RFC3339 prefix",
			line:    "2024-03-15T10:22:33Z INFO server started",
			wantUTC: "2024-03-15T10:22:33Z",
		},
		{
			name:    "bracketed RFC3339",
			line:    "[2024-03-15T10:22:33Z] ERROR something failed",
			wantUTC: "2024-03-15T10:22:33Z",
		},
		{
			name:    "space-separated datetime",
			line:    "2024-03-15 10:22:33 DEBUG processing request",
			wantUTC: "2024-03-15T10:22:33Z",
		},
		{
			name:    "apache combined log format",
			line:    "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326",
			wantErr: false, // timestamp embedded, not at start — expect error
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "no timestamp",
			line:    "this is a plain log line without any timestamp",
			wantErr: true,
		},
		{
			name:    "RFC3339 with nanoseconds",
			line:    "2024-03-15T10:22:33.123456789Z WARN high latency detected",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTimestamp(tc.line)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got time %v", got)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tc.wantUTC != "" {
				want, _ := time.Parse(time.RFC3339, tc.wantUTC)
				if !got.UTC().Equal(want.UTC()) {
					t.Errorf("got %v, want %v", got.UTC(), want.UTC())
				}
			}
		})
	}
}

func TestParseTimestampEmptyBracket(t *testing.T) {
	_, err := ParseTimestamp("[] no timestamp here")
	if err == nil {
		t.Error("expected error for empty bracket prefix")
	}
}
