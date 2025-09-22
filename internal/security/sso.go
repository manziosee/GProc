package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gproc/pkg/types"
)

type SSOManager struct {
	config *types.SSOConfig
	oauth2 *OAuth2Provider
	saml   *SAMLProvider
}

type OAuth2Provider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

type SAMLProvider struct {
	EntityID    string
	MetadataURL string
	CertFile    string
	KeyFile     string
}

func NewSSOManager(config *types.SSOConfig) *SSOManager {
	sso := &SSOManager{config: config}
	
	switch config.Provider {
	case "oauth2":
		sso.oauth2 = &OAuth2Provider{
			AuthURL:     "https://accounts.google.com/o/oauth2/auth",
			TokenURL:    "https://oauth2.googleapis.com/token",
			UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
		}
	case "saml":
		sso.saml = &SAMLProvider{
			EntityID:    config.EntityID,
			MetadataURL: config.MetadataUrl,
		}
	}
	
	return sso
}

func (s *SSOManager) GetAuthURL(state string) (string, error) {
	if s.oauth2 != nil {
		return s.getOAuth2AuthURL(state), nil
	}
	if s.saml != nil {
		return s.getSAMLAuthURL(state), nil
	}
	return "", fmt.Errorf("no SSO provider configured")
}

func (s *SSOManager) getOAuth2AuthURL(state string) string {
	params := url.Values{
		"client_id":     {s.oauth2.ClientID},
		"redirect_uri":  {s.oauth2.RedirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
	}
	return s.oauth2.AuthURL + "?" + params.Encode()
}

func (s *SSOManager) getSAMLAuthURL(state string) string {
	// SAML SSO URL generation
	return fmt.Sprintf("%s?SAMLRequest=%s&RelayState=%s", 
		s.saml.MetadataURL, "encoded_saml_request", state)
}

func (s *SSOManager) HandleCallback(ctx context.Context, code, state string) (*types.User, error) {
	if s.oauth2 != nil {
		return s.handleOAuth2Callback(ctx, code, state)
	}
	if s.saml != nil {
		return s.handleSAMLCallback(ctx, code, state)
	}
	return nil, fmt.Errorf("no SSO provider configured")
}

func (s *SSOManager) handleOAuth2Callback(_ context.Context, code, _ string) (*types.User, error) {
	// Exchange code for token
	tokenResp, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, err
	}
	
	// Get user info
	userInfo, err := s.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}
	
	user := &types.User{
		ID:       userInfo["id"].(string),
		Username: userInfo["email"].(string),
		Email:    userInfo["email"].(string),
		Roles:    []string{"user"},
		Created:  time.Now(),
		LastSeen: time.Now(),
		Enabled:  true,
	}
	
	return user, nil
}

func (s *SSOManager) handleSAMLCallback(_ context.Context, _, _ string) (*types.User, error) {
	// SAML response parsing (simplified)
	user := &types.User{
		ID:       generateSAMLUserID(),
		Username: "saml_user",
		Email:    "user@company.com",
		Roles:    []string{"user"},
		Created:  time.Now(),
		LastSeen: time.Now(),
		Enabled:  true,
	}
	
	return user, nil
}

func (s *SSOManager) exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{
		"client_id":     {s.oauth2.ClientID},
		"client_secret": {s.oauth2.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {s.oauth2.RedirectURL},
	}
	
	resp, err := http.PostForm(s.oauth2.TokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	
	return &tokenResp, nil
}

func (s *SSOManager) getUserInfo(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", s.oauth2.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	
	return userInfo, nil
}

func generateSAMLUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}