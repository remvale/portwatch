package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	// AllowedPorts is the list of ports that are expected to be bound.
	AllowedPorts []uint16 `json:"allowed_ports"`

	// ScanInterval is how often the port scanner runs.
	ScanInterval Duration `json:"scan_interval"`

	// LogLevel controls verbosity ("info", "warn", "error").
	LogLevel string `json:"log_level"`

	// AlertOnConflict enables alerting when two processes bind the same port.
	AlertOnConflict bool `json:"alert_on_conflict"`
}

// Duration is a wrapper around time.Duration for JSON unmarshalling.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	d.Duration = parsed
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		AllowedPorts:    []uint16{},
		ScanInterval:    Duration{Duration: 5 * time.Second},
		LogLevel:        "info",
		AlertOnConflict: true,
	}
}

// LoadFromFile reads and parses a JSON config file at the given path.
// Missing fields fall back to defaults.
func LoadFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks that the config values are sensible.
func (c *Config) Validate() error {
	if c.ScanInterval.Duration <= 0 {
		return fmt.Errorf("scan_interval must be positive, got %s", c.ScanInterval.Duration)
	}
	switch c.LogLevel {
	case "info", "warn", "error", "debug":
	default:
		return fmt.Errorf("log_level must be one of info/warn/error/debug, got %q", c.LogLevel)
	}
	return nil
}
