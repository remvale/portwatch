package alerting

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds SMTP configuration for the email alerter.
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	To       []string
}

// emailAlerter sends alerts via email using SMTP.
type emailAlerter struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailAlerter creates a new Alerter that sends emails via SMTP.
func NewEmailAlerter(cfg EmailConfig) Alerter {
	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	}
	return &emailAlerter{cfg: cfg, auth: auth}
}

func (e *emailAlerter) Send(alert Alert) error {
	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)
	subject := fmt.Sprintf("[portwatch] %s alert: %s", alert.Level, alert.Title)
	body := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\r\n\r\nDetails: %v",
		strings.Join(e.cfg.To, ", "),
		e.cfg.From,
		subject,
		alert.Message,
		alert.Details,
	)
	return smtp.SendMail(addr, e.auth, e.cfg.From, e.cfg.To, []byte(body))
}
