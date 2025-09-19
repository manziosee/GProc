package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Phase 5: Security & Compliance

func rbacCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rbac <action> [options]",
		Short: "Role-Based Access Control management",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "init":
				if err := manager.InitRBAC(); err != nil {
					fmt.Printf("Error initializing RBAC: %v\n", err)
					return
				}
				fmt.Println("RBAC initialized")
				
			case "user":
				if len(args) < 4 {
					fmt.Println("Usage: rbac user <add|remove|list> <username> [password] [roles...]")
					return
				}
				subAction := args[1]
				username := args[2]
				
				switch subAction {
				case "add":
					if len(args) < 4 {
						fmt.Println("Usage: rbac user add <username> <password> [roles...]")
						return
					}
					password := args[3]
					roles := args[4:]
					if err := manager.AddUser(username, password, roles); err != nil {
						fmt.Printf("Error adding user: %v\n", err)
						return
					}
					fmt.Printf("Added user %s with roles: %v\n", username, roles)
					
				case "remove":
					if err := manager.RemoveUser(username); err != nil {
						fmt.Printf("Error removing user: %v\n", err)
						return
					}
					fmt.Printf("Removed user %s\n", username)
					
				case "list":
					users := manager.ListUsers()
					if len(users) == 0 {
						fmt.Println("No users found")
						return
					}
					fmt.Println("Users:")
					for _, user := range users {
						status := "enabled"
						if !user.Enabled {
							status = "disabled"
						}
						fmt.Printf("  %s - %v (%s)\n", user.Username, user.Roles, status)
					}
				}
				
			case "role":
				if len(args) < 3 {
					fmt.Println("Usage: rbac role <create|delete|list> [role-name] [permissions...]")
					return
				}
				subAction := args[1]
				
				switch subAction {
				case "create":
					roleName := args[2]
					permissions := args[3:]
					if err := manager.CreateRole(roleName, permissions); err != nil {
						fmt.Printf("Error creating role: %v\n", err)
						return
					}
					fmt.Printf("Created role %s with permissions: %v\n", roleName, permissions)
					
				case "delete":
					roleName := args[2]
					if err := manager.DeleteRole(roleName); err != nil {
						fmt.Printf("Error deleting role: %v\n", err)
						return
					}
					fmt.Printf("Deleted role %s\n", roleName)
					
				case "list":
					roles := manager.ListRoles()
					if len(roles) == 0 {
						fmt.Println("No roles found")
						return
					}
					fmt.Println("Roles:")
					for _, role := range roles {
						fmt.Printf("  %s - %v\n", role.Name, role.Permissions)
					}
				}
				
			default:
				fmt.Println("Usage: rbac <init|user|role>")
			}
		},
	}
	
	return cmd
}

func auditCmd() *cobra.Command {
	var since string
	var user string
	var action string
	
	cmd := &cobra.Command{
		Use:   "audit <action> [options]",
		Short: "Audit logging and compliance",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdAction := args[0]
			
			switch cmdAction {
			case "enable":
				if err := manager.EnableAuditLogging(); err != nil {
					fmt.Printf("Error enabling audit logging: %v\n", err)
					return
				}
				fmt.Println("Audit logging enabled")
				
			case "disable":
				if err := manager.DisableAuditLogging(); err != nil {
					fmt.Printf("Error disabling audit logging: %v\n", err)
					return
				}
				fmt.Println("Audit logging disabled")
				
			case "logs":
				logs := manager.GetAuditLogs(since, user, action)
				if len(logs) == 0 {
					fmt.Println("No audit logs found")
					return
				}
				fmt.Println("Audit logs:")
				for _, log := range logs {
					fmt.Printf("  %s - %s: %s (%s)\n", 
						log.Timestamp.Format("2006-01-02 15:04:05"), 
						log.User, log.Action, log.Resource)
				}
				
			case "export":
				if err := manager.ExportAuditLogs(); err != nil {
					fmt.Printf("Error exporting audit logs: %v\n", err)
					return
				}
				fmt.Println("Audit logs exported")
				
			default:
				fmt.Println("Usage: audit <enable|disable|logs|export>")
			}
		},
	}
	
	cmd.Flags().StringVar(&since, "since", "", "Show logs since timestamp")
	cmd.Flags().StringVar(&user, "user", "", "Filter by user")
	cmd.Flags().StringVar(&action, "action", "", "Filter by action")
	
	return cmd
}

func secretsCmd() *cobra.Command {
	var vault string
	var path string
	
	cmd := &cobra.Command{
		Use:   "secrets <action> [options]",
		Short: "Secrets management integration",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "init":
				if err := manager.InitSecretsManager(vault); err != nil {
					fmt.Printf("Error initializing secrets manager: %v\n", err)
					return
				}
				fmt.Printf("Secrets manager initialized with %s\n", vault)
				
			case "set":
				if len(args) < 3 {
					fmt.Println("Usage: secrets set <key> <value>")
					return
				}
				key := args[1]
				value := args[2]
				if err := manager.SetSecret(key, value, path); err != nil {
					fmt.Printf("Error setting secret: %v\n", err)
					return
				}
				fmt.Printf("Secret %s set\n", key)
				
			case "get":
				if len(args) < 2 {
					fmt.Println("Usage: secrets get <key>")
					return
				}
				key := args[1]
				value, err := manager.GetSecret(key, path)
				if err != nil {
					fmt.Printf("Error getting secret: %v\n", err)
					return
				}
				fmt.Printf("%s: %s\n", key, value)
				
			case "list":
				secrets := manager.ListSecrets(path)
				if len(secrets) == 0 {
					fmt.Println("No secrets found")
					return
				}
				fmt.Println("Secrets:")
				for _, secret := range secrets {
					fmt.Printf("  %s\n", secret)
				}
				
			default:
				fmt.Println("Usage: secrets <init|set|get|list>")
			}
		},
	}
	
	cmd.Flags().StringVar(&vault, "vault", "hashicorp", "Vault backend (hashicorp/aws/azure)")
	cmd.Flags().StringVar(&path, "path", "", "Secret path")
	
	return cmd
}

func tlsCmd() *cobra.Command {
	var cert, key, ca string
	var generate bool
	
	cmd := &cobra.Command{
		Use:   "tls <action> [options]",
		Short: "TLS/mTLS configuration for secure communication",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "setup":
				if generate {
					if err := manager.GenerateTLSCerts(); err != nil {
						fmt.Printf("Error generating TLS certificates: %v\n", err)
						return
					}
					fmt.Println("TLS certificates generated")
				} else {
					if err := manager.SetupTLS(cert, key, ca); err != nil {
						fmt.Printf("Error setting up TLS: %v\n", err)
						return
					}
					fmt.Println("TLS configured")
				}
				
			case "status":
				status := manager.TLSStatus()
				fmt.Printf("TLS status: %s\n", status)
				
			case "rotate":
				if err := manager.RotateTLSCerts(); err != nil {
					fmt.Printf("Error rotating TLS certificates: %v\n", err)
					return
				}
				fmt.Println("TLS certificates rotated")
				
			default:
				fmt.Println("Usage: tls <setup|status|rotate>")
			}
		},
	}
	
	cmd.Flags().StringVar(&cert, "cert", "", "Certificate file path")
	cmd.Flags().StringVar(&key, "key", "", "Private key file path")
	cmd.Flags().StringVar(&ca, "ca", "", "CA certificate file path")
	cmd.Flags().BoolVar(&generate, "generate", false, "Generate self-signed certificates")
	
	return cmd
}