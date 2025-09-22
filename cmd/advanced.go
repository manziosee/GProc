package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gproc/pkg/types"
)

func startGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start-group <group> <processes>",
		Short: "Start a group of processes",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			groupName := args[0]
			processNames := strings.Split(args[1], ",")
			
			for _, name := range processNames {
				name = strings.TrimSpace(name)
				if err := manager.StartByName(name); err != nil {
					fmt.Printf("Error starting %s: %v\n", name, err)
				}
			}
			fmt.Printf("Started group %s\n", groupName)
		},
	}
}

func stopGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop-group <group>",
		Short: "Stop a group of processes",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			groupName := args[0]
			processes := manager.List()
			
			for _, proc := range processes {
				if proc.Group == groupName {
					if err := manager.Stop(proc.ID); err != nil {
						fmt.Printf("Error stopping %s: %v\n", proc.ID, err)
					}
				}
			}
			fmt.Printf("Stopped group %s\n", groupName)
		},
	}
}

// scheduleCmd implemented in enhanced.go

func webCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start web dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting web dashboard on port %d...\n", port)
			if err := manager.StartWebDashboard(port); err != nil {
				fmt.Printf("Error starting web dashboard: %v\n", err)
				return
			}
		},
	}
	cmd.Flags().IntVar(&port, "port", 3000, "Web dashboard port")
	return cmd
}

func templateCmd() *cobra.Command {
	var command, workingDir string
	var envVars []string
	var autoRestart bool
	var maxRestarts int
	
	cmd := &cobra.Command{
		Use:   "template create <name>",
		Short: "Create a process template",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "create" {
				fmt.Println("Usage: template create <name>")
				return
			}
			
			env := make(map[string]string)
			for _, e := range envVars {
				if parts := strings.SplitN(e, "=", 2); len(parts) == 2 {
					env[parts[0]] = parts[1]
				}
			}
			
			template := &types.ProcessTemplate{
				Name:        args[1],
				Command:     command,
				WorkingDir:  workingDir,
				Env:         env,
				AutoRestart: autoRestart,
				MaxRestarts: maxRestarts,
			}
			
			if err := manager.SaveTemplate(template); err != nil {
				fmt.Printf("Error saving template: %v\n", err)
				return
			}
			fmt.Printf("Created template %s\n", args[1])
		},
	}
	
	cmd.Flags().StringVar(&command, "command", "", "Command to execute")
	cmd.Flags().StringVar(&workingDir, "cwd", "", "Working directory")
	cmd.Flags().StringSliceVar(&envVars, "env", []string{}, "Environment variables")
	cmd.Flags().BoolVar(&autoRestart, "auto-restart", true, "Auto restart")
	cmd.Flags().IntVar(&maxRestarts, "max-restarts", 5, "Max restarts")
	
	return cmd
}

func startTemplateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start-template <template> <name>",
		Short: "Start process from template",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			templateName := args[0]
			processName := args[1]
			
			if err := manager.StartFromTemplate(templateName, processName); err != nil {
				fmt.Printf("Error starting from template: %v\n", err)
				return
			}
			fmt.Printf("Started %s from template %s\n", processName, templateName)
		},
	}
}

func clusterCmd() *cobra.Command {
	var instances int
	var port int
	
	cmd := &cobra.Command{
		Use:   "cluster start <name> <command>",
		Short: "Start multiple instances of a process",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "start" {
				fmt.Println("Usage: cluster start <name> <command>")
				return
			}
			
			baseName := args[1]
			command := args[2]
			cmdArgs := args[3:]
			
			for i := 0; i < instances; i++ {
				processName := fmt.Sprintf("%s-%d", baseName, i)
				processPort := port + i
				
				env := map[string]string{
					"PORT": fmt.Sprintf("%d", processPort),
				}
				
				proc := &types.Process{
					ID:          processName,
					Name:        processName,
					Command:     command,
					Args:        cmdArgs,
					Env:         env,
					AutoRestart: true,
					MaxRestarts: 5,
				}
				
				if err := manager.Start(proc); err != nil {
					fmt.Printf("Error starting %s: %v\n", processName, err)
				}
			}
			
			fmt.Printf("Started cluster %s with %d instances\n", baseName, instances)
		},
	}
	
	cmd.Flags().IntVar(&instances, "instances", 4, "Number of instances")
	cmd.Flags().IntVar(&port, "port", 8080, "Base port number")
	
	return cmd
}

func startFromConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start-from-config <config-file>",
		Short: "Start processes from configuration file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configFile := args[0]
			
			if err := manager.StartFromConfig(configFile); err != nil {
				fmt.Printf("Error starting from config: %v\n", err)
				return
			}
			fmt.Printf("Started processes from %s\n", configFile)
		},
	}
}