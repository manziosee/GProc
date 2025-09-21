package api

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gproc/internal/process"
	"gproc/internal/security"
	"gproc/pkg/types"
)

// gRPC Service Implementation
type GRPCServer struct {
	config  *types.GRPCConfig
	server  *grpc.Server
	manager *process.Manager
	rbac    *security.RBACManager
	token   *security.TokenManager
}

func NewGRPCServer(config *types.GRPCConfig, manager *process.Manager, rbac *security.RBACManager, token *security.TokenManager) *GRPCServer {
	return &GRPCServer{
		config:  config,
		manager: manager,
		rbac:    rbac,
		token:   token,
	}
}

func (gs *GRPCServer) Start(ctx context.Context) error {
	if !gs.config.Enabled {
		return nil
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", gs.config.Port, err)
	}

	// Create gRPC server with interceptors
	gs.server = grpc.NewServer(
		grpc.UnaryInterceptor(gs.authInterceptor),
		grpc.StreamInterceptor(gs.streamAuthInterceptor),
	)

	// Register services
	RegisterProcessServiceServer(gs.server, gs)
	RegisterClusterServiceServer(gs.server, gs)
	RegisterMetricsServiceServer(gs.server, gs)

	go func() {
		if err := gs.server.Serve(lis); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		}
	}()

	fmt.Printf("gRPC server started on port %d\n", gs.config.Port)
	return nil
}

func (gs *GRPCServer) Stop() {
	if gs.server != nil {
		gs.server.GracefulStop()
	}
}

// Authentication interceptor
func (gs *GRPCServer) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip auth for health check
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	user, err := gs.authenticateFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	// Add user to context
	ctx = context.WithValue(ctx, "user", user)
	return handler(ctx, req)
}

func (gs *GRPCServer) streamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	user, err := gs.authenticateFromContext(ss.Context())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	// Create new context with user
	ctx := context.WithValue(ss.Context(), "user", user)
	wrapped := &wrappedStream{ServerStream: ss, ctx: ctx}
	return handler(srv, wrapped)
}

func (gs *GRPCServer) authenticateFromContext(ctx context.Context) (*types.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, fmt.Errorf("missing authorization token")
	}

	token := tokens[0]
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	return gs.token.ValidateToken(token)
}

// Process Service Implementation
func (gs *GRPCServer) ListProcesses(ctx context.Context, req *ListProcessesRequest) (*ListProcessesResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "process", "read", "*") {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	processes := gs.manager.List()
	var grpcProcesses []*Process

	for _, proc := range processes {
		grpcProcesses = append(grpcProcesses, &Process{
			Id:        proc.ID,
			Name:      proc.Name,
			Command:   proc.Command,
			Args:      proc.Args,
			Status:    string(proc.Status),
			Pid:       int32(proc.PID),
			StartTime: proc.StartTime.Unix(),
			Restarts:  int32(proc.Restarts),
		})
	}

	return &ListProcessesResponse{Processes: grpcProcesses}, nil
}

func (gs *GRPCServer) GetProcess(ctx context.Context, req *GetProcessRequest) (*GetProcessResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "process", "read", fmt.Sprintf("process:%s", req.Id)) {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	proc := gs.manager.Get(req.Id)
	if proc == nil {
		return nil, status.Errorf(codes.NotFound, "process not found")
	}

	return &GetProcessResponse{
		Process: &Process{
			Id:        proc.ID,
			Name:      proc.Name,
			Command:   proc.Command,
			Args:      proc.Args,
			Status:    string(proc.Status),
			Pid:       int32(proc.PID),
			StartTime: proc.StartTime.Unix(),
			Restarts:  int32(proc.Restarts),
		},
	}, nil
}

func (gs *GRPCServer) StartProcess(ctx context.Context, req *StartProcessRequest) (*StartProcessResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "process", "execute", fmt.Sprintf("process:%s", req.Id)) {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	// Start process logic would go here
	return &StartProcessResponse{Success: true, Message: "Process started"}, nil
}

func (gs *GRPCServer) StopProcess(ctx context.Context, req *StopProcessRequest) (*StopProcessResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "process", "execute", fmt.Sprintf("process:%s", req.Id)) {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	if err := gs.manager.Stop(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to stop process: %v", err)
	}

	return &StopProcessResponse{Success: true, Message: "Process stopped"}, nil
}

func (gs *GRPCServer) RestartProcess(ctx context.Context, req *RestartProcessRequest) (*RestartProcessResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "process", "execute", fmt.Sprintf("process:%s", req.Id)) {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	if err := gs.manager.Restart(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to restart process: %v", err)
	}

	return &RestartProcessResponse{Success: true, Message: "Process restarted"}, nil
}

func (gs *GRPCServer) StreamLogs(req *StreamLogsRequest, stream ProcessService_StreamLogsServer) error {
	user := stream.Context().Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "process", "read", fmt.Sprintf("process:%s", req.ProcessId)) {
		return status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	// Simulate log streaming
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			logEntry := &LogEntry{
				Timestamp: time.Now().Unix(),
				Level:     "INFO",
				Message:   fmt.Sprintf("Log message for process %s", req.ProcessId),
			}

			if err := stream.Send(logEntry); err != nil {
				return err
			}
		}
	}
}

// Cluster Service Implementation
func (gs *GRPCServer) GetClusterStatus(ctx context.Context, req *GetClusterStatusRequest) (*GetClusterStatusResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "cluster", "read", "*") {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	return &GetClusterStatusResponse{
		Leader:      "node-1",
		NodeCount:   2,
		Healthy:     true,
		LastElection: time.Now().Add(-1 * time.Hour).Unix(),
	}, nil
}

func (gs *GRPCServer) ListNodes(ctx context.Context, req *ListNodesRequest) (*ListNodesResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "cluster", "read", "*") {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	nodes := []*ClusterNode{
		{Id: "node-1", Address: "localhost:8080", Role: "leader", Status: "active"},
		{Id: "node-2", Address: "localhost:8081", Role: "follower", Status: "active"},
	}

	return &ListNodesResponse{Nodes: nodes}, nil
}

// Metrics Service Implementation
func (gs *GRPCServer) GetMetrics(ctx context.Context, req *GetMetricsRequest) (*GetMetricsResponse, error) {
	user := ctx.Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "metrics", "read", "*") {
		return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	return &GetMetricsResponse{
		ProcessesRunning: 5,
		ProcessesTotal:   8,
		CpuUsage:        34.5,
		MemoryUsage:     67.2,
		Uptime:          "2d 14h 32m",
	}, nil
}

func (gs *GRPCServer) StreamMetrics(req *StreamMetricsRequest, stream MetricsService_StreamMetricsServer) error {
	user := stream.Context().Value("user").(*types.User)
	if !gs.rbac.Authorize(user, "metrics", "read", "*") {
		return status.Errorf(codes.PermissionDenied, "insufficient permissions")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			metrics := &MetricsSnapshot{
				Timestamp:        time.Now().Unix(),
				ProcessesRunning: 5,
				ProcessesTotal:   8,
				CpuUsage:        34.5 + float32(time.Now().Unix()%10),
				MemoryUsage:     67.2 + float32(time.Now().Unix()%5),
			}

			if err := stream.Send(metrics); err != nil {
				return err
			}
		}
	}
}

// Wrapped stream for context injection
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// Protocol Buffer message definitions would typically be in separate .proto files
// For this example, I'm defining them as Go structs

type ListProcessesRequest struct{}

type ListProcessesResponse struct {
	Processes []*Process
}

type Process struct {
	Id        string
	Name      string
	Command   string
	Args      []string
	Status    string
	Pid       int32
	StartTime int64
	Restarts  int32
}

type GetProcessRequest struct {
	Id string
}

type GetProcessResponse struct {
	Process *Process
}

type StartProcessRequest struct {
	Id string
}

type StartProcessResponse struct {
	Success bool
	Message string
}

type StopProcessRequest struct {
	Id string
}

type StopProcessResponse struct {
	Success bool
	Message string
}

type RestartProcessRequest struct {
	Id string
}

type RestartProcessResponse struct {
	Success bool
	Message string
}

type StreamLogsRequest struct {
	ProcessId string
}

type LogEntry struct {
	Timestamp int64
	Level     string
	Message   string
}

type GetClusterStatusRequest struct{}

type GetClusterStatusResponse struct {
	Leader       string
	NodeCount    int32
	Healthy      bool
	LastElection int64
}

type ListNodesRequest struct{}

type ListNodesResponse struct {
	Nodes []*ClusterNode
}

type ClusterNode struct {
	Id      string
	Address string
	Role    string
	Status  string
}

type GetMetricsRequest struct{}

type GetMetricsResponse struct {
	ProcessesRunning int32
	ProcessesTotal   int32
	CpuUsage        float32
	MemoryUsage     float32
	Uptime          string
}

type StreamMetricsRequest struct{}

type MetricsSnapshot struct {
	Timestamp        int64
	ProcessesRunning int32
	ProcessesTotal   int32
	CpuUsage        float32
	MemoryUsage     float32
}

// Service interfaces (would be generated from .proto files)
type ProcessServiceServer interface {
	ListProcesses(context.Context, *ListProcessesRequest) (*ListProcessesResponse, error)
	GetProcess(context.Context, *GetProcessRequest) (*GetProcessResponse, error)
	StartProcess(context.Context, *StartProcessRequest) (*StartProcessResponse, error)
	StopProcess(context.Context, *StopProcessRequest) (*StopProcessResponse, error)
	RestartProcess(context.Context, *RestartProcessRequest) (*RestartProcessResponse, error)
	StreamLogs(*StreamLogsRequest, ProcessService_StreamLogsServer) error
}

type ProcessService_StreamLogsServer interface {
	Send(*LogEntry) error
	grpc.ServerStream
}

type ClusterServiceServer interface {
	GetClusterStatus(context.Context, *GetClusterStatusRequest) (*GetClusterStatusResponse, error)
	ListNodes(context.Context, *ListNodesRequest) (*ListNodesResponse, error)
}

type MetricsServiceServer interface {
	GetMetrics(context.Context, *GetMetricsRequest) (*GetMetricsResponse, error)
	StreamMetrics(*StreamMetricsRequest, MetricsService_StreamMetricsServer) error
}

type MetricsService_StreamMetricsServer interface {
	Send(*MetricsSnapshot) error
	grpc.ServerStream
}

// Registration functions (would be generated)
func RegisterProcessServiceServer(s *grpc.Server, srv ProcessServiceServer) {
	// Implementation would be generated from .proto files
}

func RegisterClusterServiceServer(s *grpc.Server, srv ClusterServiceServer) {
	// Implementation would be generated from .proto files
}

func RegisterMetricsServiceServer(s *grpc.Server, srv MetricsServiceServer) {
	// Implementation would be generated from .proto files
}