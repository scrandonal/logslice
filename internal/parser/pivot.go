package parser

import (
	"sort"
	"time"
)

// PivotEntry holds aggregated counts for a single field value across time buckets.
type PivotEntry struct {
	Value   string
	Buckets map[string]int
	Total   int
}

// PivotResult is the full pivot table result.
type PivotResult struct {
	Field   string
	Buckets []string // ordered bucket labels
	Rows    []*PivotEntry
}

// PivotOptions controls how the pivot table is built.
type PivotOptions struct {
	Field      string
	BucketSize time.Duration
	MaxValues  int // 0 = unlimited
}

// NewPivot builds a pivot table grouping log lines by a field value and time bucket.
func NewPivot(lines []Line, opts PivotOptions) *PivotResult {
	if opts.BucketSize <= 0 {
		opts.BucketSize = time.Minute
	}

	bucketSet := map[string]struct{}{}
	valueMap := map[string]*PivotEntry{}

	for _, l := range lines {
		val := extractField(l.Raw, opts.Field)
		if val == "" {
			val = "(none)"
		}

		bucket := "(no time)"
		if l.Timestamp != nil {
			t := l.Timestamp.Truncate(opts.BucketSize)
			bucket = t.UTC().Format("2006-01-02T15:04:05Z")
		}

		bucketSet[bucket] = struct{}{}

		entry, ok := valueMap[val]
		if !ok {
			entry = &PivotEntry{Value: val, Buckets: map[string]int{}}
			valueMap[val] = entry
		}
		entry.Buckets[bucket]++
		entry.Total++
	}

	buckets := make([]string, 0, len(bucketSet))
	for b := range bucketSet {
		buckets = append(buckets, b)
	}
	sort.Strings(buckets)

	rows := make([]*PivotEntry, 0, len(valueMap))
	for _, e := range valueMap {
		rows = append(rows, e)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Total != rows[j].Total {
			return rows[i].Total > rows[j].Total
		}
		return rows[i].Value < rows[j].Value
	})

	if opts.MaxValues > 0 && len(rows) > opts.MaxValues {
		rows = rows[:opts.MaxValues]
	}

	return &PivotResult{
		Field:   opts.Field,
		Buckets: buckets,
		Rows:    rows,
	}
}
