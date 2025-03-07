package ui

import (
	"fmt"

	"discover/agents/docker"
	"discover/agents/kubernetes"
	"discover/agents/systemd"
	"discover/state"
)

// CaptureSystemState gathers and stores the current state of all resources
func CaptureSystemState() error {
	fmt.Println("Capturing system state...")
	
	// Gather data from all agents
	dockerProjects := docker.GetDockerComposeProjects()
	k8sConfigs := kubernetes.GetKubernetesConfigs()
	systemdServices := systemd.GetSystemdServices()
	
	// Update system state
	err := state.UpdateSystemState(dockerProjects, k8sConfigs, systemdServices)
	if err != nil {
		return fmt.Errorf("error updating system state: %v", err)
	}
	
	fmt.Printf("System state captured and saved to %s\n", state.GetStateFilePath())
	return nil
}