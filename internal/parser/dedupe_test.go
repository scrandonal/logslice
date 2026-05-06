package parser

import (
	"fmt"
	"testing"
)

func makeDedupeLine(raw string) Line {
	return Line{Raw: raw}
}

func TestDeduplicatorNoDuplicates(t *testing.T) {
	d := NewDeduplicator(10)
	lines := []string{"alpha", "beta", "gamma"}
	for _, s := range lines {
		if d.IsDuplicate(makeDedupeLine(s)) {
			t.Errorf("expected %q to be unique on first encounter", s)
		}
	}
}

func TestDeduplicatorDetectsDuplicates(t *testing.T) {
	d := NewDeduplicator(10)
	l := makeDedupeLine("repeated line")
	if d.IsDuplicate(l) {
		t.Fatal("first occurrence should not be a duplicate")
	}
	if !d.IsDuplicate(l) {
		t.Fatal("second occurrence should be detected as duplicate")
	}
}

func TestDeduplicatorEviction(t *testing.T) {
	cap := 3
	d := NewDeduplicator(cap)

	// Fill to capacity with distinct lines.
	for i := 0; i < cap; i++ {
		d.IsDuplicate(makeDedupeLine(fmt.Sprintf("line-%d", i)))
	}

	// Adding one more should evict line-0.
	d.IsDuplicate(makeDedupeLine("line-new"))

	// line-0 should now be treated as new (evicted from seen set).
	if d.IsDuplicate(makeDedupeLine("line-0")) {
		t.Error("line-0 should have been evicted and treated as new")
	}
}

func TestDedupeLinesHelper(t *testing.T) {
	input := []Line{
		makeDedupeLine("foo"),
		makeDedupeLine("bar"),
		makeDedupeLine("foo"),
		makeDedupeLine("baz"),
		makeDedupeLine("bar"),
	}

	got := DedupeLines(input, 100)
	want := []string{"foo", "bar", "baz"}

	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i].Raw != w {
			t.Errorf("index %d: want %q, got %q", i, w, got[i].Raw)
		}
	}
}

func TestDedupeLinesEmpty(t *testing.T) {
	got := DedupeLines(nil, 10)
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d lines", len(got))
	}
}

func TestNewDeduplicatorZeroCapacity(t *testing.T) {
	d := NewDeduplicator(0)
	if d.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", d.capacity)
	}
}
