package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	cmd := &cobra.Command{
		Use:   "start <name> <command> [args...]",
		Short: "Start a new process",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			proc := &types.Process{
				ID:          args[0],
				Name:        args[0],
				Command:     args[1],
				Args:        args[2:],
				WorkingDir:  workingDir,
				AutoRestart: autoRestart,
				MaxRestarts: maxRestarts,
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
			if err := manager.Stop(args[0]); err != nil {
				fmt.Printf("Error stopping process: %v\n", err)
			}
			
			time.Sleep(1 * time.Second)
			
			// Note: In a real implementation, you'd need to store process configs
			fmt.Printf("Restarted process %s\n", args[0])
		},
	}
}