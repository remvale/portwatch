package snapshot

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// PortSnapshot holds a point-in-time view of active port bindings.
type PortSnapshot struct {
	Timestamp time.Time
	Ports     []portscanner.PortEntry
}

// Store maintains the latest and previous snapshots for diff computation.
type Store struct {
	mu       sync.RWMutex
	previous *PortSnapshot
	latest   *PortSnapshot
}

// NewStore returns an initialised snapshot Store.
func NewStore() *Store {
	return &Store{}
}

// Update records a new snapshot, rotating the previous one.
func (s *Store) Update(entries []portscanner.PortEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.previous = s.latest
	s.latest = &PortSnapshot{
		Timestamp: time.Now(),
		Ports:     entries,
	}
}

// Diff returns ports that are new (appeared) and ports that are gone (disappeared)
// compared to the previous snapshot. Returns nil slices if no previous snapshot exists.
func (s *Store) Diff() (appeared, disappeared []portscanner.PortEntry) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.previous == nil || s.latest == nil {
		return nil, nil
	}

	prevSet := toSet(s.previous.Ports)
	currSet := toSet(s.latest.Ports)

	for _, e := range s.latest.Ports {
		if _, ok := prevSet[key(e)]; !ok {
			appeared = append(appeared, e)
		}
	}

	for _, e := range s.previous.Ports {
		if _, ok := currSet[key(e)]; !ok {
			disappeared = append(disappeared, e)
		}
	}

	return appeared, disappeared
}

// Latest returns the most recent snapshot or nil.
func (s *Store) Latest() *PortSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latest
}

func key(e portscanner.PortEntry) string {
	return e.Protocol + "|" + e.LocalAddress + "|" + e.State
}

func toSet(entries []portscanner.PortEntry) map[string]struct{} {
	m := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		m[key(e)] = struct{}{}
	}
	return m
}
