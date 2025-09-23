package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gproc/internal/api"
	"gproc/internal/security"
	"gproc/pkg/types"
)

func daemonCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run GProc as a daemon service with API server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting GProc daemon with API server...")
			
			// Initialize RBAC and security
			rbacConfig := &types.RBACConfig{
				Enabled: true,
				Roles:   []types.Role{},
				Users:   []types.User{},
			}
			
			jwtConfig := &types.JWTConfig{
				Secret:     "gproc-demo-secret-key-2024",
				Expiration: 24 * time.Hour,
				Issuer:     "gproc",
			}
			
			rbacManager := security.NewRBACManager(rbacConfig)
			tokenManager := security.NewTokenManager(jwtConfig)
			
			// Initialize API server
			apiConfig := &types.RESTConfig{
				Enabled: true,
				Port:    8080,
				Prefix:  "/api/v1",
			}
			
			apiServer := api.NewRESTServer(apiConfig, manager, rbacManager, tokenManager)
			
			// Start API server
			ctx := context.Background()
			if err := apiServer.Start(ctx); err != nil {
				fmt.Printf("Failed to start API server: %v\n", err)
				return
			}
			
			// Handle shutdown signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			
			fmt.Println("GProc daemon started with API server on :8080")
			fmt.Println("API endpoints available at http://localhost:8080/api/v1")
			fmt.Println("Health check: http://localhost:8080/api/v1/health")
			fmt.Println("Press Ctrl+C to stop.")
			
			// Wait for shutdown signal
			<-sigChan
			fmt.Println("Shutting down daemon...")
			
			// Stop API server
			if err := apiServer.Stop(); err != nil {
				fmt.Printf("Error stopping API server: %v\n", err)
			}
			
			// Stop all processes gracefully
			processes := manager.List()
			for _, proc := range processes {
				if proc.Status == "running" {
					manager.Stop(proc.ID)
				}
			}
			
			fmt.Println("Daemon stopped.")
		},
	}
}