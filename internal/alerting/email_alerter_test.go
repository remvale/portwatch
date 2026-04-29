package alerting

import (
	"net"
	"net/smtp"
	"strings"
	"testing"
	"time"
)

// startFakeSMTP starts a minimal fake SMTP server on a random port and
// returns the address and a channel that receives the raw message data.
func startFakeSMTP(t *testing.T) (addr string, msgCh chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake SMTP: %v", err)
	}
	msgCh = make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		defer ln.Close()
		buf := make([]byte, 4096)
		conn.Write([]byte("220 fake SMTP ready\r\n"))
		var collected strings.Builder
		for {
			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, _ := conn.Read(buf)
			if n == 0 {
				break
			}
			line := string(buf[:n])
			collected.WriteString(line)
			switch {
			case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"):
				conn.Write([]byte("250 OK\r\n"))
			case strings.HasPrefix(line, "MAIL FROM"):
				conn.Write([]byte("250 OK\r\n"))
			case strings.HasPrefix(line, "RCPT TO"):
				conn.Write([]byte("250 OK\r\n"))
			case strings.TrimSpace(line) == "DATA":
				conn.Write([]byte("354 Start input\r\n"))
			case strings.Contains(line, "\r\n.\r\n"):
				conn.Write([]byte("250 OK\r\n"))
			case strings.HasPrefix(line, "QUIT"):
				conn.Write([]byte("221 Bye\r\n"))
				msgCh <- collected.String()
				return
			}
		}
	}()
	return ln.Addr().String(), msgCh
}

func TestEmailAlerter_Send_BadAddress(t *testing.T) {
	cfg := EmailConfig{
		SMTPHost: "127.0.0.1",
		SMTPPort: 1, // nothing listening
		From:     "portwatch@example.com",
		To:       []string{"admin@example.com"},
	}
	a := NewEmailAlerter(cfg)
	alert := NewUnexpectedBindingAlert("127.0.0.1", 9999)
	err := a.Send(alert)
	if err == nil {
		t.Fatal("expected error sending to bad address, got nil")
	}
}

func TestEmailAlerter_Send_NoAuth(t *testing.T) {
	_ = smtp.PlainAuth // ensure import used
	cfg := EmailConfig{
		SMTPHost: "127.0.0.1",
		SMTPPort: 1,
		From:     "portwatch@example.com",
		To:       []string{"admin@example.com"},
	}
	a := NewEmailAlerter(cfg)
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}
