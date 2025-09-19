package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gproc/pkg/types"
)

// Phase 1: Zero-downtime reloads, Dependencies, Config wizard, TUI, Snapshots

func reloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reload <name>",
		Short: "Zero-downtime reload of a process",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.ZeroDowntimeReload(args[0]); err != nil {
				fmt.Printf("Error reloading process: %v\n", err)
				return
			}
			fmt.Printf("Reloaded process %s with zero downtime\n", args[0])
		},
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize GProc configuration wizard",
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.ConfigWizard(); err != nil {
				fmt.Printf("Error running config wizard: %v\n", err)
				return
			}
			fmt.Println("Configuration wizard completed!")
		},
	}
}

func topCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "top",
		Short: "Interactive TUI dashboard (like htop)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.StartTUI(); err != nil {
				fmt.Printf("Error starting TUI: %v\n", err)
				return
			}
		},
	}
}

func snapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot <action> [name]",
		Short: "Create, list, or restore process snapshots",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "create":
				if len(args) < 2 {
					fmt.Println("Usage: snapshot create <name>")
					return
				}
				if err := manager.CreateSnapshot(args[1]); err != nil {
					fmt.Printf("Error creating snapshot: %v\n", err)
					return
				}
				fmt.Printf("Created snapshot %s\n", args[1])
				
			case "list":
				snapshots := manager.ListSnapshots()
				if len(snapshots) == 0 {
					fmt.Println("No snapshots found")
					return
				}
				fmt.Println("Available snapshots:")
				for _, s := range snapshots {
					fmt.Printf("  %s - %s (%d processes)\n", s.ID, s.Timestamp.Format("2006-01-02 15:04:05"), len(s.Processes))
				}
				
			case "restore":
				if len(args) < 2 {
					fmt.Println("Usage: snapshot restore <name>")
					return
				}
				if err := manager.RestoreSnapshot(args[1]); err != nil {
					fmt.Printf("Error restoring snapshot: %v\n", err)
					return
				}
				fmt.Printf("Restored snapshot %s\n", args[1])
				
			default:
				fmt.Println("Usage: snapshot <create|list|restore> [name]")
			}
		},
	}
	
	return cmd
}

func dependsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "depends <process> <dependency>",
		Short: "Add dependency (process B starts only if A is healthy)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			processName := args[0]
			dependency := args[1]
			
			if err := manager.AddDependency(processName, dependency); err != nil {
				fmt.Printf("Error adding dependency: %v\n", err)
				return
			}
			fmt.Printf("Added dependency: %s depends on %s\n", processName, dependency)
		},
	}
}

func blueGreenCmd() *cobra.Command {
	var bluePort, greenPort int
	var healthPath string
	
	cmd := &cobra.Command{
		Use:   "blue-green <action> <process>",
		Short: "Blue/Green deployment management",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			processName := args[1]
			
			switch action {
			case "setup":
				config := &types.BlueGreenConfig{
					Enabled:    true,
					ActiveSlot: "blue",
					BluePort:   bluePort,
					GreenPort:  greenPort,
					HealthPath: healthPath,
				}
				if err := manager.SetupBlueGreen(processName, config); err != nil {
					fmt.Printf("Error setting up blue/green: %v\n", err)
					return
				}
				fmt.Printf("Blue/Green deployment setup for %s\n", processName)
				
			case "switch":
				if err := manager.SwitchBlueGreen(processName); err != nil {
					fmt.Printf("Error switching blue/green: %v\n", err)
					return
				}
				fmt.Printf("Switched blue/green deployment for %s\n", processName)
				
			case "status":
				status, err := manager.BlueGreenStatus(processName)
				if err != nil {
					fmt.Printf("Error getting blue/green status: %v\n", err)
					return
				}
				fmt.Printf("Blue/Green status for %s: %s\n", processName, status)
				
			default:
				fmt.Println("Usage: blue-green <setup|switch|status> <process>")
			}
		},
	}
	
	cmd.Flags().IntVar(&bluePort, "blue-port", 8080, "Blue slot port")
	cmd.Flags().IntVar(&greenPort, "green-port", 8081, "Green slot port")
	cmd.Flags().StringVar(&healthPath, "health-path", "/health", "Health check path")
	
	return cmd
}