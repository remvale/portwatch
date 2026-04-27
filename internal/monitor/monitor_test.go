package monitor_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/portscanner"
)

// stubScanner returns a fixed list of PortEntry values.
type stubScanner struct {
	entries []portscanner.PortEntry
	err     error
}

func (s *stubScanner) Scan() ([]portscanner.PortEntry, error) {
	return s.entries, s.err
}

// captureAlerter records every alert it receives.
type captureAlerter struct {
	alerts []alerting.Alert
}

func (c *captureAlerter) Send(a alerting.Alert) error {
	c.alerts = append(c.alerts, a)
	return nil
}

func TestMonitor_UnexpectedBinding(t *testing.T) {
	scanner := &stubScanner{
		entries: []portscanner.PortEntry{
			{LocalAddr: "127.0.0.1", LocalPort: 9999, PID: 42},
		},
	}
	cap := &captureAlerter{}
	cfg := monitor.Config{
		Interval:     10 * time.Millisecond,
		AllowedPorts: map[int]bool{80: true, 443: true},
	}

	m := monitor.New(cfg, portscanner.WrapScanner(scanner), cap)
	done := make(chan struct{})
	go m.Run(done)

	time.Sleep(30 * time.Millisecond)
	close(done)
	time.Sleep(10 * time.Millisecond)

	if len(cap.alerts) == 0 {
		t.Fatal("expected at least one unexpected-binding alert, got none")
	}
	for _, a := range cap.alerts {
		if a.Level != alerting.LevelWarn {
			t.Errorf("expected level Warn, got %v", a.Level)
		}
	}
}

func TestMonitor_AllowedPort_NoAlert(t *testing.T) {
	scanner := &stubScanner{
		entries: []portscanner.PortEntry{
			{LocalAddr: "0.0.0.0", LocalPort: 80, PID: 1},
		},
	}
	cap := &captureAlerter{}
	cfg := monitor.Config{
		Interval:     10 * time.Millisecond,
		AllowedPorts: map[int]bool{80: true},
	}

	m := monitor.New(cfg, portscanner.WrapScanner(scanner), cap)
	done := make(chan struct{})
	go m.Run(done)

	time.Sleep(30 * time.Millisecond)
	close(done)
	time.Sleep(10 * time.Millisecond)

	for _, a := range cap.alerts {
		if a.Type == alerting.TypeUnexpectedBinding {
			t.Errorf("did not expect an unexpected-binding alert for allowed port 80")
		}
	}
}
