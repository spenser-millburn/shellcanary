package dockerUI

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"discover/agents/docker"
)

// ShowDockerMenu handles the Docker project menu
func ShowDockerMenu(projectName string) {
	// Get all containers in the selected project
	containers, err := docker.GetDockerContainers(projectName)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Add an option to view logs for all containers and back option
	containerOptions := []string{"🔄 All Containers", "⬅️ Back"}
	containerOptions = append(containerOptions, containers...)
	
	// Create a prompt for selecting a container
	containerPrompt := promptui.Select{
		Label: fmt.Sprintf("🔍 Select a container in project '%s' to view logs", projectName),
		Items: containerOptions,
	}
	
	_, containerSelection, err := containerPrompt.Run()
	if err != nil {
		fmt.Printf("Container selection failed: %v\n", err)
		return
	}
	
	// Handle back option
	if containerSelection == "⬅️ Back" {
		return
	}
	
	var logs string
	if containerSelection == "🔄 All Containers" {
		// Get logs for all containers in the project
		logs = docker.GetAllProjectLogs(projectName)
		fmt.Println(logs)
	} else {
		// Get logs for the selected container
		logs = docker.GetDockerLogs(projectName, containerSelection)
		fmt.Println(logs)
	}
}