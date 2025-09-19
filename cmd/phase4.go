package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Phase 4: Cloud & Container Integration

func dockerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker <action> [options]",
		Short: "Docker/Podman container management",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "run":
				if len(args) < 3 {
					fmt.Println("Usage: docker run <name> <image> [args...]")
					return
				}
				name := args[1]
				image := args[2]
				dockerArgs := args[3:]
				
				if err := manager.RunDockerContainer(name, image, dockerArgs); err != nil {
					fmt.Printf("Error running container: %v\n", err)
					return
				}
				fmt.Printf("Started container %s from image %s\n", name, image)
				
			case "stop":
				if len(args) < 2 {
					fmt.Println("Usage: docker stop <name>")
					return
				}
				name := args[1]
				if err := manager.StopDockerContainer(name); err != nil {
					fmt.Printf("Error stopping container: %v\n", err)
					return
				}
				fmt.Printf("Stopped container %s\n", name)
				
			case "list":
				containers := manager.ListDockerContainers()
				if len(containers) == 0 {
					fmt.Println("No containers running")
					return
				}
				fmt.Println("Running containers:")
				for _, container := range containers {
					fmt.Printf("  %s - %s (%s)\n", container.Name, container.Image, container.Status)
				}
				
			case "logs":
				if len(args) < 2 {
					fmt.Println("Usage: docker logs <name>")
					return
				}
				name := args[1]
				if err := manager.DockerContainerLogs(name); err != nil {
					fmt.Printf("Error getting container logs: %v\n", err)
					return
				}
				
			default:
				fmt.Println("Usage: docker <run|stop|list|logs>")
			}
		},
	}
	
	return cmd
}

func k8sCmd() *cobra.Command {
	var namespace string
	var kubeconfig string
	
	cmd := &cobra.Command{
		Use:   "k8s <action> [options]",
		Short: "Kubernetes integration and operator mode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "operator":
				if err := manager.StartK8sOperator(namespace, kubeconfig); err != nil {
					fmt.Printf("Error starting K8s operator: %v\n", err)
					return
				}
				fmt.Println("Kubernetes operator started")
				
			case "deploy":
				if len(args) < 2 {
					fmt.Println("Usage: k8s deploy <manifest-file>")
					return
				}
				manifestFile := args[1]
				if err := manager.DeployToK8s(manifestFile, namespace); err != nil {
					fmt.Printf("Error deploying to K8s: %v\n", err)
					return
				}
				fmt.Printf("Deployed %s to namespace %s\n", manifestFile, namespace)
				
			case "status":
				status := manager.K8sStatus(namespace)
				fmt.Printf("Kubernetes status in namespace %s:\n%s\n", namespace, status)
				
			case "sync":
				if err := manager.SyncWithK8s(namespace); err != nil {
					fmt.Printf("Error syncing with K8s: %v\n", err)
					return
				}
				fmt.Println("Synced with Kubernetes")
				
			default:
				fmt.Println("Usage: k8s <operator|deploy|status|sync>")
			}
		},
	}
	
	cmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Kubeconfig file path")
	
	return cmd
}

func hybridCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hybrid <action> [options]",
		Short: "Hybrid bare-metal + container management",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "setup":
				if err := manager.SetupHybridMode(); err != nil {
					fmt.Printf("Error setting up hybrid mode: %v\n", err)
					return
				}
				fmt.Println("Hybrid mode setup completed")
				
			case "balance":
				if err := manager.BalanceHybridWorkloads(); err != nil {
					fmt.Printf("Error balancing hybrid workloads: %v\n", err)
					return
				}
				fmt.Println("Hybrid workloads balanced")
				
			case "migrate":
				if len(args) < 3 {
					fmt.Println("Usage: hybrid migrate <process> <to-container|to-bare-metal>")
					return
				}
				processName := args[1]
				target := args[2]
				
				if err := manager.MigrateProcess(processName, target); err != nil {
					fmt.Printf("Error migrating process: %v\n", err)
					return
				}
				fmt.Printf("Migrated %s to %s\n", processName, target)
				
			case "status":
				status := manager.HybridStatus()
				fmt.Printf("Hybrid deployment status:\n%s\n", status)
				
			default:
				fmt.Println("Usage: hybrid <setup|balance|migrate|status>")
			}
		},
	}
	
	return cmd
}