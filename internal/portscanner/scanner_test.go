package portscanner

import (
	"testing"
)

func TestNewScanner(t *testing.T) {
	s := NewScanner()
	if s == nil {
		t.Fatal("expected non-nil Scanner")
	}
}

func TestParseHexPort_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"0100007F:0050", 80},
		{"0100007F:01BB", 443},
		{"00000000:270F", 9999},
	}

	for _, tt := range tests {
		port, err := parseHexPort(tt.input)
		if err != nil {
			t.Errorf("parseHexPort(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if port != tt.expected {
			t.Errorf("parseHexPort(%q) = %d, want %d", tt.input, port, tt.expected)
		}
	}
}

func TestParseHexPort_Invalid(t *testing.T) {
	invalidInputs := []string{
		"",
		"noport",
		"addr:ZZZZ",
	}

	for _, input := range invalidInputs {
		_, err := parseHexPort(input)
		if err == nil {
			t.Errorf("parseHexPort(%q) expected error, got nil", input)
		}
	}
}

func TestParseHexAddr_Valid(t *testing.T) {
	addr := parseHexAddr("0100007F:0050")
	if addr != "0100007F" {
		t.Errorf("parseHexAddr expected '0100007F', got %q", addr)
	}
}

func TestParseHexAddr_Invalid(t *testing.T) {
	addr := parseHexAddr("badformat")
	if addr != "unknown" {
		t.Errorf("parseHexAddr expected 'unknown' for bad input, got %q", addr)
	}
}

func TestPortEntry_Fields(t *testing.T) {
	entry := PortEntry{
		Protocol:     "tcp",
		LocalAddress: "0100007F",
		LocalPort:    8080,
		PID:          1234,
		State:        "0A",
	}

	if entry.Protocol != "tcp" {
		t.Errorf("expected protocol 'tcp', got %q", entry.Protocol)
	}
	if entry.LocalPort != 8080 {
		t.Errorf("expected port 8080, got %d", entry.LocalPort)
	}
	if entry.PID != 1234 {
		t.Errorf("expected PID 1234, got %d", entry.PID)
	}
}
