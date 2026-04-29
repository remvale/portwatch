package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
	PID     int    `json:"pid,omitempty"`
	Time    string `json:"time"`
}

// WebhookAlerter sends alerts to an HTTP endpoint as JSON POST requests.
type WebhookAlerter struct {
	URL    string
	client *http.Client
}

// NewWebhookAlerter creates a WebhookAlerter that posts to the given URL.
// timeout controls how long each HTTP request may take.
func NewWebhookAlerter(url string, timeout time.Duration) *WebhookAlerter {
	return &WebhookAlerter{
		URL: url,
		client: &http.Client{Timeout: timeout},
	}
}

// Send marshals the alert to JSON and POSTs it to the configured URL.
func (w *WebhookAlerter) Send(a Alert) error {
	payload := WebhookPayload{
		Level:   a.Level.String(),
		Message: a.Message,
		Port:    a.Port,
		Proto:   a.Proto,
		PID:     a.PID,
		Time:    a.Timestamp.UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook_alerter: marshal: %w", err)
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook_alerter: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook_alerter: unexpected status %d from %s", resp.StatusCode, w.URL)
	}
	return nil
}
