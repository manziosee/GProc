package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func daemonCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run GProc as a daemon service",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting GProc daemon...")
			
			// Handle shutdown signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			
			// Keep daemon running
			go func() {
				<-sigChan
				fmt.Println("Shutting down daemon...")
				
				// Stop all processes gracefully
				processes := manager.List()
				for _, proc := range processes {
					if proc.Status == "running" {
						manager.Stop(proc.ID)
					}
				}
				os.Exit(0)
			}()
			
			fmt.Println("GProc daemon started. Press Ctrl+C to stop.")
			select {} // Block forever
		},
	}
}