package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/portscanner"
)

// Config holds the configuration for the Monitor.
type Config struct {
	// Interval is how often the monitor polls port bindings.
	Interval time.Duration
	// AllowedPorts is the set of ports considered expected.
	AllowedPorts map[int]bool
}

// Monitor watches port bindings and emits alerts via an Alerter.
type Monitor struct {
	cfg     Config
	scanner *portscanner.Scanner
	alerter alerting.Alerter
	prev    map[string]portscanner.PortEntry
}

// New creates a new Monitor.
func New(cfg Config, scanner *portscanner.Scanner, alerter alerting.Alerter) *Monitor {
	return &Monitor{
		cfg:     cfg,
		scanner: scanner,
		alerter: alerter,
		prev:    make(map[string]portscanner.PortEntry),
	}
}

// Run starts the monitoring loop. It blocks until the done channel is closed.
func (m *Monitor) Run(done <-chan struct{}) {
	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	log.Println("[monitor] starting port watch loop")
	for {
		select {
		case <-done:
			log.Println("[monitor] stopping")
			return
		case <-ticker.C:
			m.scan()
		}
	}
}

// scan performs a single scan cycle.
func (m *Monitor) scan() {
	entries, err := m.scanner.Scan()
	if err != nil {
		log.Printf("[monitor] scan error: %v", err)
		return
	}

	current := make(map[string]portscanner.PortEntry, len(entries))
	for _, e := range entries {
		key := e.Key()
		current[key] = e

		// Alert on unexpected new bindings.
		if _, seen := m.prev[key]; !seen {
			if !m.cfg.AllowedPorts[e.LocalPort] {
				alert := alerting.NewUnexpectedBindingAlert(e.LocalAddr, e.LocalPort, e.PID)
				if err := m.alerter.Send(alert); err != nil {
					log.Printf("[monitor] alert send error: %v", err)
				}
			}
		}
	}

	// Detect port conflicts (same port, multiple PIDs).
	portPIDs := make(map[int][]int)
	for _, e := range entries {
		portPIDs[e.LocalPort] = append(portPIDs[e.LocalPort], e.PID)
	}
	for port, pids := range portPIDs {
		if len(pids) > 1 {
			alert := alerting.NewConflictAlert(port, pids)
			if err := m.alerter.Send(alert); err != nil {
				log.Printf("[monitor] conflict alert send error: %v", err)
			}
		}
	}

	m.prev = current
}
