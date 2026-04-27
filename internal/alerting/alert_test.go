package alerting

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAlertLevelString(t *testing.T) {
	cases := []struct {
		level    AlertLevel
		expected string
	}{
		{AlertInfo, "INFO"},
		{AlertWarning, "WARNING"},
		{AlertCritical, "CRITICAL"},
		{AlertLevel(99), "UNKNOWN"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.expected {
			t.Errorf("AlertLevel(%d).String() = %q, want %q", c.level, got, c.expected)
		}
	}
}

func TestNewUnexpectedBindingAlert(t *testing.T) {
	before := time.Now()
	a := NewUnexpectedBindingAlert(8080, 1234, "nginx")
	after := time.Now()

	if a.Level != AlertWarning {
		t.Errorf("expected WARNING, got %s", a.Level)
	}
	if a.Port != 8080 || a.PID != 1234 || a.Process != "nginx" {
		t.Errorf("unexpected alert fields: %+v", a)
	}
	if a.Timestamp.Before(before) || a.Timestamp.After(after) {
		t.Errorf("timestamp out of expected range")
	}
}

func TestNewConflictAlert(t *testing.T) {
	a := NewConflictAlert(443, 5678, "apache")
	if a.Level != AlertCritical {
		t.Errorf("expected CRITICAL, got %s", a.Level)
	}
	if !strings.Contains(a.Message, "443") {
		t.Errorf("expected port in message, got %q", a.Message)
	}
}

func TestLoggerAlerterSend(t *testing.T) {
	var buf bytes.Buffer
	la := NewLoggerAlerter(&buf)
	a := NewUnexpectedBindingAlert(3000, 42, "node")

	if err := la.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "WARNING") || !strings.Contains(output, "3000") {
		t.Errorf("unexpected log output: %q", output)
	}
}

func TestMultiAlerterSend(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	ma := NewMultiAlerter(NewLoggerAlerter(&buf1), NewLoggerAlerter(&buf2))
	a := NewConflictAlert(80, 99, "caddy")

	if err := ma.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, buf := range []*bytes.Buffer{&buf1, &buf2} {
		if !strings.Contains(buf.String(), "CRITICAL") {
			t.Errorf("expected CRITICAL in output, got %q", buf.String())
		}
	}
}
