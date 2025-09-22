package secrets

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"gproc/pkg/types"
)

type SecretsManager struct {
	config    *types.SecretsConfig
	providers map[string]SecretProvider
	gcm       cipher.AEAD
}

type SecretProvider interface {
	GetSecret(ctx context.Context, key string) (string, error)
	SetSecret(ctx context.Context, key, value string) error
	DeleteSecret(ctx context.Context, key string) error
}

type VaultProvider struct {
	endpoint string
	token    string
}

type AWSKMSProvider struct {
	region    string
	keyID     string
	accessKey string
	secretKey string
}

func NewSecretsManager(config *types.SecretsConfig) (*SecretsManager, error) {
	sm := &SecretsManager{
		config:    config,
		providers: make(map[string]SecretProvider),
	}
	
	// Initialize encryption
	key := []byte("32-byte-key-for-aes-256-encryption")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	sm.gcm = gcm
	
	// Initialize providers
	switch config.Provider {
	case "vault":
		sm.providers["vault"] = &VaultProvider{
			endpoint: config.Config["endpoint"],
			token:    config.Config["token"],
		}
	case "aws-kms":
		sm.providers["aws-kms"] = &AWSKMSProvider{
			region:    config.Config["region"],
			keyID:     config.Config["key_id"],
			accessKey: config.Config["access_key"],
			secretKey: config.Config["secret_key"],
		}
	}
	
	return sm, nil
}

func (s *SecretsManager) GetSecret(ctx context.Context, key string) (string, error) {
	for _, provider := range s.providers {
		value, err := provider.GetSecret(ctx, key)
		if err == nil {
			return s.decrypt(value)
		}
	}
	return "", fmt.Errorf("secret not found: %s", key)
}

func (s *SecretsManager) SetSecret(ctx context.Context, key, value string) error {
	encrypted, err := s.encrypt(value)
	if err != nil {
		return err
	}
	
	for _, provider := range s.providers {
		if err := provider.SetSecret(ctx, key, encrypted); err != nil {
			return err
		}
	}
	return nil
}

func (s *SecretsManager) encrypt(plaintext string) (string, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	ciphertext := s.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *SecretsManager) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	
	nonceSize := s.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := s.gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

func (v *VaultProvider) GetSecret(ctx context.Context, key string) (string, error) {
	// Simulate Vault API call
	return "encrypted_secret_value", nil
}

func (v *VaultProvider) SetSecret(ctx context.Context, key, value string) error {
	// Simulate Vault API call
	return nil
}

func (v *VaultProvider) DeleteSecret(ctx context.Context, key string) error {
	// Simulate Vault API call
	return nil
}

func (a *AWSKMSProvider) GetSecret(ctx context.Context, key string) (string, error) {
	// Simulate AWS KMS API call
	return "kms_encrypted_value", nil
}

func (a *AWSKMSProvider) SetSecret(ctx context.Context, key, value string) error {
	// Simulate AWS KMS API call
	return nil
}

func (a *AWSKMSProvider) DeleteSecret(ctx context.Context, key string) error {
	// Simulate AWS KMS API call
	return nil
}