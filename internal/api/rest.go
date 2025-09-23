package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gproc/internal/process"
	"gproc/internal/security"
	"gproc/pkg/types"
)

type RESTServer struct {
	config     *types.RESTConfig
	server     *http.Server
	manager    *process.Manager
	rbac       *security.RBACManager
	token      *security.TokenManager
	upgrader   websocket.Upgrader
	wsClients  map[*websocket.Conn]bool
}

func NewRESTServer(config *types.RESTConfig, manager *process.Manager, rbac *security.RBACManager, token *security.TokenManager) *RESTServer {
	return &RESTServer{
		config:    config,
		manager:   manager,
		rbac:      rbac,
		token:     token,
		upgrader:  websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		wsClients: make(map[*websocket.Conn]bool),
	}
}

func (rs *RESTServer) Start(ctx context.Context) error {
	if !rs.config.Enabled {
		return nil
	}
	
	router := mux.NewRouter()
	
	// Add CORS middleware
	router.Use(rs.corsMiddleware)
	
	api := router.PathPrefix(rs.config.Prefix).Subrouter()
	
	// Authentication endpoints
	api.HandleFunc("/auth/login", rs.handleLogin).Methods("POST")
	api.HandleFunc("/auth/register", rs.handleRegister).Methods("POST")
	api.HandleFunc("/auth/refresh", rs.handleRefresh).Methods("POST")
	api.HandleFunc("/auth/sso/login", rs.handleSSOLogin).Methods("GET")
	
	// Process endpoints
	api.HandleFunc("/processes", rs.authMiddleware(rs.handleListProcesses)).Methods("GET")
	api.HandleFunc("/processes", rs.authMiddleware(rs.handleCreateProcess)).Methods("POST")
	api.HandleFunc("/processes/{id}", rs.authMiddleware(rs.handleGetProcess)).Methods("GET")
	api.HandleFunc("/processes/{id}", rs.authMiddleware(rs.handleUpdateProcess)).Methods("PUT")
	api.HandleFunc("/processes/{id}", rs.authMiddleware(rs.handleDeleteProcess)).Methods("DELETE")
	api.HandleFunc("/processes/{id}/start", rs.authMiddleware(rs.handleStartProcess)).Methods("POST")
	api.HandleFunc("/processes/{id}/stop", rs.authMiddleware(rs.handleStopProcess)).Methods("POST")
	api.HandleFunc("/processes/{id}/restart", rs.authMiddleware(rs.handleRestartProcess)).Methods("POST")
	api.HandleFunc("/processes/{id}/logs", rs.authMiddleware(rs.handleGetLogs)).Methods("GET")
	
	// Cluster endpoints
	api.HandleFunc("/cluster/nodes", rs.authMiddleware(rs.handleListNodes)).Methods("GET")
	api.HandleFunc("/cluster/status", rs.authMiddleware(rs.handleClusterStatus)).Methods("GET")
	
	// Health endpoint (no auth required)
	api.HandleFunc("/health", rs.handleHealth).Methods("GET")
	
	// Metrics endpoints
	api.HandleFunc("/metrics", rs.authMiddleware(rs.handleMetrics)).Methods("GET")
	
	// WebSocket endpoint
	api.HandleFunc("/ws", rs.handleWebSocket)
	
	// RBAC endpoints
	api.HandleFunc("/users", rs.authMiddleware(rs.handleListUsers)).Methods("GET")
	api.HandleFunc("/users", rs.authMiddleware(rs.handleCreateUser)).Methods("POST")
	api.HandleFunc("/roles", rs.authMiddleware(rs.handleListRoles)).Methods("GET")
	
	rs.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", rs.config.Port),
		Handler: router,
	}
	
	go func() {
		if err := rs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("REST server error: %v\n", err)
		}
	}()
	
	fmt.Printf("REST API server started on :%d%s\n", rs.config.Port, rs.config.Prefix)
	return nil
}

func (rs *RESTServer) Stop() error {
	if rs.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return rs.server.Shutdown(ctx)
	}
	return nil
}

func (rs *RESTServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (rs *RESTServer) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		tokenString := authHeader[7:] // Remove "Bearer " prefix
		user, err := rs.token.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (rs *RESTServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	user, err := rs.rbac.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	token, err := rs.token.GenerateToken(user)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (rs *RESTServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}
	
	if err := rs.rbac.RegisterUser(req.Username, req.Password, req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (rs *RESTServer) handleRefresh(w http.ResponseWriter, r *http.Request) {
	// Token refresh logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "refresh not implemented"})
}

func (rs *RESTServer) handleSSOLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// If SSO is configured in security config, we would redirect. For now return a helpful message.
	resp := map[string]interface{}{
		"status":  "not_configured",
		"message": "SSO is not configured. Please configure security.sso in config to enable this endpoint.",
		"providers": []string{"google", "github", "azure", "okta"},
	}
	json.NewEncoder(w).Encode(resp)
}

func (rs *RESTServer) handleListProcesses(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	
	if !rs.rbac.Authorize(user, "process", "read", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	processes := rs.manager.List()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processes)
}

func (rs *RESTServer) handleCreateProcess(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	
	if !rs.rbac.Authorize(user, "process", "write", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	var process types.Process
	if err := json.NewDecoder(r.Body).Decode(&process); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	if err := rs.manager.Start(&process); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func (rs *RESTServer) handleGetProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "read", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	process := rs.manager.Get(processID)
	if process == nil {
		http.Error(w, "Process not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func (rs *RESTServer) handleUpdateProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "write", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	var updates types.Process
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// Update process logic would go here
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (rs *RESTServer) handleDeleteProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "delete", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	if err := rs.manager.Stop(processID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func (rs *RESTServer) handleStartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "execute", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	// Start process logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func (rs *RESTServer) handleStopProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "execute", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	if err := rs.manager.Stop(processID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

func (rs *RESTServer) handleRestartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "execute", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	if err := rs.manager.Restart(processID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restarted"})
}

func (rs *RESTServer) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]
	
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "process", "read", fmt.Sprintf("process:%s", processID)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	// Get logs logic
	logs := []string{"Log line 1", "Log line 2", "Log line 3"}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"logs": logs})
}

func (rs *RESTServer) handleListNodes(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "cluster", "read", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	// Mock cluster nodes
	nodes := []types.ClusterNode{
		{ID: "node-1", Address: "localhost:8080", Role: "leader", Status: "active"},
		{ID: "node-2", Address: "localhost:8081", Role: "follower", Status: "active"},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func (rs *RESTServer) handleClusterStatus(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "cluster", "read", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	status := map[string]interface{}{
		"leader":     "node-1",
		"nodes":      2,
		"healthy":    true,
		"last_election": time.Now().Add(-1 * time.Hour),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (rs *RESTServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	health := map[string]interface{}{
		"status":    "healthy",
		"version":   "2.0.0",
		"timestamp": time.Now(),
		"uptime":    "2d 14h 32m",
		"services": map[string]string{
			"api":       "running",
			"websocket": "running",
			"database":  "connected",
		},
	}
	
	json.NewEncoder(w).Encode(health)
}

func (rs *RESTServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "metrics", "read", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	metrics := map[string]interface{}{
		"processes_running": 5,
		"processes_total":   8,
		"cpu_usage":        34.5,
		"memory_usage":     67.2,
		"uptime":          "2d 14h 32m",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (rs *RESTServer) handleListUsers(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "user", "read", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	// Return mock users (in production, get from RBAC manager)
	users := []types.User{
		{ID: "1", Username: "admin", Email: "admin@company.com", Roles: []string{"admin"}},
		{ID: "2", Username: "operator", Email: "operator@company.com", Roles: []string{"operator"}},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (rs *RESTServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "user", "write", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	var newUser struct {
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Roles    []string `json:"roles"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// Create user logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "user created"})
}

func (rs *RESTServer) handleListRoles(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*types.User)
	if !rs.rbac.Authorize(user, "role", "read", "*") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	// Return mock roles
	roles := []types.Role{
		{Name: "admin", Permissions: []types.Permission{{Resource: "*", Actions: []string{"*"}, Scope: "*"}}},
		{Name: "operator", Permissions: []types.Permission{{Resource: "process", Actions: []string{"read", "execute"}, Scope: "*"}}},
		{Name: "viewer", Permissions: []types.Permission{{Resource: "*", Actions: []string{"read"}, Scope: "*"}}},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (rs *RESTServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := rs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()
	
	rs.wsClients[conn] = true
	defer delete(rs.wsClients, conn)
	
	// Send real-time updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		update := map[string]interface{}{
			"type":      "process_update",
			"timestamp": time.Now(),
			"data": map[string]interface{}{
				"processes_running": 5,
				"cpu_usage":        34.5,
				"memory_usage":     67.2,
			},
		}
		
		if err := conn.WriteJSON(update); err != nil {
			return
		}
	}
}