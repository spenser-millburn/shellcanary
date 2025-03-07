package main

import (
	"fmt"
	
	"github.com/shellcanary/discover/lib"
)

func main() {
	// Create a new discover instance
	d := discover.New()
	
	// Capture current system state (Docker, Kubernetes, systemd)
	fmt.Println("Capturing system state...")
	if err := d.CaptureSystemState(); err != nil {
		fmt.Printf("Error capturing state: %v\n", err)
		return
	}
	
	// Display discovered resources
	fmt.Println("\n=== System Resources ===\n")
	
	// Docker projects
	fmt.Println("=== Docker Projects ===")
	dockerProjects := d.GetDockerProjects()
	if len(dockerProjects) == 0 {
		fmt.Println("No Docker projects found")
	}
	
	for _, project := range dockerProjects {
		fmt.Printf("* %s (%d containers)\n", project.Name, project.Containers)
		for _, container := range project.ContainerDetails {
			fmt.Printf("  - %s: %s\n", container.Name, container.Status)
		}
	}
	
	// Kubernetes resources
	fmt.Println("\n=== Kubernetes Configurations ===")
	k8sConfigs := d.GetKubernetesConfigs()
	if len(k8sConfigs) == 0 {
		fmt.Println("No Kubernetes configurations found")
	}
	
	for _, config := range k8sConfigs {
		fmt.Printf("* %s (%s)\n", config.Name, config.Status)
		for _, ns := range config.Namespaces {
			fmt.Printf("  Namespace: %s\n", ns.Name)
			for _, deployment := range ns.Deployments {
				fmt.Printf("    - %s: %s (%d/%d ready)\n", 
					deployment.Name, deployment.Status, deployment.Ready, deployment.Replicas)
			}
		}
	}
	
	// Systemd services
	fmt.Println("\n=== Active Systemd Services ===")
	systemdServices := d.GetSystemdServices()
	activeCount := 0
	
	for _, service := range systemdServices {
		if service.Status == "active" {
			fmt.Printf("* %s: %s (%s)\n", service.Name, service.Status, service.Description)
			activeCount++
			
			// Limit to 10 services to avoid overwhelming output
			if activeCount >= 10 {
				remaining := 0
				for _, s := range systemdServices {
					if s.Status == "active" {
						remaining++
					}
				}
				fmt.Printf("... and %d more active services\n", remaining-10)
				break
			}
		}
	}
	
	// Save state to file
	fmt.Println("\nSaving system state to file...")
	if err := d.SaveStateToFile(); err != nil {
		fmt.Printf("Error saving state: %v\n", err)
	} else {
		fmt.Println("State saved successfully")
	}
}