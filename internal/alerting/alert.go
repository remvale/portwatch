package alerting

import (
	"fmt"
	"time"
)

// AlertLevel represents the severity of an alert.
type AlertLevel int

const (
	AlertInfo AlertLevel = iota
	AlertWarning
	AlertCritical
)

func (l AlertLevel) String() string {
	switch l {
	case AlertInfo:
		return "INFO"
	case AlertWarning:
		return "WARNING"
	case AlertCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Alert represents a port-related alert event.
type Alert struct {
	Level     AlertLevel
	Message   string
	Port      int
	PID       int
	Process   string
	Timestamp time.Time
}

func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s | port=%d pid=%d process=%s ts=%s",
		a.Level, a.Message, a.Port, a.PID, a.Process,
		a.Timestamp.Format(time.RFC3339),
	)
}

// Alerter defines the interface for sending alerts.
type Alerter interface {
	Send(alert Alert) error
}

// NewUnexpectedBindingAlert creates an alert for an unexpected port binding.
func NewUnexpectedBindingAlert(port, pid int, process string) Alert {
	return Alert{
		Level:     AlertWarning,
		Message:   fmt.Sprintf("unexpected binding detected on port %d", port),
		Port:      port,
		PID:       pid,
		Process:   process,
		Timestamp: time.Now(),
	}
}

// NewConflictAlert creates an alert for a port conflict.
func NewConflictAlert(port, pid int, process string) Alert {
	return Alert{
		Level:     AlertCritical,
		Message:   fmt.Sprintf("port conflict detected on port %d", port),
		Port:      port,
		PID:       pid,
		Process:   process,
		Timestamp: time.Now(),
	}
}
