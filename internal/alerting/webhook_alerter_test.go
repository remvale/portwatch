package alerting

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookAlerter_Send_Success(t *testing.T) {
	var received WebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	alerter := NewWebhookAlerter(server.URL, 5*time.Second)
	alert := NewUnexpectedBindingAlert(8080, "tcp", 1234, "suspicious")

	if err := alerter.Send(alert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", received.Proto)
	}
	if received.PID != 1234 {
		t.Errorf("expected pid 1234, got %d", received.PID)
	}
	if received.Level == "" {
		t.Error("expected non-empty level")
	}
	if received.Time == "" {
		t.Error("expected non-empty time")
	}
}

func TestWebhookAlerter_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	alerter := NewWebhookAlerter(server.URL, 5*time.Second)
	alert := NewConflictAlert(9090, "udp")

	if err := alerter.Send(alert); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookAlerter_Send_BadURL(t *testing.T) {
	alerter := NewWebhookAlerter("http://127.0.0.1:0/no-server", 500*time.Millisecond)
	alert := NewConflictAlert(3000, "tcp")

	if err := alerter.Send(alert); err == nil {
		t.Fatal("expected connection error, got nil")
	}
}
