package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gproc/pkg/types"
)

type RBACManager struct {
	config *types.RBACConfig
	users  map[string]*types.User
	roles  map[string]*types.Role
}

func NewRBACManager(config *types.RBACConfig) *RBACManager {
	rbac := &RBACManager{
		config: config,
		users:  make(map[string]*types.User),
		roles:  make(map[string]*types.Role),
	}
	
	// Load roles
	for _, role := range config.Roles {
		rbac.roles[role.Name] = &role
	}
	
	// Load users
	for _, user := range config.Users {
		rbac.users[user.Username] = &user
	}
	
	return rbac
}

func (r *RBACManager) Authenticate(username, password string) (*types.User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	// In production, verify password hash
	if password == "admin" || password == "password" {
		user.LastSeen = time.Now()
		return user, nil
	}
	
	return nil, fmt.Errorf("invalid credentials")
}

func (r *RBACManager) Authorize(user *types.User, resource, action, scope string) bool {
	if !r.config.Enabled {
		return true
	}
	
	for _, roleName := range user.Roles {
		role, exists := r.roles[roleName]
		if !exists {
			continue
		}
		
		for _, perm := range role.Permissions {
			if r.matchesPermission(perm, resource, action, scope) {
				return true
			}
		}
	}
	
	return false
}

func (r *RBACManager) matchesPermission(perm types.Permission, resource, action, scope string) bool {
	if perm.Resource != "*" && perm.Resource != resource {
		return false
	}
	
	actionAllowed := false
	for _, allowedAction := range perm.Actions {
		if allowedAction == "*" || allowedAction == action {
			actionAllowed = true
			break
		}
	}
	if !actionAllowed {
		return false
	}
	
	if perm.Scope == "*" {
		return true
	}
	
	if strings.Contains(perm.Scope, ":") {
		parts := strings.SplitN(perm.Scope, ":", 2)
		if len(parts) == 2 && parts[0] == strings.Split(scope, ":")[0] {
			return parts[1] == "*" || parts[1] == strings.Split(scope, ":")[1]
		}
	}
	
	return perm.Scope == scope
}

type TokenManager struct {
	config *types.JWTConfig
}

func NewTokenManager(config *types.JWTConfig) *TokenManager {
	return &TokenManager{config: config}
}

func (t *TokenManager) GenerateToken(user *types.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"user":  user.Username,
		"email": user.Email,
		"roles": user.Roles,
		"iss":   t.config.Issuer,
		"exp":   time.Now().Add(t.config.Expiration).Unix(),
		"iat":   time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.config.Secret))
}

func (t *TokenManager) ValidateToken(tokenString string) (*types.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.config.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	
	user := &types.User{
		ID:       claims["sub"].(string),
		Username: claims["user"].(string),
		Email:    claims["email"].(string),
		Roles:    interfaceToStringSlice(claims["roles"]),
		LastSeen: time.Now(),
	}
	
	return user, nil
}

func interfaceToStringSlice(i interface{}) []string {
	if slice, ok := i.([]interface{}); ok {
		result := make([]string, len(slice))
		for i, v := range slice {
			result[i] = v.(string)
		}
		return result
	}
	return []string{}
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}