package security

import (
	"fmt"
	"strings"
	"time"
)

type HierarchicalRBAC struct {
	organizations map[string]*Organization
	users         map[string]*HierarchicalUser
	permissions   *PermissionEngine
}

type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Teams       []*Team   `json:"teams"`
	Admins      []string  `json:"admins"`
	Created     time.Time `json:"created"`
	Settings    *OrgSettings `json:"settings"`
}

type Team struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	OrgID        string     `json:"org_id"`
	Projects     []*Project `json:"projects"`
	Members      []string   `json:"members"`
	Leads        []string   `json:"leads"`
	Created      time.Time  `json:"created"`
}

type Project struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	TeamID      string      `json:"team_id"`
	Processes   []*Process  `json:"processes"`
	Maintainers []string    `json:"maintainers"`
	Created     time.Time   `json:"created"`
	Environment string      `json:"environment"`
}

type Process struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ProjectID string    `json:"project_id"`
	Owners    []string  `json:"owners"`
	Created   time.Time `json:"created"`
}

type HierarchicalUser struct {
	ID           string                 `json:"id"`
	Username     string                 `json:"username"`
	Email        string                 `json:"email"`
	OrgID        string                 `json:"org_id"`
	Roles        map[string][]string    `json:"roles"` // scope -> roles
	Permissions  []HierarchicalPermission `json:"permissions"`
	Created      time.Time              `json:"created"`
	LastSeen     time.Time              `json:"last_seen"`
}

type HierarchicalPermission struct {
	Resource    string   `json:"resource"`
	Actions     []string `json:"actions"`
	Scope       string   `json:"scope"`
	Inherited   bool     `json:"inherited"`
	GrantedBy   string   `json:"granted_by"`
}

type OrgSettings struct {
	DefaultRoles      []string `json:"default_roles"`
	RequireApproval   bool     `json:"require_approval"`
	MaxTeams          int      `json:"max_teams"`
	MaxProjects       int      `json:"max_projects"`
}

type PermissionEngine struct {
	roleHierarchy map[string][]string // role -> inherited roles
	rolePermissions map[string][]HierarchicalPermission
}

func NewHierarchicalRBAC() *HierarchicalRBAC {
	rbac := &HierarchicalRBAC{
		organizations: make(map[string]*Organization),
		users:         make(map[string]*HierarchicalUser),
		permissions:   NewPermissionEngine(),
	}
	
	// Initialize default role hierarchy
	rbac.permissions.SetupDefaultHierarchy()
	
	return rbac
}

func NewPermissionEngine() *PermissionEngine {
	return &PermissionEngine{
		roleHierarchy:   make(map[string][]string),
		rolePermissions: make(map[string][]HierarchicalPermission),
	}
}

func (p *PermissionEngine) SetupDefaultHierarchy() {
	// Define role hierarchy: higher roles inherit lower role permissions
	p.roleHierarchy["org_admin"] = []string{"team_lead", "project_maintainer", "process_owner", "viewer"}
	p.roleHierarchy["team_lead"] = []string{"project_maintainer", "process_owner", "viewer"}
	p.roleHierarchy["project_maintainer"] = []string{"process_owner", "viewer"}
	p.roleHierarchy["process_owner"] = []string{"viewer"}
	
	// Define permissions for each role
	p.rolePermissions["org_admin"] = []HierarchicalPermission{
		{Resource: "*", Actions: []string{"*"}, Scope: "org"},
	}
	
	p.rolePermissions["team_lead"] = []HierarchicalPermission{
		{Resource: "team", Actions: []string{"read", "write", "manage"}, Scope: "team"},
		{Resource: "project", Actions: []string{"read", "write", "create", "delete"}, Scope: "team"},
		{Resource: "process", Actions: []string{"read", "write", "start", "stop", "restart"}, Scope: "team"},
	}
	
	p.rolePermissions["project_maintainer"] = []HierarchicalPermission{
		{Resource: "project", Actions: []string{"read", "write", "manage"}, Scope: "project"},
		{Resource: "process", Actions: []string{"read", "write", "start", "stop", "restart", "create", "delete"}, Scope: "project"},
	}
	
	p.rolePermissions["process_owner"] = []HierarchicalPermission{
		{Resource: "process", Actions: []string{"read", "write", "start", "stop", "restart"}, Scope: "process"},
	}
	
	p.rolePermissions["viewer"] = []HierarchicalPermission{
		{Resource: "*", Actions: []string{"read"}, Scope: "*"},
	}
}

func (h *HierarchicalRBAC) CreateOrganization(name string, adminUserID string) (*Organization, error) {
	orgID := fmt.Sprintf("org-%d", time.Now().Unix())
	
	org := &Organization{
		ID:      orgID,
		Name:    name,
		Teams:   []*Team{},
		Admins:  []string{adminUserID},
		Created: time.Now(),
		Settings: &OrgSettings{
			DefaultRoles:    []string{"viewer"},
			RequireApproval: false,
			MaxTeams:        50,
			MaxProjects:     200,
		},
	}
	
	h.organizations[orgID] = org
	
	// Grant org admin role to creator
	if user, exists := h.users[adminUserID]; exists {
		user.OrgID = orgID
		h.GrantRole(adminUserID, "org_admin", fmt.Sprintf("org:%s", orgID))
	}
	
	return org, nil
}

func (h *HierarchicalRBAC) CreateTeam(orgID, name string, leadUserID string) (*Team, error) {
	org, exists := h.organizations[orgID]
	if !exists {
		return nil, fmt.Errorf("organization not found: %s", orgID)
	}
	
	teamID := fmt.Sprintf("team-%d", time.Now().Unix())
	
	team := &Team{
		ID:       teamID,
		Name:     name,
		OrgID:    orgID,
		Projects: []*Project{},
		Members:  []string{leadUserID},
		Leads:    []string{leadUserID},
		Created:  time.Now(),
	}
	
	org.Teams = append(org.Teams, team)
	
	// Grant team lead role
	h.GrantRole(leadUserID, "team_lead", fmt.Sprintf("team:%s", teamID))
	
	return team, nil
}

func (h *HierarchicalRBAC) CreateProject(teamID, name, environment string, maintainerUserID string) (*Project, error) {
	team := h.findTeam(teamID)
	if team == nil {
		return nil, fmt.Errorf("team not found: %s", teamID)
	}
	
	projectID := fmt.Sprintf("proj-%d", time.Now().Unix())
	
	project := &Project{
		ID:          projectID,
		Name:        name,
		TeamID:      teamID,
		Processes:   []*Process{},
		Maintainers: []string{maintainerUserID},
		Created:     time.Now(),
		Environment: environment,
	}
	
	team.Projects = append(team.Projects, project)
	
	// Grant project maintainer role
	h.GrantRole(maintainerUserID, "project_maintainer", fmt.Sprintf("project:%s", projectID))
	
	return project, nil
}

func (h *HierarchicalRBAC) CreateProcess(projectID, name string, ownerUserID string) (*Process, error) {
	project := h.findProject(projectID)
	if project == nil {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}
	
	processID := fmt.Sprintf("proc-%d", time.Now().Unix())
	
	process := &Process{
		ID:        processID,
		Name:      name,
		ProjectID: projectID,
		Owners:    []string{ownerUserID},
		Created:   time.Now(),
	}
	
	project.Processes = append(project.Processes, process)
	
	// Grant process owner role
	h.GrantRole(ownerUserID, "process_owner", fmt.Sprintf("process:%s", processID))
	
	return process, nil
}

func (h *HierarchicalRBAC) GrantRole(userID, role, scope string) error {
	user, exists := h.users[userID]
	if !exists {
		return fmt.Errorf("user not found: %s", userID)
	}
	
	if user.Roles == nil {
		user.Roles = make(map[string][]string)
	}
	
	// Add role to scope
	user.Roles[scope] = append(user.Roles[scope], role)
	
	// Calculate effective permissions
	h.calculateEffectivePermissions(user)
	
	return nil
}

func (h *HierarchicalRBAC) calculateEffectivePermissions(user *HierarchicalUser) {
	user.Permissions = []HierarchicalPermission{}
	
	// Collect permissions from all roles and scopes
	for scope, roles := range user.Roles {
		for _, role := range roles {
			// Get direct permissions
			if perms, exists := h.permissions.rolePermissions[role]; exists {
				for _, perm := range perms {
					effectivePerm := perm
					effectivePerm.Scope = scope
					effectivePerm.GrantedBy = role
					user.Permissions = append(user.Permissions, effectivePerm)
				}
			}
			
			// Get inherited permissions
			if inheritedRoles, exists := h.permissions.roleHierarchy[role]; exists {
				for _, inheritedRole := range inheritedRoles {
					if inheritedPerms, exists := h.permissions.rolePermissions[inheritedRole]; exists {
						for _, perm := range inheritedPerms {
							effectivePerm := perm
							effectivePerm.Scope = scope
							effectivePerm.Inherited = true
							effectivePerm.GrantedBy = fmt.Sprintf("%s (via %s)", inheritedRole, role)
							user.Permissions = append(user.Permissions, effectivePerm)
						}
					}
				}
			}
		}
	}
}

func (h *HierarchicalRBAC) CheckPermission(userID, resource, action, targetScope string) bool {
	user, exists := h.users[userID]
	if !exists {
		return false
	}
	
	// Check each permission
	for _, perm := range user.Permissions {
		if h.matchesPermission(perm, resource, action, targetScope) {
			return true
		}
	}
	
	return false
}

func (h *HierarchicalRBAC) matchesPermission(perm HierarchicalPermission, resource, action, targetScope string) bool {
	// Check resource match
	if perm.Resource != "*" && perm.Resource != resource {
		return false
	}
	
	// Check action match
	actionMatch := false
	for _, allowedAction := range perm.Actions {
		if allowedAction == "*" || allowedAction == action {
			actionMatch = true
			break
		}
	}
	if !actionMatch {
		return false
	}
	
	// Check scope match with hierarchy
	return h.scopeMatches(perm.Scope, targetScope)
}

func (h *HierarchicalRBAC) scopeMatches(permScope, targetScope string) bool {
	if permScope == "*" {
		return true
	}
	
	// Extract scope type and ID
	permParts := strings.SplitN(permScope, ":", 2)
	targetParts := strings.SplitN(targetScope, ":", 2)
	
	if len(permParts) != 2 || len(targetParts) != 2 {
		return permScope == targetScope
	}
	
	permType, permID := permParts[0], permParts[1]
	targetType, targetID := targetParts[0], targetParts[1]
	
	// Exact match
	if permScope == targetScope {
		return true
	}
	
	// Hierarchical match: higher scope can access lower scope
	switch permType {
	case "org":
		return h.orgContains(permID, targetType, targetID)
	case "team":
		return h.teamContains(permID, targetType, targetID)
	case "project":
		return h.projectContains(permID, targetType, targetID)
	}
	
	return false
}

func (h *HierarchicalRBAC) orgContains(orgID, targetType, targetID string) bool {
	switch targetType {
	case "org":
		return orgID == targetID
	case "team":
		team := h.findTeam(targetID)
		return team != nil && team.OrgID == orgID
	case "project":
		project := h.findProject(targetID)
		if project == nil {
			return false
		}
		team := h.findTeam(project.TeamID)
		return team != nil && team.OrgID == orgID
	case "process":
		process := h.findProcess(targetID)
		if process == nil {
			return false
		}
		project := h.findProject(process.ProjectID)
		if project == nil {
			return false
		}
		team := h.findTeam(project.TeamID)
		return team != nil && team.OrgID == orgID
	}
	return false
}

func (h *HierarchicalRBAC) teamContains(teamID, targetType, targetID string) bool {
	switch targetType {
	case "team":
		return teamID == targetID
	case "project":
		project := h.findProject(targetID)
		return project != nil && project.TeamID == teamID
	case "process":
		process := h.findProcess(targetID)
		if process == nil {
			return false
		}
		project := h.findProject(process.ProjectID)
		return project != nil && project.TeamID == teamID
	}
	return false
}

func (h *HierarchicalRBAC) projectContains(projectID, targetType, targetID string) bool {
	switch targetType {
	case "project":
		return projectID == targetID
	case "process":
		process := h.findProcess(targetID)
		return process != nil && process.ProjectID == projectID
	}
	return false
}

func (h *HierarchicalRBAC) findTeam(teamID string) *Team {
	for _, org := range h.organizations {
		for _, team := range org.Teams {
			if team.ID == teamID {
				return team
			}
		}
	}
	return nil
}

func (h *HierarchicalRBAC) findProject(projectID string) *Project {
	for _, org := range h.organizations {
		for _, team := range org.Teams {
			for _, project := range team.Projects {
				if project.ID == projectID {
					return project
				}
			}
		}
	}
	return nil
}

func (h *HierarchicalRBAC) findProcess(processID string) *Process {
	for _, org := range h.organizations {
		for _, team := range org.Teams {
			for _, project := range team.Projects {
				for _, process := range project.Processes {
					if process.ID == processID {
						return process
					}
				}
			}
		}
	}
	return nil
}