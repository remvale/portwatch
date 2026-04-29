package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Config holds the full portwatch configuration.
type Config struct {
	ScanInterval  time.Duration `json:"-"`
	RawInterval   string        `json:"scan_interval"`
	AllowedPorts  []int         `json:"allowed_ports"`
	LogLevel      string        `json:"log_level"`
	WebhookURL    string        `json:"webhook_url"`
	EmailConfig   EmailConfig   `json:"email"`
}

// EmailConfig mirrors alerting.EmailConfig for JSON unmarshalling.
type EmailConfig struct {
	Enabled  bool     `json:"enabled"`
	SMTPHost string   `json:"smtp_host"`
	SMTPPort int      `json:"smtp_port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

var validLogLevels = map[string]bool{
	"debug": true, "info": true, "warn": true, "error": true,
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		RawInterval:  "5s",
		ScanInterval: 5 * time.Second,
		AllowedPorts: []int{22, 80, 443},
		LogLevel:     "info",
	}
}

// LoadFromFile reads and validates a JSON config file.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.RawInterval != "" {
		d, err := time.ParseDuration(cfg.RawInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid scan_interval %q: %w", cfg.RawInterval, err)
		}
		cfg.ScanInterval = d
	}
	if !validLogLevels[cfg.LogLevel] {
		return nil, fmt.Errorf("invalid log_level %q: must be one of debug, info, warn, error", cfg.LogLevel)
	}
	if cfg.EmailConfig.Enabled && cfg.EmailConfig.SMTPPort == 0 {
		cfg.EmailConfig.SMTPPort = 587
	}
	return cfg, nil
}
