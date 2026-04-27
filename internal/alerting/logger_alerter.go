package alerting

import (
	"fmt"
	"io"
	"os"
)

// LoggerAlerter writes alerts to an io.Writer (default: stderr).
type LoggerAlerter struct {
	out io.Writer
}

// NewLoggerAlerter creates a LoggerAlerter writing to the given writer.
// If w is nil, os.Stderr is used.
func NewLoggerAlerter(w io.Writer) *LoggerAlerter {
	if w == nil {
		w = os.Stderr
	}
	return &LoggerAlerter{out: w}
}

// Send writes the alert as a formatted line to the configured writer.
func (l *LoggerAlerter) Send(alert Alert) error {
	_, err := fmt.Fprintln(l.out, alert.String())
	return err
}

// MultiAlerter fans out an alert to multiple Alerter implementations.
type MultiAlerter struct {
	alerters []Alerter
}

// NewMultiAlerter creates a MultiAlerter from the provided alerters.
func NewMultiAlerter(alerters ...Alerter) *MultiAlerter {
	return &MultiAlerter{alerters: alerters}
}

// Send dispatches the alert to all registered alerters, collecting errors.
func (m *MultiAlerter) Send(alert Alert) error {
	var firstErr error
	for _, a := range m.alerters {
		if err := a.Send(alert); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
