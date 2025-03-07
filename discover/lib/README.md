# Discover Library

A Go library for discovering and monitoring system resources including Docker containers, Kubernetes deployments, and systemd services.

## Installation

```bash
go get github.com/shellcanary/discover
```

## Features

- Discover Docker Compose projects and containers
- Monitor Kubernetes contexts, namespaces, and deployments
- Track systemd services
- Retrieve logs from various resources
- Persist system state to JSON file

## Usage

```go
package main

import (
	"fmt"
	
	"github.com/shellcanary/discover/lib"
)

func main() {
	// Create a new discover instance
	d := discover.New()
	
	// Capture current system state (Docker, Kubernetes, systemd)
	if err := d.CaptureSystemState(); err != nil {
		fmt.Printf("Error capturing state: %v\n", err)
		return
	}
	
	// Get Docker projects
	dockerProjects := d.GetDockerProjects()
	for _, project := range dockerProjects {
		fmt.Printf("Docker project: %s (%d containers)\n", project.Name, project.Containers)
		
		// Get logs for a specific container
		if len(project.ContainerDetails) > 0 {
			container := project.ContainerDetails[0]
			logs := d.GetDockerLogs(project.Name, container.Name)
			fmt.Printf("Logs for container %s:\n%s\n", container.Name, logs)
		}
	}
	
	// Get Kubernetes configs
	k8sConfigs := d.GetKubernetesConfigs()
	for _, config := range k8sConfigs {
		fmt.Printf("Kubernetes context: %s (%s)\n", config.Name, config.Status)
		
		for _, ns := range config.Namespaces {
			fmt.Printf("  Namespace: %s\n", ns.Name)
			
			for _, deployment := range ns.Deployments {
				fmt.Printf("    Deployment: %s (%s)\n", deployment.Name, deployment.Status)
				
				// Get logs for a deployment
				logs := d.GetKubernetesLogs(config.Name, deployment.Name)
				fmt.Printf("    Logs: %s\n", logs)
			}
		}
	}
	
	// Get systemd services
	systemdServices := d.GetSystemdServices()
	for _, service := range systemdServices {
		fmt.Printf("Systemd service: %s (%s)\n", service.Name, service.Status)
		
		// Get service status
		status, err := d.GetSystemdServiceStatus(service.Name)
		if err == nil {
			fmt.Printf("  Status: %s, Type: %s\n", status.ActiveState, status.Type)
		}
		
		// Get logs for a service
		logs := d.GetSystemdServiceLogs(service.Name)
		fmt.Printf("  Logs: %s\n", logs)
	}
	
	// Save state to file
	if err := d.SaveStateToFile(); err != nil {
		fmt.Printf("Error saving state: %v\n", err)
	}
}
```

## API Reference

### Core Functions

- `New()` - Create a new Discover instance
- `CaptureSystemState()` - Capture current state of all resources
- `LoadStateFromFile()` - Load system state from state file
- `SaveStateToFile()` - Save current state to state file

### Docker Functions

- `GetDockerProjects()` - Get Docker Compose projects
- `GetDockerLogs(projectName, containerName)` - Get logs for a container
- `GetAllDockerProjectLogs(projectName)` - Get logs for all containers in a project

### Kubernetes Functions

- `GetKubernetesConfigs()` - Get Kubernetes contexts and configurations
- `GetKubernetesLogs(contextName, deploymentName)` - Get logs for a deployment

### Systemd Functions

- `GetSystemdServices()` - Get systemd services
- `GetSystemdServiceStatus(serviceName)` - Get detailed service status
- `GetSystemdServiceLogs(serviceName)` - Get logs for a service
- `RestartSystemdService(serviceName)` - Restart a systemd service

## License

[MIT License](LICENSE)