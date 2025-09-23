package security

import (
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
	// Allow nil config by providing a safe default
	if config == nil {
		config = &types.RBACConfig{Enabled: false, Roles: []types.Role{}, Users: []types.User{}}
	}

	rbac := &RBACManager{
		config: config,
		users:  make(map[string]*types.User),
		roles:  make(map[string]*types.Role),
	}

	// Create default roles
	rbac.createDefaultRoles()
	
	// Load roles (avoid taking address of range variable)
	for i := range config.Roles {
		role := config.Roles[i]
		rbac.roles[role.Name] = &role
	}

	// Create default users
	rbac.createDefaultUsers()
	
	// Load users (avoid taking address of range variable)
	for i := range config.Users {
		u := config.Users[i]
		rbac.users[u.Username] = &u
	}

	return rbac
}

func (r *RBACManager) Authenticate(username, password string) (*types.User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Simple password check (in production, use proper password hashing)
	if user.Password != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	user.LastSeen = time.Now()
	return user, nil
}

func (r *RBACManager) Authorize(user *types.User, resource, action, scope string) bool {
	// All users have full access
	return true
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

// RegisterUser creates a new user account
func (r *RBACManager) RegisterUser(username, password, email string) error {
    if username == "" || password == "" {
        return fmt.Errorf("username and password are required")
    }
    if _, exists := r.users[username]; exists {
        return fmt.Errorf("user already exists")
    }
    user := &types.User{
        ID:       generateID(),
        Username: username,
        Password: password, // In production, hash this
        Email:    email,
        Roles:    []string{"user"},
        Created:  time.Now(),
        LastSeen: time.Now(),
        Enabled:  true,
    }
    r.users[username] = user
    return nil
}

// AddUser registers a new user. Password handling is a stub for now.
func (r *RBACManager) AddUser(username, password string, roles []string) error {
    if username == "" {
        return fmt.Errorf("username is required")
    }
    if _, exists := r.users[username]; exists {
        return fmt.Errorf("user %s already exists", username)
    }
    user := &types.User{
        ID:       generateID(),
        Username: username,
        Email:    "",
        Roles:    roles,
        Created:  time.Now(),
        LastSeen: time.Now(),
    }
    r.users[username] = user
    return nil
}

// RemoveUser deletes a user by username.
func (r *RBACManager) RemoveUser(username string) error {
    if _, exists := r.users[username]; !exists {
        return fmt.Errorf("user %s not found", username)
    }
    delete(r.users, username)
    return nil
}

// GetUsers returns all users currently known to the RBAC manager.
func (r *RBACManager) GetUsers() []*types.User {
    out := make([]*types.User, 0, len(r.users))
    for _, u := range r.users {
        out = append(out, u)
    }
    return out
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

	// Safe extraction to avoid panics on malformed tokens
	id, _ := claims["sub"].(string)
	username, _ := claims["user"].(string)
	email, _ := claims["email"].(string)
	roles := interfaceToStringSlice(claims["roles"])
	if id == "" || username == "" {
		return nil, fmt.Errorf("invalid token claims: missing sub/user")
	}

	user := &types.User{
		ID:       id,
		Username: username,
		Email:    email,
		Roles:    roles,
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

func (r *RBACManager) createDefaultRoles() {
	// Admin role with full access
	adminRole := &types.Role{
		Name: "admin",
		Permissions: []types.Permission{
			{Resource: "*", Actions: []string{"*"}, Scope: "*"},
		},
	}
	r.roles["admin"] = adminRole
	
	// Operator role with process management
	operatorRole := &types.Role{
		Name: "operator",
		Permissions: []types.Permission{
			{Resource: "process", Actions: []string{"read", "write", "execute"}, Scope: "*"},
			{Resource: "metrics", Actions: []string{"read"}, Scope: "*"},
			{Resource: "cluster", Actions: []string{"read"}, Scope: "*"},
		},
	}
	r.roles["operator"] = operatorRole
	
	// Viewer role with read-only access
	viewerRole := &types.Role{
		Name: "viewer",
		Permissions: []types.Permission{
			{Resource: "*", Actions: []string{"read"}, Scope: "*"},
		},
	}
	r.roles["viewer"] = viewerRole
}

func (r *RBACManager) createDefaultUsers() {
	// No default users - users must register
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}