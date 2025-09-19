package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type RBACManager struct {
	users       map[string]*types.User
	roles       map[string]*Role
	sessions    map[string]*Session
	auditLogs   []AuditLog
	auditEnabled bool
}

type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address"`
}

type AuditLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	IPAddress string    `json:"ip_address"`
	Success   bool      `json:"success"`
}

func NewRBACManager() *RBACManager {
	rbac := &RBACManager{
		users:     make(map[string]*types.User),
		roles:     make(map[string]*Role),
		sessions:  make(map[string]*Session),
		auditLogs: []AuditLog{},
		auditEnabled: true,
	}
	
	// Initialize default roles
	rbac.initDefaultRoles()
	
	return rbac
}

func (r *RBACManager) initDefaultRoles() {
	// Admin role - full access
	r.roles["admin"] = &Role{
		Name: "admin",
		Permissions: []string{
			"process.create", "process.read", "process.update", "process.delete",
			"cluster.manage", "user.manage", "system.configure",
		},
	}
	
	// Operator role - process management
	r.roles["operator"] = &Role{
		Name: "operator",
		Permissions: []string{
			"process.create", "process.read", "process.update", "process.delete",
			"process.logs", "process.metrics",
		},
	}
	
	// Viewer role - read-only
	r.roles["viewer"] = &Role{
		Name: "viewer",
		Permissions: []string{
			"process.read", "process.logs", "process.metrics",
		},
	}
}

func (r *RBACManager) AddUser(username, password string, roleNames []string) error {
	// Validate roles exist
	for _, roleName := range roleNames {
		if _, exists := r.roles[roleName]; !exists {
			return fmt.Errorf("role %s does not exist", roleName)
		}
	}
	
	// Hash password
	hashedPassword := r.hashPassword(password)
	
	user := &types.User{
		Username: username,
		Password: hashedPassword,
		Roles:    roleNames,
		Enabled:  true,
	}
	
	r.users[username] = user
	
	r.auditLog("system", "user.create", username, "", true)
	
	return nil
}

func (r *RBACManager) RemoveUser(username string) error {
	if _, exists := r.users[username]; !exists {
		return fmt.Errorf("user %s not found", username)
	}
	
	delete(r.users, username)
	
	// Invalidate all sessions for this user
	for token, session := range r.sessions {
		if session.UserID == username {
			delete(r.sessions, token)
		}
	}
	
	r.auditLog("system", "user.delete", username, "", true)
	
	return nil
}

func (r *RBACManager) Authenticate(username, password string, ipAddress string) (string, error) {
	user, exists := r.users[username]
	if !exists || !user.Enabled {
		r.auditLog(username, "auth.login", "", ipAddress, false)
		return "", fmt.Errorf("authentication failed")
	}
	
	if !r.verifyPassword(password, user.Password) {
		r.auditLog(username, "auth.login", "", ipAddress, false)
		return "", fmt.Errorf("authentication failed")
	}
	
	// Generate session token
	token := r.generateToken()
	session := &Session{
		Token:     token,
		UserID:    username,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IPAddress: ipAddress,
	}
	
	r.sessions[token] = session
	
	r.auditLog(username, "auth.login", "", ipAddress, true)
	
	return token, nil
}

func (r *RBACManager) ValidateToken(token string) (*types.User, error) {
	session, exists := r.sessions[token]
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}
	
	if time.Now().After(session.ExpiresAt) {
		delete(r.sessions, token)
		return nil, fmt.Errorf("token expired")
	}
	
	user, exists := r.users[session.UserID]
	if !exists || !user.Enabled {
		delete(r.sessions, token)
		return nil, fmt.Errorf("user not found or disabled")
	}
	
	return user, nil
}

func (r *RBACManager) HasPermission(username, permission string) bool {
	user, exists := r.users[username]
	if !exists || !user.Enabled {
		return false
	}
	
	for _, roleName := range user.Roles {
		role, exists := r.roles[roleName]
		if !exists {
			continue
		}
		
		for _, perm := range role.Permissions {
			if perm == permission || perm == "*" {
				return true
			}
		}
	}
	
	return false
}

func (r *RBACManager) CreateRole(name string, permissions []string) error {
	r.roles[name] = &Role{
		Name:        name,
		Permissions: permissions,
	}
	
	r.auditLog("system", "role.create", name, "", true)
	
	return nil
}

func (r *RBACManager) DeleteRole(name string) error {
	if _, exists := r.roles[name]; !exists {
		return fmt.Errorf("role %s not found", name)
	}
	
	// Check if any users have this role
	for _, user := range r.users {
		for _, roleName := range user.Roles {
			if roleName == name {
				return fmt.Errorf("cannot delete role %s: still assigned to users", name)
			}
		}
	}
	
	delete(r.roles, name)
	
	r.auditLog("system", "role.delete", name, "", true)
	
	return nil
}

func (r *RBACManager) GetUsers() []*types.User {
	users := make([]*types.User, 0, len(r.users))
	for _, user := range r.users {
		// Don't return password hash
		userCopy := *user
		userCopy.Password = ""
		users = append(users, &userCopy)
	}
	return users
}

func (r *RBACManager) GetRoles() []*Role {
	roles := make([]*Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	return roles
}

func (r *RBACManager) GetAuditLogs(since string, user string, action string) []AuditLog {
	// Filter audit logs based on criteria
	var filtered []AuditLog
	
	for _, log := range r.auditLogs {
		if user != "" && log.User != user {
			continue
		}
		if action != "" && log.Action != action {
			continue
		}
		// Add time filtering logic here if needed
		
		filtered = append(filtered, log)
	}
	
	return filtered
}

func (r *RBACManager) auditLog(user, action, resource, ipAddress string, success bool) {
	if !r.auditEnabled {
		return
	}
	
	log := AuditLog{
		ID:        r.generateToken()[:8],
		Timestamp: time.Now(),
		User:      user,
		Action:    action,
		Resource:  resource,
		IPAddress: ipAddress,
		Success:   success,
	}
	
	r.auditLogs = append(r.auditLogs, log)
	
	// Keep only last 10000 logs
	if len(r.auditLogs) > 10000 {
		r.auditLogs = r.auditLogs[1000:]
	}
}

func (r *RBACManager) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password + "gproc-salt"))
	return hex.EncodeToString(hash[:])
}

func (r *RBACManager) verifyPassword(password, hash string) bool {
	return r.hashPassword(password) == hash
}

func (r *RBACManager) generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}