package parser

import (
	"hash/fnv"
	"sync"
)

// Deduplicator filters out duplicate log lines within a sliding seen-set.
// It hashes each line's raw content and discards lines whose hash was
// already observed within the last `capacity` entries.
type Deduplicator struct {
	mu       sync.Mutex
	capacity int
	seen     map[uint64]struct{}
	order    []uint64
}

// NewDeduplicator returns a Deduplicator that remembers up to capacity
// distinct line hashes. When full, the oldest hash is evicted.
func NewDeduplicator(capacity int) *Deduplicator {
	if capacity <= 0 {
		capacity = 1024
	}
	return &Deduplicator{
		capacity: capacity,
		seen:     make(map[uint64]struct{}, capacity),
		order:    make([]uint64, 0, capacity),
	}
}

// IsDuplicate reports whether line has been seen before and, if not,
// records it so future calls with the same content return true.
func (d *Deduplicator) IsDuplicate(line Line) bool {
	h := hashLine(line.Raw)

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[h]; exists {
		return true
	}

	// Evict oldest entry when at capacity.
	if len(d.order) >= d.capacity {
		old := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, old)
	}

	d.seen[h] = struct{}{}
	d.order = append(d.order, h)
	return false
}

// DedupeLines filters lines, returning only those not seen before.
func DedupeLines(lines []Line, capacity int) []Line {
	d := NewDeduplicator(capacity)
	out := make([]Line, 0, len(lines))
	for _, l := range lines {
		if !d.IsDuplicate(l) {
			out = append(out, l)
		}
	}
	return out
}

func hashLine(raw string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(raw))
	return h.Sum64()
}
