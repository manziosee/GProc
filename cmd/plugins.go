package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Plugin System for extensibility

func pluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin <action> [options]",
		Short: "Plugin management system",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "install":
				if len(args) < 2 {
					fmt.Println("Usage: plugin install <plugin-path>")
					return
				}
				pluginPath := args[1]
				if err := manager.InstallPlugin(pluginPath); err != nil {
					fmt.Printf("Error installing plugin: %v\n", err)
					return
				}
				fmt.Printf("Plugin installed from %s\n", pluginPath)
				
			case "list":
				plugins := manager.ListPlugins()
				if len(plugins) == 0 {
					fmt.Println("No plugins installed")
					return
				}
				fmt.Println("Installed plugins:")
				for _, plugin := range plugins {
					status := "enabled"
					if !plugin.Enabled {
						status = "disabled"
					}
					fmt.Printf("  %s - %s (%s)\n", plugin.Name, plugin.Path, status)
					fmt.Printf("    Events: %v\n", plugin.Events)
				}
				
			case "enable":
				if len(args) < 2 {
					fmt.Println("Usage: plugin enable <plugin-name>")
					return
				}
				pluginName := args[1]
				if err := manager.EnablePlugin(pluginName); err != nil {
					fmt.Printf("Error enabling plugin: %v\n", err)
					return
				}
				fmt.Printf("Plugin %s enabled\n", pluginName)
				
			case "disable":
				if len(args) < 2 {
					fmt.Println("Usage: plugin disable <plugin-name>")
					return
				}
				pluginName := args[1]
				if err := manager.DisablePlugin(pluginName); err != nil {
					fmt.Printf("Error disabling plugin: %v\n", err)
					return
				}
				fmt.Printf("Plugin %s disabled\n", pluginName)
				
			case "remove":
				if len(args) < 2 {
					fmt.Println("Usage: plugin remove <plugin-name>")
					return
				}
				pluginName := args[1]
				if err := manager.RemovePlugin(pluginName); err != nil {
					fmt.Printf("Error removing plugin: %v\n", err)
					return
				}
				fmt.Printf("Plugin %s removed\n", pluginName)
				
			case "create":
				if len(args) < 2 {
					fmt.Println("Usage: plugin create <plugin-name>")
					return
				}
				pluginName := args[1]
				if err := manager.CreatePluginTemplate(pluginName); err != nil {
					fmt.Printf("Error creating plugin template: %v\n", err)
					return
				}
				fmt.Printf("Plugin template created for %s\n", pluginName)
				
			default:
				fmt.Println("Usage: plugin <install|list|enable|disable|remove|create>")
			}
		},
	}
	
	return cmd
}

func hookCmd() *cobra.Command {
	var event, script string
	
	cmd := &cobra.Command{
		Use:   "hook <action> [options]",
		Short: "Event hooks management (pre-start, post-stop, on-failure)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action := args[0]
			
			switch action {
			case "add":
				if len(args) < 2 {
					fmt.Println("Usage: hook add <process-name> --event <event> --script <script>")
					return
				}
				processName := args[1]
				if event == "" || script == "" {
					fmt.Println("Both --event and --script are required")
					return
				}
				
				if err := manager.AddHook(processName, event, script); err != nil {
					fmt.Printf("Error adding hook: %v\n", err)
					return
				}
				fmt.Printf("Added %s hook for %s: %s\n", event, processName, script)
				
			case "list":
				if len(args) < 2 {
					fmt.Println("Usage: hook list <process-name>")
					return
				}
				processName := args[1]
				hooks := manager.ListHooks(processName)
				if len(hooks) == 0 {
					fmt.Printf("No hooks found for %s\n", processName)
					return
				}
				fmt.Printf("Hooks for %s:\n", processName)
				for _, hook := range hooks {
					fmt.Printf("  %s: %s\n", hook.Event, hook.Script)
				}
				
			case "remove":
				if len(args) < 2 {
					fmt.Println("Usage: hook remove <process-name> --event <event>")
					return
				}
				processName := args[1]
				if event == "" {
					fmt.Println("--event is required")
					return
				}
				
				if err := manager.RemoveHook(processName, event); err != nil {
					fmt.Printf("Error removing hook: %v\n", err)
					return
				}
				fmt.Printf("Removed %s hook for %s\n", event, processName)
				
			case "test":
				if len(args) < 2 {
					fmt.Println("Usage: hook test <process-name> --event <event>")
					return
				}
				processName := args[1]
				if event == "" {
					fmt.Println("--event is required")
					return
				}
				
				if err := manager.TestHook(processName, event); err != nil {
					fmt.Printf("Error testing hook: %v\n", err)
					return
				}
				fmt.Printf("Tested %s hook for %s\n", event, processName)
				
			default:
				fmt.Println("Usage: hook <add|list|remove|test>")
			}
		},
	}
	
	cmd.Flags().StringVar(&event, "event", "", "Hook event (pre-start, post-start, pre-stop, post-stop, on-failure, on-restart)")
	cmd.Flags().StringVar(&script, "script", "", "Script to execute")
	
	return cmd
}