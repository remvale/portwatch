package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func entry(proto, addr, state string) portscanner.PortEntry {
	return portscanner.PortEntry{Protocol: proto, LocalAddress: addr, State: state}
}

func TestStore_UpdateAndLatest(t *testing.T) {
	s := NewStore()

	if s.Latest() != nil {
		t.Fatal("expected nil latest on empty store")
	}

	entries := []portscanner.PortEntry{entry("tcp", "0.0.0.0:8080", "LISTEN")}
	s.Update(entries)

	latest := s.Latest()
	if latest == nil {
		t.Fatal("expected non-nil latest after update")
	}
	if len(latest.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(latest.Ports))
	}
}

func TestStore_Diff_NoPrevious(t *testing.T) {
	s := NewStore()
	s.Update([]portscanner.PortEntry{entry("tcp", "0.0.0.0:80", "LISTEN")})

	appeared, disappeared := s.Diff()
	if appeared != nil || disappeared != nil {
		t.Fatal("expected nil diffs when no previous snapshot")
	}
}

func TestStore_Diff_Appeared(t *testing.T) {
	s := NewStore()
	s.Update([]portscanner.PortEntry{entry("tcp", "0.0.0.0:80", "LISTEN")})
	s.Update([]portscanner.PortEntry{
		entry("tcp", "0.0.0.0:80", "LISTEN"),
		entry("tcp", "0.0.0.0:9090", "LISTEN"),
	})

	appeared, disappeared := s.Diff()
	if len(appeared) != 1 {
		t.Fatalf("expected 1 appeared, got %d", len(appeared))
	}
	if appeared[0].LocalAddress != "0.0.0.0:9090" {
		t.Errorf("unexpected appeared address: %s", appeared[0].LocalAddress)
	}
	if len(disappeared) != 0 {
		t.Fatalf("expected 0 disappeared, got %d", len(disappeared))
	}
}

func TestStore_Diff_Disappeared(t *testing.T) {
	s := NewStore()
	s.Update([]portscanner.PortEntry{
		entry("tcp", "0.0.0.0:80", "LISTEN"),
		entry("tcp", "0.0.0.0:443", "LISTEN"),
	})
	s.Update([]portscanner.PortEntry{entry("tcp", "0.0.0.0:80", "LISTEN")})

	appeared, disappeared := s.Diff()
	if len(appeared) != 0 {
		t.Fatalf("expected 0 appeared, got %d", len(appeared))
	}
	if len(disappeared) != 1 {
		t.Fatalf("expected 1 disappeared, got %d", len(disappeared))
	}
	if disappeared[0].LocalAddress != "0.0.0.0:443" {
		t.Errorf("unexpected disappeared address: %s", disappeared[0].LocalAddress)
	}
}

func TestStore_Diff_NoChange(t *testing.T) {
	s := NewStore()
	ports := []portscanner.PortEntry{entry("tcp", "0.0.0.0:80", "LISTEN")}
	s.Update(ports)
	s.Update(ports)

	appeared, disappeared := s.Diff()
	if len(appeared) != 0 || len(disappeared) != 0 {
		t.Fatal("expected no diff when ports unchanged")
	}
}
