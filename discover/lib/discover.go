// Package discover provides functionality to discover and monitor various system resources
// including Docker containers, Kubernetes deployments, and systemd services.
package discover

import (
	"fmt"

	"github.com/shellcanary/discover/lib/agents/docker"
	"github.com/shellcanary/discover/lib/agents/kubernetes"
	"github.com/shellcanary/discover/lib/agents/systemd"
	"github.com/shellcanary/discover/lib/models"
	"github.com/shellcanary/discover/lib/state"
)

// Discover represents the main library API for resource discovery
type Discover struct {
	State models.SystemState
}

// New creates a new Discover instance
func New() *Discover {
	return &Discover{
		State: models.SystemState{},
	}
}

// CaptureSystemState captures the current state of all resources
func (d *Discover) CaptureSystemState() error {
	// Gather data from all agents
	dockerProjects := docker.GetDockerComposeProjects()
	k8sConfigs := kubernetes.GetKubernetesConfigs()
	systemdServices := systemd.GetSystemdServices()
	
	// Update the local state
	d.State.DockerProjects = dockerProjects
	d.State.KubernetesConfigs = k8sConfigs
	d.State.SystemdServices = systemdServices
	
	// Update system state file
	err := state.UpdateSystemState(dockerProjects, k8sConfigs, systemdServices)
	if err != nil {
		return fmt.Errorf("error updating system state: %v", err)
	}
	
	return nil
}

// GetDockerProjects returns Docker compose projects
func (d *Discover) GetDockerProjects() []models.DockerProject {
	return docker.GetDockerComposeProjects()
}

// GetDockerLogs retrieves logs for a specific container in a project
func (d *Discover) GetDockerLogs(projectName, containerName string) string {
	return docker.GetDockerLogs(projectName, containerName)
}

// GetAllDockerProjectLogs retrieves logs for all containers in a project
func (d *Discover) GetAllDockerProjectLogs(projectName string) string {
	return docker.GetAllProjectLogs(projectName)
}

// GetKubernetesConfigs returns Kubernetes configurations
func (d *Discover) GetKubernetesConfigs() []models.KubernetesConfig {
	return kubernetes.GetKubernetesConfigs()
}

// GetKubernetesLogs retrieves logs for a specific deployment
func (d *Discover) GetKubernetesLogs(contextName, deploymentName string) string {
	return kubernetes.GetKubernetesLogs(contextName, deploymentName)
}

// GetSystemdServices returns systemd services
func (d *Discover) GetSystemdServices() []models.SystemdService {
	return systemd.GetSystemdServices()
}

// GetSystemdServiceStatus retrieves detailed status of a service
func (d *Discover) GetSystemdServiceStatus(serviceName string) (models.SystemdServiceDetail, error) {
	return systemd.GetSystemdServiceStatus(serviceName)
}

// GetSystemdServiceLogs retrieves logs for a specific service
func (d *Discover) GetSystemdServiceLogs(serviceName string) string {
	return systemd.GetSystemdServiceLogs(serviceName)
}

// RestartSystemdService attempts to restart a systemd service
func (d *Discover) RestartSystemdService(serviceName string) error {
	return systemd.RestartSystemdService(serviceName)
}

// LoadStateFromFile loads system state from the state file
func (d *Discover) LoadStateFromFile() error {
	loadedState, err := state.LoadState()
	if err != nil {
		return err
	}
	
	d.State = loadedState
	return nil
}

// SaveStateToFile saves the current state to the state file
func (d *Discover) SaveStateToFile() error {
	return state.SaveState(d.State)
}