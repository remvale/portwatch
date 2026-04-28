package snapshot

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

func TestReport_NilSnapshot(t *testing.T) {
	var buf bytes.Buffer
	if err := Report(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no snapshot available") {
		t.Errorf("expected 'no snapshot available', got: %s", buf.String())
	}
}

func TestReport_WithPorts(t *testing.T) {
	snap := &PortSnapshot{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Ports: []portscanner.PortEntry{
			{Protocol: "tcp", LocalAddress: "0.0.0.0:80", State: "LISTEN"},
			{Protocol: "tcp", LocalAddress: "0.0.0.0:443", State: "LISTEN"},
		},
	}

	var buf bytes.Buffer
	if err := Report(&buf, snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "0.0.0.0:80") {
		t.Errorf("expected port 80 in output")
	}
	if !strings.Contains(out, "0.0.0.0:443") {
		t.Errorf("expected port 443 in output")
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Errorf("expected timestamp in output")
	}
}

func TestDiffReport_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	if err := DiffReport(&buf, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes detected") {
		t.Errorf("expected 'no changes detected'")
	}
}

func TestDiffReport_WithChanges(t *testing.T) {
	appeared := []portscanner.PortEntry{
		{Protocol: "tcp", LocalAddress: "0.0.0.0:9090", State: "LISTEN"},
	}
	disappeared := []portscanner.PortEntry{
		{Protocol: "tcp", LocalAddress: "0.0.0.0:8080", State: "LISTEN"},
	}

	var buf bytes.Buffer
	if err := DiffReport(&buf, appeared, disappeared); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+ tcp") {
		t.Errorf("expected appeared entry with '+' prefix")
	}
	if !strings.Contains(out, "- tcp") {
		t.Errorf("expected disappeared entry with '-' prefix")
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected appeared port 9090")
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected disappeared port 8080")
	}
}
