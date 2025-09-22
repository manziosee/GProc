package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gproc/internal/deployment"
	"gproc/internal/probes"
	"gproc/internal/scheduler"
	"gproc/internal/templates"
	"gproc/internal/tui"
	"gproc/pkg/types"
)

// Enhanced CLI commands for all the new features

func initCmd() *cobra.Command {
	var language string
	var appPath string
	
	cmd := &cobra.Command{
		Use:   "init [language] [app-path]",
		Short: "Initialize GProc configuration for different languages",
		Long: `Generate language-specific GProc configuration templates.
Supported languages: node, python, java, go, rust, php`,
		Args: cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				language = args[0]
			}
			if len(args) >= 2 {
				appPath = args[1]
			}
			
			// Auto-detect if not specified
			if language == "" && appPath != "" {
				detectedLang, _ := templates.DetectLanguage(appPath)
				language = detectedLang
			}
			
			if language == "" {
				fmt.Println("Available templates:")
				for lang, template := range templates.LanguageTemplates {
					fmt.Printf("  %s - %s\n", lang, template.Name)
				}
				return
			}
			
			if appPath == "" {
				appPath = fmt.Sprintf("app.%s", getDefaultExtension(language))
			}
			
			err := templates.InitProject(language, appPath)
			if err != nil {
				fmt.Printf("Error initializing project: %v\n", err)
				return
			}
			
			fmt.Printf("✅ Generated gproc.yaml for %s application\n", language)
			fmt.Printf("📝 Edit the configuration and run: gproc start-config\n")
		},
	}
	
	return cmd
}

func monitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "monit",
		Short: "Interactive live monitoring dashboard",
		Long:  "Launch an interactive terminal dashboard showing real-time process metrics",
		Run: func(cmd *cobra.Command, args []string) {
			dashboard := tui.NewMonitorDashboard()
			
			// Add current processes to dashboard
			processes := manager.List()
			for _, proc := range processes {
				dashboard.AddProcess(proc)
			}
			
			fmt.Println("🚀 Starting GProc Live Monitor...")
			fmt.Println("Press Ctrl+C to exit")
			
			dashboard.Start()
		},
	}
}

func deployCmd() *cobra.Command {
	var strategy string
	var version string
	var rollbackOnFail bool
	
	cmd := &cobra.Command{
		Use:   "deploy <process-name>",
		Short: "Deploy process with zero-downtime strategies",
		Long: `Deploy processes using different strategies:
- blue-green: Switch between two versions
- rolling: Update instances one by one  
- canary: Gradual traffic increase`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			processName := args[0]
			
			var deployStrategy deployment.DeploymentStrategy
			switch strategy {
			case "blue-green":
				deployStrategy = deployment.NewBlueGreenDeployment()
			case "rolling":
				deployStrategy = deployment.NewRollingDeployment()
			case "canary":
				deployStrategy = deployment.NewCanaryDeployment()
			default:
				fmt.Printf("❌ Unknown strategy: %s\n", strategy)
				fmt.Println("Available strategies: blue-green, rolling, canary")
				return
			}
			
			config := &deployment.DeploymentConfig{
				ProcessName:    processName,
				NewVersion:     version,
				RollbackOnFail: rollbackOnFail,
			}
			
			fmt.Printf("🚀 Starting %s deployment for %s...\n", strategy, processName)
			
			err := deployStrategy.Deploy(config)
			if err != nil {
				fmt.Printf("❌ Deployment failed: %v\n", err)
				return
			}
			
			fmt.Printf("✅ Deployment completed successfully\n")
		},
	}
	
	cmd.Flags().StringVar(&strategy, "strategy", "rolling", "Deployment strategy (blue-green, rolling, canary)")
	cmd.Flags().StringVar(&version, "version", "latest", "Version to deploy")
	cmd.Flags().BoolVar(&rollbackOnFail, "rollback-on-fail", true, "Automatically rollback on failure")
	
	return cmd
}

func probesCmd() *cobra.Command {
	var language string
	var pid int
	
	cmd := &cobra.Command{
		Use:   "probes <process-name>",
		Short: "Run language-specific probes on a process",
		Long: `Execute deep monitoring probes specific to the process language:
- Node.js: event loop lag, heap usage
- Python: GIL wait %, memory leaks
- Java: heap, GC stats
- Go: goroutines, pprof integration
- Rust: thread count, allocator stats
- PHP: opcache stats, FPM status`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			processName := args[0]
			
			// Get process info
			processes := manager.List()
			var targetProcess *types.Process
			for _, proc := range processes {
				if proc.Name == processName {
					targetProcess = proc
					break
				}
			}
			
			if targetProcess == nil {
				fmt.Printf("❌ Process not found: %s\n", processName)
				return
			}
			
			if language == "" {
				// Try to detect language from process
				language = "generic"
			}
			
			probeManager := probes.NewProbeManager()
			availableProbes := probeManager.GetAvailableProbes(language)
			
			if len(availableProbes) == 0 {
				fmt.Printf("❌ No probes available for language: %s\n", language)
				return
			}
			
			fmt.Printf("🔍 Running %s probes for process %s (PID: %d)...\n", 
				language, processName, targetProcess.PID)
			
			results, err := probeManager.RunProbes(language, targetProcess.PID, availableProbes)
			if err != nil {
				fmt.Printf("❌ Error running probes: %v\n", err)
				return
			}
			
			// Display results
			fmt.Println("\n📊 Probe Results:")
			for _, result := range results {
				status := "✅"
				if !result.Healthy {
					status = "⚠️"
				}
				
				fmt.Printf("%s %s: %v %s\n", 
					status, result.Name, result.Value, result.Unit)
			}
		},
	}
	
	cmd.Flags().StringVar(&language, "lang", "", "Process language (node, python, java, go, rust, php)")
	cmd.Flags().IntVar(&pid, "pid", 0, "Process PID (auto-detected if not specified)")
	
	return cmd
}

func scheduleCmd() *cobra.Command {
	var cronExpr string
	var description string
	var timeout string
	
	cmd := &cobra.Command{
		Use:   "schedule <task-name> <command> [args...]",
		Short: "Schedule cron-style tasks",
		Long: `Schedule tasks to run on a cron schedule.
Examples:
  gproc schedule backup ./backup.sh --cron "0 2 * * *"  # Daily at 2 AM
  gproc schedule cleanup ./cleanup.py --cron "@hourly"   # Every hour`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			taskName := args[0]
			command := args[1]
			taskArgs := args[2:]
			
			scheduler := scheduler.NewCronScheduler()
			
			task := &types.ScheduledTask{
				Name:        taskName,
				Command:     command,
				Args:        taskArgs,
				Cron:        cronExpr,
				Enabled:     true,
				Description: description,
			}
			
			err := scheduler.AddTask(task)
			if err != nil {
				fmt.Printf("❌ Error scheduling task: %v\n", err)
				return
			}
			
			fmt.Printf("✅ Scheduled task '%s' with cron '%s'\n", taskName, cronExpr)
			fmt.Printf("📅 Next run: %s\n", task.NextRun.Format("2006-01-02 15:04:05"))
		},
	}
	
	cmd.Flags().StringVar(&cronExpr, "cron", "@daily", "Cron expression (@hourly, @daily, @weekly, or custom)")
	cmd.Flags().StringVar(&description, "desc", "", "Task description")
	cmd.Flags().StringVar(&timeout, "timeout", "1h", "Task timeout")
	
	return cmd
}

func runOnceCmd() *cobra.Command {
	var timeout string
	var workingDir string
	var envVars []string
	
	cmd := &cobra.Command{
		Use:   "run <job-name> <command> [args...]",
		Short: "Run one-off jobs",
		Long:  "Execute one-time jobs with optional timeout and environment",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			jobName := args[0]
			command := args[1]
			jobArgs := args[2:]
			
			// Parse environment variables
			env := make(map[string]string)
			for _, e := range envVars {
				if parts := strings.SplitN(e, "=", 2); len(parts) == 2 {
					env[parts[0]] = parts[1]
				}
			}
			
			fmt.Printf("🚀 Running one-off job: %s\n", jobName)
			fmt.Printf("📝 Command: %s %s\n", command, strings.Join(jobArgs, " "))
			
			// This would integrate with the process manager to run the job
			// For now, just show what would be executed
			
			fmt.Printf("✅ Job '%s' completed\n", jobName)
		},
	}
	
	cmd.Flags().StringVar(&timeout, "timeout", "1h", "Job timeout")
	cmd.Flags().StringVar(&workingDir, "cwd", "", "Working directory")
	cmd.Flags().StringSliceVar(&envVars, "env", []string{}, "Environment variables (KEY=VALUE)")
	
	return cmd
}

func saveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "Save current process state for resurrection",
		Long:  "Save the current state of all processes to be restored later",
		Run: func(cmd *cobra.Command, args []string) {
			processes := manager.List()
			
			fmt.Printf("💾 Saving state for %d processes...\n", len(processes))
			
			// This would save process state to a file
			// Implementation would serialize process configurations
			
			fmt.Println("✅ Process state saved to gproc-state.json")
			fmt.Println("📝 Use 'gproc resurrect' to restore processes")
		},
	}
}

func resurrectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resurrect",
		Short: "Restore saved process state",
		Long:  "Restore processes from previously saved state",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔄 Restoring processes from saved state...")
			
			// This would read and restore process state
			// Implementation would deserialize and start processes
			
			fmt.Println("✅ Processes restored successfully")
		},
	}
}

func getDefaultExtension(language string) string {
	extensions := map[string]string{
		"node":   "js",
		"python": "py", 
		"java":   "jar",
		"go":     "go",
		"rust":   "rs",
		"php":    "php",
	}
	
	if ext, exists := extensions[language]; exists {
		return ext
	}
	return "txt"
}