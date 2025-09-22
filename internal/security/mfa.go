package security

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type MFAManager struct {
	config *types.MFAConfig
	users  map[string]*MFAUserData
}

type MFAUserData struct {
	UserID     string    `json:"user_id"`
	Secret     string    `json:"secret"`
	BackupCodes []string `json:"backup_codes"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

type MFAConfig struct {
	Enabled     bool   `json:"enabled"`
	Issuer      string `json:"issuer"`
	WindowSize  int    `json:"window_size"`
	BackupCodes int    `json:"backup_codes"`
}

func NewMFAManager(config *types.MFAConfig) *MFAManager {
	return &MFAManager{
		config: config,
		users:  make(map[string]*MFAUserData),
	}
}

func (m *MFAManager) GenerateSecret(userID string) (*MFASetup, error) {
	secret := generateTOTPSecret()
	backupCodes := generateBackupCodes(10)
	
	userData := &MFAUserData{
		UserID:      userID,
		Secret:      secret,
		BackupCodes: backupCodes,
		Enabled:     false,
		CreatedAt:   time.Now(),
	}
	
	m.users[userID] = userData
	
	qrCode := generateQRCodeURL(userID, secret, m.config.Issuer)
	
	return &MFASetup{
		Secret:      secret,
		QRCodeURL:   qrCode,
		BackupCodes: backupCodes,
	}, nil
}

func (m *MFAManager) EnableMFA(userID, token string) error {
	userData, exists := m.users[userID]
	if !exists {
		return fmt.Errorf("MFA not set up for user %s", userID)
	}
	
	if !m.ValidateTOTP(userData.Secret, token) {
		return fmt.Errorf("invalid TOTP token")
	}
	
	userData.Enabled = true
	return nil
}

func (m *MFAManager) DisableMFA(userID string) error {
	userData, exists := m.users[userID]
	if !exists {
		return fmt.Errorf("MFA not found for user %s", userID)
	}
	
	userData.Enabled = false
	return nil
}

func (m *MFAManager) ValidateToken(userID, token string) bool {
	userData, exists := m.users[userID]
	if !exists || !userData.Enabled {
		return true // MFA not enabled
	}
	
	// Check TOTP token
	if m.ValidateTOTP(userData.Secret, token) {
		return true
	}
	
	// Check backup codes
	for i, code := range userData.BackupCodes {
		if code == token {
			// Remove used backup code
			userData.BackupCodes = append(userData.BackupCodes[:i], userData.BackupCodes[i+1:]...)
			return true
		}
	}
	
	return false
}

func (m *MFAManager) ValidateTOTP(secret, token string) bool {
	now := time.Now().Unix() / 30
	
	// Check current window and adjacent windows for clock skew
	for i := -2; i <= 2; i++ {
		timeStep := now + int64(i)
		expectedToken := generateTOTP(secret, timeStep)
		if expectedToken == token {
			return true
		}
	}
	
	return false
}

func (m *MFAManager) GetMFAStatus(userID string) *MFAStatus {
	userData, exists := m.users[userID]
	if !exists {
		return &MFAStatus{
			Enabled:           false,
			BackupCodesCount:  0,
		}
	}
	
	return &MFAStatus{
		Enabled:          userData.Enabled,
		BackupCodesCount: len(userData.BackupCodes),
		CreatedAt:        userData.CreatedAt,
	}
}

func generateTOTPSecret() string {
	bytes := make([]byte, 20)
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

func generateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, 4)
		rand.Read(bytes)
		codes[i] = fmt.Sprintf("%08d", uint32(bytes[0])<<24|uint32(bytes[1])<<16|uint32(bytes[2])<<8|uint32(bytes[3]))
	}
	return codes
}

func generateQRCodeURL(userID, secret, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, userID, secret, issuer)
}

func generateTOTP(secret string, timeStep int64) string {
	// Simplified TOTP implementation
	hash := hmacSHA1(secret, timeStep)
	offset := hash[len(hash)-1] & 0x0F
	code := ((int(hash[offset]) & 0x7F) << 24) |
		((int(hash[offset+1]) & 0xFF) << 16) |
		((int(hash[offset+2]) & 0xFF) << 8) |
		(int(hash[offset+3]) & 0xFF)
	code = code % 1000000
	return fmt.Sprintf("%06d", code)
}

func hmacSHA1(secret string, counter int64) []byte {
	// Simplified HMAC-SHA1 implementation
	key := []byte(secret)
	data := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		data[i] = byte(counter & 0xFF)
		counter >>= 8
	}
	
	// Mock hash for demonstration
	hash := make([]byte, 20)
	for i := range hash {
		hash[i] = byte((int(key[i%len(key)]) + int(data[i%len(data)])) % 256)
	}
	return hash
}

type MFASetup struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

type MFAStatus struct {
	Enabled          bool      `json:"enabled"`
	BackupCodesCount int       `json:"backup_codes_count"`
	CreatedAt        time.Time `json:"created_at"`
}