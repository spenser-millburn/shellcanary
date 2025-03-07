package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"discover/agents/docker"
	"discover/agents/kubernetes"
	"discover/agents/systemd"
	"discover/models"
	"discover/state"
	"discover/ui/docker"
	"discover/ui/kubernetes"
	"discover/ui/systemd"
	"discover/ui/help"
)

// StartMainMenu launches the main interactive menu
func StartMainMenu() {
	// Loop through the main menu until user exits
	for {
		resourceTypes := []string{
			"🔍 All Resource Types",
			"🐳 Docker Only",
			"☸️ Kubernetes Only",
			"⚙️ Systemd Only",
			"📊 Capture System State Only",
			"❓ Help",
			"❌ Exit Application",
		}
		
		typePrompt := promptui.Select{
			Label: "Select which types of resources to search for",
			Items: resourceTypes,
		}
		
		typeIndex, typeResult, err := typePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		
		// Handle exit option
		if typeResult == "❌ Exit Application" {
			fmt.Println("Exiting application. Goodbye!")
			return
		}
		
		// Handle help option
		if typeResult == "❓ Help" {
			help.ShowHelpPage()
			PauseForUser()
			continue // Return to main menu
		}
		
		// Handle capture state only option
		if typeResult == "📊 Capture System State Only" {
			if err := CaptureSystemState(); err != nil {
				fmt.Println(err)
			}
			PauseForUser()
			continue // Return to main menu
		}
		
		// Initialize resource collections
		var dockerProjects []models.DockerProject
		var k8sConfigs []models.KubernetesConfig
		var systemdServices []models.SystemdService
		
		// Fetch only the selected resource types
		switch typeIndex {
		case 0: // All Resource Types
			fmt.Println("Searching for all resource types...")
			dockerProjects = docker.GetDockerComposeProjects()
			k8sConfigs = kubernetes.GetKubernetesConfigs()
			systemdServices = systemd.GetSystemdServices()
		case 1: // Docker Only
			fmt.Println("Searching for Docker resources...")
			dockerProjects = docker.GetDockerComposeProjects()
		case 2: // Kubernetes Only
			fmt.Println("Searching for Kubernetes resources...")
			k8sConfigs = kubernetes.GetKubernetesConfigs()
		case 3: // Systemd Only
			fmt.Println("Searching for Systemd resources...")
			systemdServices = systemd.GetSystemdServices()
		}
		
		// Capture the system state with what we've found
		if err := state.UpdateSystemState(
			dockerProjects, 
			k8sConfigs, 
			systemdServices,
		); err != nil {
			fmt.Printf("Warning: Failed to capture system state: %v\n", err)
		}
		
		showResourceSelectionMenu(dockerProjects, k8sConfigs, systemdServices)
	}
}

// showResourceSelectionMenu displays the menu for selecting specific resources
func showResourceSelectionMenu(
	dockerProjects []models.DockerProject,
	k8sConfigs []models.KubernetesConfig,
	systemdServices []models.SystemdService,
) {
	for {
		// Build the selection options based on what we've fetched
		var options []string
		options = append(options, "⬅️ Back to Main Menu")
		options = append(options, "❓ Help")
		
		for _, project := range dockerProjects {
			options = append(options, fmt.Sprintf("🐳 Docker: %s", project.Name))
		}
		
		for _, config := range k8sConfigs {
			options = append(options, fmt.Sprintf("☸️ Kubernetes: %s", config.Name))
		}
		
		for _, service := range systemdServices {
			// Only include active services to avoid cluttering the menu
			if strings.HasPrefix(service.Status, "active") {
				options = append(options, fmt.Sprintf("⚙️ Systemd: %s", service.Name))
			}
		}
		
		// Check if we have any options besides the back and help options
		if len(options) <= 2 {
			fmt.Println("No resources found for the selected type(s).")
			break // Return to main menu
		}
		
		// Create the main selection prompt
		prompt := promptui.Select{
			Label: "📋 Select a project, configuration, or service to view logs",
			Items: options,
		}
		
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		
		// Handle back option
		if result == "⬅️ Back to Main Menu" {
			break // Return to main menu
		}
		
		// Handle help option
		if result == "❓ Help" {
			help.ShowHelpPage()
			PauseForUser()
			continue
		}
		
		// Handle resource selection
		if strings.HasPrefix(result, "🐳 Docker: ") {
			dockerUI.ShowDockerMenu(strings.TrimPrefix(result, "🐳 Docker: "))
		} else if strings.HasPrefix(result, "☸️ Kubernetes: ") {
			kubernetesUI.ShowKubernetesMenu(strings.TrimPrefix(result, "☸️ Kubernetes: "))
		} else if strings.HasPrefix(result, "⚙️ Systemd: ") {
			systemdUI.ShowSystemdMenu(strings.TrimPrefix(result, "⚙️ Systemd: "))
		}
		
		// Pause after displaying content
		PauseForUser()
	}
}