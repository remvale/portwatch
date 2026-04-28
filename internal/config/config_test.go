package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ScanInterval.Duration != 5*time.Second {
		t.Errorf("expected 5s scan interval, got %s", cfg.ScanInterval.Duration)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected log level 'info', got %q", cfg.LogLevel)
	}
	if !cfg.AlertOnConflict {
		t.Error("expected AlertOnConflict to be true by default")
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	content := `{
		"allowed_ports": [80, 443, 8080],
		"scan_interval": "10s",
		"log_level": "warn",
		"alert_on_conflict": false
	}`

	f := writeTempFile(t, content)
	defer os.Remove(f)

	cfg, err := LoadFromFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.AllowedPorts) != 3 {
		t.Errorf("expected 3 allowed ports, got %d", len(cfg.AllowedPorts))
	}
	if cfg.ScanInterval.Duration != 10*time.Second {
		t.Errorf("expected 10s, got %s", cfg.ScanInterval.Duration)
	}
	if cfg.LogLevel != "warn" {
		t.Errorf("expected 'warn', got %q", cfg.LogLevel)
	}
	if cfg.AlertOnConflict {
		t.Error("expected AlertOnConflict false")
	}
}

func TestLoadFromFile_InvalidDuration(t *testing.T) {
	content := `{"scan_interval": "not-a-duration"}`
	f := writeTempFile(t, content)
	defer os.Remove(f)

	_, err := LoadFromFile(f)
	if err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
}

func TestLoadFromFile_InvalidLogLevel(t *testing.T) {
	content := `{"log_level": "verbose", "scan_interval": "5s"}`
	f := writeTempFile(t, content)
	defer os.Remove(f)

	_, err := LoadFromFile(f)
	if err == nil {
		t.Fatal("expected error for invalid log level, got nil")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestValidate_ZeroDuration(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ScanInterval = Duration{Duration: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for zero duration")
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}
