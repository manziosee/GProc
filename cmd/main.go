package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"gproc/internal/logger"
	"gproc/internal/process"
	"gproc/pkg/types"
)

var manager *process.Manager

func main() {
	logDir := "./logs"
	os.MkdirAll(logDir, 0755)
	manager = process.NewManager(logDir)

	rootCmd := &cobra.Command{
		Use:   "gproc",
		Short: "A process manager for Go applications",
	}

	rootCmd.AddCommand(
		startCmd(),
		stopCmd(),
		listCmd(),
		logsCmd(),
		restartCmd(),
		daemonCmd(),
		startGroupCmd(),
		stopGroupCmd(),
		scheduleCmd(),
		webCmd(),
		templateCmd(),
		startTemplateCmd(),
		clusterCmd(),
		startFromConfigCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startCmd() *cobra.Command {
	var autoRestart bool
	var maxRestarts int
	var workingDir string
	var envVars []string
	var group string
	var healthCheck string
	var healthInterval string
	var logMaxSize string
	var logMaxFiles int
	var memoryLimit string
	var cpuLimit float64
	var notifyEmail string
	var notifySlack string

	cmd := &cobra.Command{
		Use:   "start <name> <command> [args...]",
		Short: "Start a new process",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Parse environment variables
			env := make(map[string]string)
			for _, e := range envVars {
				if parts := strings.SplitN(e, "=", 2); len(parts) == 2 {
					env[parts[0]] = parts[1]
				}
			}
			
			// Parse health check
			var hc *types.HealthCheck
			if healthCheck != "" {
				interval, _ := time.ParseDuration(healthInterval)
				if interval == 0 {
					interval = 30 * time.Second
				}
				hc = &types.HealthCheck{
					URL:      healthCheck,
					Interval: interval,
					Timeout:  5 * time.Second,
					Retries:  3,
				}
			}
			
			// Parse log rotation
			var lr *types.LogRotation
			if logMaxSize != "" {
				lr = &types.LogRotation{
					MaxSize:  logMaxSize,
					MaxFiles: logMaxFiles,
				}
			}
			
			// Parse resource limits
			var rl *types.ResourceLimit
			if memoryLimit != "" || cpuLimit > 0 {
				memMB := 0
				if memoryLimit != "" {
					// Simple parsing: "512MB" -> 512
					if strings.HasSuffix(memoryLimit, "MB") {
						fmt.Sscanf(memoryLimit, "%dMB", &memMB)
					}
				}
				rl = &types.ResourceLimit{
					MemoryMB: memMB,
					CPULimit: cpuLimit,
				}
			}
			
			// Parse notifications
			var notif *types.Notifications
			if notifyEmail != "" || notifySlack != "" {
				notif = &types.Notifications{
					Email: notifyEmail,
					Slack: notifySlack,
				}
			}
			
			proc := &types.Process{
				ID:            args[0],
				Name:          args[0],
				Command:       args[1],
				Args:          args[2:],
				WorkingDir:    workingDir,
				Env:           env,
				Group:         group,
				AutoRestart:   autoRestart,
				MaxRestarts:   maxRestarts,
				HealthCheck:   hc,
				LogRotation:   lr,
				ResourceLimit: rl,
				Notifications: notif,
			}

			if err := manager.Start(proc); err != nil {
				fmt.Printf("Error starting process: %v\n", err)
				return
			}
			fmt.Printf("Started process %s\n", args[0])
		},
	}

	cmd.Flags().BoolVar(&autoRestart, "auto-restart", true, "Auto restart on failure")
	cmd.Flags().IntVar(&maxRestarts, "max-restarts", 5, "Maximum restart attempts")
	cmd.Flags().StringVar(&workingDir, "cwd", "", "Working directory")
	cmd.Flags().StringSliceVar(&envVars, "env", []string{}, "Environment variables (KEY=VALUE)")
	cmd.Flags().StringVar(&group, "group", "", "Process group name")
	cmd.Flags().StringVar(&healthCheck, "health-check", "", "Health check URL")
	cmd.Flags().StringVar(&healthInterval, "health-interval", "30s", "Health check interval")
	cmd.Flags().StringVar(&logMaxSize, "log-max-size", "", "Maximum log file size (e.g., 100MB)")
	cmd.Flags().IntVar(&logMaxFiles, "log-max-files", 5, "Maximum number of log files")
	cmd.Flags().StringVar(&memoryLimit, "memory-limit", "", "Memory limit (e.g., 512MB)")
	cmd.Flags().Float64Var(&cpuLimit, "cpu-limit", 0, "CPU limit percentage (e.g., 50.0)")
	cmd.Flags().StringVar(&notifyEmail, "notify-email", "", "Email for notifications")
	cmd.Flags().StringVar(&notifySlack, "notify-slack", "", "Slack webhook for notifications")

	return cmd
}

func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <name>",
		Short: "Stop a running process",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.Stop(args[0]); err != nil {
				fmt.Printf("Error stopping process: %v\n", err)
				return
			}
			fmt.Printf("Stopped process %s\n", args[0])
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all processes",
		Run: func(cmd *cobra.Command, args []string) {
			processes := manager.List()
			if len(processes) == 0 {
				fmt.Println("No processes running")
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSTATUS\tPID\tRESTARTS\tUPTIME")
			
			for _, proc := range processes {
				uptime := ""
				if proc.Status == types.StatusRunning {
					uptime = time.Since(proc.StartTime).Round(time.Second).String()
				}
				
				fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
					proc.Name, proc.Status, proc.PID, proc.Restarts, uptime)
			}
			w.Flush()
		},
	}
}

func logsCmd() *cobra.Command {
	var lines int
	
	cmd := &cobra.Command{
		Use:   "logs <name>",
		Short: "View process logs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logFile := filepath.Join("./logs", args[0]+".log")
			if err := logger.TailFile(logFile, lines); err != nil {
				fmt.Printf("Error reading logs: %v\n", err)
			}
		},
	}
	
	cmd.Flags().IntVar(&lines, "lines", 20, "Number of lines to show")
	return cmd
}

func restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart <name>",
		Short: "Restart a process",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.Restart(args[0]); err != nil {
				fmt.Printf("Error restarting process: %v\n", err)
				return
			}
			fmt.Printf("Restarted process %s\n", args[0])
		},
	}
}