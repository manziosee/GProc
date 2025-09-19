package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Phase 2: Monitoring, Observability, Alerts, Metrics

func metricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics <action> [process]",
		Short: "View process metrics and historical data",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "show":
				if len(args) < 2 {
					// Show all processes metrics
					if err := manager.ShowAllMetrics(); err != nil {
						fmt.Printf("Error showing metrics: %v\n", err)
						return
					}
				} else {
					// Show specific process metrics
					processName := args[1]
					if err := manager.ShowProcessMetrics(processName); err != nil {
						fmt.Printf("Error showing metrics for %s: %v\n", processName, err)
						return
					}
				}
				
			case "history":
				if len(args) < 2 {
					fmt.Println("Usage: metrics history <process>")
					return
				}
				processName := args[1]
				if err := manager.ShowMetricsHistory(processName); err != nil {
					fmt.Printf("Error showing metrics history: %v\n", err)
					return
				}
				
			case "export":
				if err := manager.ExportMetrics(); err != nil {
					fmt.Printf("Error exporting metrics: %v\n", err)
					return
				}
				fmt.Println("Metrics exported successfully")
				
			default:
				fmt.Println("Usage: metrics <show|history|export> [process]")
			}
		},
	}
	
	return cmd
}

func alertsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts <action> [options]",
		Short: "Manage alerts and notifications",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "list":
				alerts := manager.ListAlerts()
				if len(alerts) == 0 {
					fmt.Println("No active alerts")
					return
				}
				fmt.Println("Active alerts:")
				for _, alert := range alerts {
					status := "ACTIVE"
					if alert.Acknowledged {
						status = "ACK"
					}
					fmt.Printf("  [%s] %s - %s: %s (%s)\n", 
						status, alert.Severity, alert.ProcessID, alert.Message, 
						alert.Timestamp.Format("15:04:05"))
				}
				
			case "ack":
				if len(args) < 2 {
					fmt.Println("Usage: alerts ack <alert-id>")
					return
				}
				alertID := args[1]
				if err := manager.AcknowledgeAlert(alertID); err != nil {
					fmt.Printf("Error acknowledging alert: %v\n", err)
					return
				}
				fmt.Printf("Acknowledged alert %s\n", alertID)
				
			case "clear":
				if err := manager.ClearAlerts(); err != nil {
					fmt.Printf("Error clearing alerts: %v\n", err)
					return
				}
				fmt.Println("All alerts cleared")
				
			case "config":
				if err := manager.ConfigureAlerting(); err != nil {
					fmt.Printf("Error configuring alerts: %v\n", err)
					return
				}
				fmt.Println("Alert configuration updated")
				
			default:
				fmt.Println("Usage: alerts <list|ack|clear|config>")
			}
		},
	}
	
	return cmd
}

func profileCmd() *cobra.Command {
	var duration string
	var output string
	
	cmd := &cobra.Command{
		Use:   "profile <process>",
		Short: "Profile a Go process (pprof-like)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			processName := args[0]
			
			if err := manager.ProfileProcess(processName, duration, output); err != nil {
				fmt.Printf("Error profiling process: %v\n", err)
				return
			}
			fmt.Printf("Profiling completed for %s, output saved to %s\n", processName, output)
		},
	}
	
	cmd.Flags().StringVar(&duration, "duration", "30s", "Profiling duration")
	cmd.Flags().StringVar(&output, "output", "", "Output file (default: process-name.prof)")
	
	return cmd
}

func dashboardCmd() *cobra.Command {
	var customConfig string
	
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Start enhanced dashboard with charts and metrics",
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.StartEnhancedDashboard(customConfig); err != nil {
				fmt.Printf("Error starting enhanced dashboard: %v\n", err)
				return
			}
		},
	}
	
	cmd.Flags().StringVar(&customConfig, "config", "", "Custom dashboard configuration")
	
	return cmd
}