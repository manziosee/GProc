package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Phase 3: Distributed & Remote Management

func clusterMgmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster-mgmt <action> [options]",
		Short: "Cluster management (master/agent mode)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "init-master":
				if err := manager.InitClusterMaster(); err != nil {
					fmt.Printf("Error initializing cluster master: %v\n", err)
					return
				}
				fmt.Println("Cluster master initialized")
				
			case "join":
				if len(args) < 2 {
					fmt.Println("Usage: cluster-mgmt join <master-address>")
					return
				}
				masterAddr := args[1]
				if err := manager.JoinCluster(masterAddr); err != nil {
					fmt.Printf("Error joining cluster: %v\n", err)
					return
				}
				fmt.Printf("Joined cluster at %s\n", masterAddr)
				
			case "nodes":
				nodes := manager.ListClusterNodes()
				if len(nodes) == 0 {
					fmt.Println("No cluster nodes found")
					return
				}
				fmt.Println("Cluster nodes:")
				for _, node := range nodes {
					fmt.Printf("  %s - %s (%s)\n", node.ID, node.Address, node.Status)
				}
				
			case "leave":
				if err := manager.LeaveCluster(); err != nil {
					fmt.Printf("Error leaving cluster: %v\n", err)
					return
				}
				fmt.Println("Left cluster")
				
			default:
				fmt.Println("Usage: cluster-mgmt <init-master|join|nodes|leave>")
			}
		},
	}
	
	return cmd
}

func remoteCmd() *cobra.Command {
	var remote string
	
	cmd := &cobra.Command{
		Use:   "remote <command> [args...]",
		Short: "Execute commands on remote GProc nodes",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			command := args[0]
			cmdArgs := args[1:]
			
			if err := manager.ExecuteRemoteCommand(remote, command, cmdArgs); err != nil {
				fmt.Printf("Error executing remote command: %v\n", err)
				return
			}
		},
	}
	
	cmd.Flags().StringVar(&remote, "remote", "", "Remote node address")
	
	return cmd
}

func agentCmd() *cobra.Command {
	var masterAddr string
	var port int
	
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Start GProc in agent mode",
		Run: func(cmd *cobra.Command, args []string) {
			if err := manager.StartAgent(masterAddr, port); err != nil {
				fmt.Printf("Error starting agent: %v\n", err)
				return
			}
		},
	}
	
	cmd.Flags().StringVar(&masterAddr, "master", "", "Master node address")
	cmd.Flags().IntVar(&port, "port", 9090, "Agent port")
	
	return cmd
}

func discoveryCmd() *cobra.Command {
	var backend string
	var address string
	
	cmd := &cobra.Command{
		Use:   "discovery <action>",
		Short: "Service discovery integration (Consul/Etcd)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "register":
				if err := manager.RegisterServiceDiscovery(backend, address); err != nil {
					fmt.Printf("Error registering service discovery: %v\n", err)
					return
				}
				fmt.Printf("Registered with %s at %s\n", backend, address)
				
			case "deregister":
				if err := manager.DeregisterServiceDiscovery(); err != nil {
					fmt.Printf("Error deregistering service discovery: %v\n", err)
					return
				}
				fmt.Println("Deregistered from service discovery")
				
			case "status":
				status := manager.ServiceDiscoveryStatus()
				fmt.Printf("Service discovery status: %s\n", status)
				
			default:
				fmt.Println("Usage: discovery <register|deregister|status>")
			}
		},
	}
	
	cmd.Flags().StringVar(&backend, "backend", "consul", "Service discovery backend (consul/etcd)")
	cmd.Flags().StringVar(&address, "address", "localhost:8500", "Service discovery address")
	
	return cmd
}