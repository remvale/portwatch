package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.LogLevel != "info" {
		t.Errorf("expected log_level info, got %s", cfg.LogLevel)
	}
	if len(cfg.AllowedPorts) == 0 {
		t.Error("expected non-empty allowed_ports")
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	content := `{"scan_interval":"10s","allowed_ports":[22,80],"log_level":"debug"}`
	f := writeTempConfig(t, content)
	cfg, err := LoadFromFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval.Seconds() != 10 {
		t.Errorf("expected 10s interval, got %v", cfg.ScanInterval)
	}
}

func TestLoadFromFile_InvalidDuration(t *testing.T) {
	content := `{"scan_interval":"notaduration","log_level":"info"}`
	f := writeTempConfig(t, content)
	_, err := LoadFromFile(f)
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestLoadFromFile_InvalidLogLevel(t *testing.T) {
	content := `{"scan_interval":"5s","log_level":"verbose"}`
	f := writeTempConfig(t, content)
	_, err := LoadFromFile(f)
	if err == nil {
		t.Fatal("expected error for invalid log level")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFromFile_EmailDefaultPort(t *testing.T) {
	content := `{"scan_interval":"5s","log_level":"info","email":{"enabled":true,"smtp_host":"smtp.example.com","from":"a@b.com","to":["c@d.com"]}}`
	f := writeTempConfig(t, content)
	cfg, err := LoadFromFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EmailConfig.SMTPPort != 587 {
		t.Errorf("expected default smtp port 587, got %d", cfg.EmailConfig.SMTPPort)
	}
}

func TestLoadFromFile_EmailExplicitPort(t *testing.T) {
	content := `{"scan_interval":"5s","log_level":"info","email":{"enabled":true,"smtp_host":"smtp.example.com","smtp_port":465,"from":"a@b.com","to":["c@d.com"]}}`
	f := writeTempConfig(t, content)
	cfg, err := LoadFromFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EmailConfig.SMTPPort != 465 {
		t.Errorf("expected smtp port 465, got %d", cfg.EmailConfig.SMTPPort)
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}
