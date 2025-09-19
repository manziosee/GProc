package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"gproc/pkg/types"
)

type AlertManager struct {
	alerts    []types.Alert
	config    *AlertConfig
	notifiers map[string]Notifier
}

type AlertConfig struct {
	EmailEnabled bool   `json:"email_enabled"`
	SMTPServer   string `json:"smtp_server"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPass     string `json:"smtp_pass"`
	SlackEnabled bool   `json:"slack_enabled"`
	SlackWebhook string `json:"slack_webhook"`
	SMSEnabled   bool   `json:"sms_enabled"`
	TwilioSID    string `json:"twilio_sid"`
	TwilioToken  string `json:"twilio_token"`
	TwilioFrom   string `json:"twilio_from"`
}

type Notifier interface {
	Send(alert *types.Alert) error
}

type EmailNotifier struct {
	config *AlertConfig
}

type SlackNotifier struct {
	webhookURL string
}

type SMSNotifier struct {
	config *AlertConfig
}

func NewAlertManager(config *AlertConfig) *AlertManager {
	am := &AlertManager{
		alerts:    []types.Alert{},
		config:    config,
		notifiers: make(map[string]Notifier),
	}
	
	if config.EmailEnabled {
		am.notifiers["email"] = &EmailNotifier{config: config}
	}
	
	if config.SlackEnabled {
		am.notifiers["slack"] = &SlackNotifier{webhookURL: config.SlackWebhook}
	}
	
	if config.SMSEnabled {
		am.notifiers["sms"] = &SMSNotifier{config: config}
	}
	
	return am
}

func (am *AlertManager) TriggerAlert(processID, alertType, message, severity string) error {
	alert := types.Alert{
		ID:           fmt.Sprintf("alert-%d", time.Now().Unix()),
		ProcessID:    processID,
		Type:         alertType,
		Message:      message,
		Severity:     severity,
		Timestamp:    time.Now(),
		Acknowledged: false,
	}
	
	am.alerts = append(am.alerts, alert)
	
	// Send notifications
	for _, notifier := range am.notifiers {
		go notifier.Send(&alert)
	}
	
	return nil
}

func (am *AlertManager) GetAlerts() []types.Alert {
	return am.alerts
}

func (am *AlertManager) AcknowledgeAlert(alertID string) error {
	for i := range am.alerts {
		if am.alerts[i].ID == alertID {
			am.alerts[i].Acknowledged = true
			return nil
		}
	}
	return fmt.Errorf("alert %s not found", alertID)
}

func (am *AlertManager) ClearAlerts() {
	am.alerts = []types.Alert{}
}

// Email Notifier Implementation
func (e *EmailNotifier) Send(alert *types.Alert) error {
	if !e.config.EmailEnabled {
		return nil
	}
	
	subject := fmt.Sprintf("[GProc Alert] %s - %s", alert.Severity, alert.ProcessID)
	body := fmt.Sprintf(`
Alert Details:
- Process: %s
- Type: %s
- Severity: %s
- Message: %s
- Time: %s

This is an automated alert from GProc.
`, alert.ProcessID, alert.Type, alert.Severity, alert.Message, alert.Timestamp.Format("2006-01-02 15:04:05"))
	
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	
	auth := smtp.PlainAuth("", e.config.SMTPUser, e.config.SMTPPass, e.config.SMTPServer)
	addr := fmt.Sprintf("%s:%d", e.config.SMTPServer, e.config.SMTPPort)
	
	return smtp.SendMail(addr, auth, e.config.SMTPUser, []string{"admin@example.com"}, []byte(msg))
}

// Slack Notifier Implementation
func (s *SlackNotifier) Send(alert *types.Alert) error {
	payload := map[string]interface{}{
		"text": fmt.Sprintf("🚨 *GProc Alert*\n*Process:* %s\n*Severity:* %s\n*Message:* %s", 
			alert.ProcessID, alert.Severity, alert.Message),
		"username": "GProc",
		"icon_emoji": ":warning:",
	}
	
	jsonPayload, _ := json.Marshal(payload)
	
	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// SMS Notifier Implementation (Twilio)
func (sms *SMSNotifier) Send(alert *types.Alert) error {
	if !sms.config.SMSEnabled {
		return nil
	}
	
	message := fmt.Sprintf("GProc Alert: %s - %s (%s)", 
		alert.ProcessID, alert.Message, alert.Severity)
	
	// Twilio API call would go here
	fmt.Printf("SMS Alert: %s\n", message)
	
	return nil
}