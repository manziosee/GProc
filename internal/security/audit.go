package security

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gproc/pkg/types"
)

type AuditLogger struct {
	config    *types.AuditConfig
	fileLog   *os.File
}

type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Target    string    `json:"target"`
	Result    string    `json:"result"` // success, failure, error
	Details   string    `json:"details,omitempty"`
	IP        string    `json:"ip,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

func NewAuditLogger(config *types.AuditConfig) (*AuditLogger, error) {
	if !config.Enabled {
		return &AuditLogger{config: config}, nil
	}

	al := &AuditLogger{config: config}

	// Initialize file logging
	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open audit log file: %v", err)
		}
		al.fileLog = file
	}

	// Note: syslog is not supported cross-platform in stdlib; we'll simulate syslog format to file/stdout

	return al, nil
}

func (al *AuditLogger) LogEvent(event AuditEvent) error {
	if !al.config.Enabled {
		return nil
	}

	event.Timestamp = time.Now()

	switch al.config.Format {
	case "json":
		return al.logJSON(event)
	case "syslog":
		return al.logSyslog(event)
	default:
		return al.logPlainText(event)
	}
}

func (al *AuditLogger) logJSON(event AuditEvent) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if al.fileLog != nil {
		_, err = al.fileLog.WriteString(string(jsonData) + "\n")
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func (al *AuditLogger) logSyslog(event AuditEvent) error {
	message := fmt.Sprintf("user=%s action=%s resource=%s target=%s result=%s",
		event.Username, event.Action, event.Resource, event.Target, event.Result)
	if al.fileLog != nil {
		_, err := al.fileLog.WriteString(message + "\n")
		return err
	}
	fmt.Println(message)
	return nil
}

func (al *AuditLogger) logPlainText(event AuditEvent) error {
	message := fmt.Sprintf("[%s] %s (%s) %s %s:%s -> %s",
		event.Timestamp.Format(time.RFC3339),
		event.Username, event.UserID,
		event.Action, event.Resource, event.Target,
		event.Result)

	if event.Details != "" {
		message += " - " + event.Details
	}

	if al.fileLog != nil {
		_, err := al.fileLog.WriteString(message + "\n")
		return err
	}

	fmt.Println(message)
	return nil
}

func (al *AuditLogger) LogUserLogin(user *types.User, ip, userAgent string, success bool) {
	result := "success"
	if !success {
		result = "failure"
	}

	al.LogEvent(AuditEvent{
		UserID:    user.ID,
		Username:  user.Username,
		Action:    "login",
		Resource:  "auth",
		Target:    "session",
		Result:    result,
		IP:        ip,
		UserAgent: userAgent,
	})
}

func (al *AuditLogger) LogProcessAction(user *types.User, action, processName, result string) {
	al.LogEvent(AuditEvent{
		UserID:   user.ID,
		Username: user.Username,
		Action:   action,
		Resource: "process",
		Target:   processName,
		Result:   result,
	})
}

func (al *AuditLogger) LogConfigChange(user *types.User, configType, target, result string) {
	al.LogEvent(AuditEvent{
		UserID:   user.ID,
		Username: user.Username,
		Action:   "modify",
		Resource: "config",
		Target:   fmt.Sprintf("%s:%s", configType, target),
		Result:   result,
	})
}

func (al *AuditLogger) LogClusterAction(user *types.User, action, nodeID, result string) {
	al.LogEvent(AuditEvent{
		UserID:   user.ID,
		Username: user.Username,
		Action:   action,
		Resource: "cluster",
		Target:   nodeID,
		Result:   result,
	})
}

func (al *AuditLogger) Close() error {
	if al.fileLog != nil {
		return al.fileLog.Close()
	}
	return nil
}