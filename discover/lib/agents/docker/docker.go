package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/shellcanary/discover/lib/models"
)

// GetDockerComposeProjects returns a list of running Docker Compose projects
func GetDockerComposeProjects() []models.DockerProject {
	var projects []models.DockerProject

	// Get more container details with a better format
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}|{{.Status}}|{{.Labels}}")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Warning: Docker command failed, might not be installed or running")
		return projects
	}

	containers := strings.Split(strings.TrimSpace(string(output)), "\n")
	projectMap := make(map[string]models.DockerProject)
	
	for _, container := range containers {
		parts := strings.SplitN(container, "|", 3)
		if len(parts) < 3 {
			continue
		}
		containerName, containerStatus, labelsStr := parts[0], parts[1], parts[2]

		// Parse all labels into a map
		labelsMap := make(map[string]string)
		labelPairs := strings.Split(labelsStr, ",")
		for _, pair := range labelPairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				labelsMap[kv[0]] = kv[1]
			}
		}

		projectName := labelsMap["com.docker.compose.project"]
		if projectName == "" {
			continue // Skip if not a docker-compose container
		}

		projectPath := "Unknown"
		if path := labelsMap["com.docker.compose.project.working_dir"]; path != "" {
			projectPath = path
		}

		containerInfo := models.ContainerInfo{
			Name:   containerName,
			Status: containerStatus,
			Labels: labelsMap,
		}

		if proj, exists := projectMap[projectName]; exists {
			proj.Containers++
			proj.ContainerDetails = append(proj.ContainerDetails, containerInfo)
			projectMap[projectName] = proj
		} else {
			projectMap[projectName] = models.DockerProject{
				Name:             projectName, 
				Path:             projectPath, 
				Containers:       1, 
				Status:           "Running",
				ContainerDetails: []models.ContainerInfo{containerInfo},
			}
		}
	}

	for _, project := range projectMap {
		projects = append(projects, project)
	}

	return projects
}

// GetComposeCommand determines which Docker Compose command variant is available
func GetComposeCommand() (string, []string) {
	// Check if 'docker compose' plugin is available
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return "docker", []string{"compose"}
	}
	
	// Fallback to 'docker-compose' if 'docker compose' is not available
	cmd = exec.Command("docker-compose", "version")
	if err := cmd.Run(); err == nil {
		return "docker-compose", []string{}
	}
	
	// If neither is available, return empty strings
	return "", nil
}

// GetDockerContainers retrieves all containers in a Docker Compose project
func GetDockerContainers(projectName string) ([]string, error) {
	baseCmd, args := GetComposeCommand()
	if baseCmd == "" {
		return nil, fmt.Errorf("neither 'docker compose' nor 'docker-compose' is available")
	}
	
	// Construct command to list services in the project
	cmdArgs := append(args, "-p", projectName, "ps", "--services")
	cmd := exec.Command(baseCmd, cmdArgs...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error retrieving containers for Docker project %s: %v", projectName, err)
	}
	
	containers := strings.Fields(string(output))
	if len(containers) == 0 {
		return nil, fmt.Errorf("no containers found in project %s", projectName)
	}
	
	return containers, nil
}

// GetDockerLogs retrieves logs for a specific container in a Docker Compose project
func GetDockerLogs(projectName string, containerName string) string {
	baseCmd, args := GetComposeCommand()
	if baseCmd == "" {
		return "Neither 'docker compose' nor 'docker-compose' is available on this system."
	}
	
	// Construct full command based on which compose variant we're using
	cmdArgs := append(args, "-p", projectName, "logs", containerName)
	cmd := exec.Command(baseCmd, cmdArgs...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error retrieving logs for container %s in project %s: %v", containerName, projectName, err)
	}
	return string(output)
}

// GetAllProjectLogs retrieves logs for all containers in a project
func GetAllProjectLogs(projectName string) string {
	baseCmd, args := GetComposeCommand()
	if baseCmd == "" {
		return "Neither 'docker compose' nor 'docker-compose' is available on this system."
	}
	
	// Construct command for all logs
	cmdArgs := append(args, "-p", projectName, "logs")
	cmd := exec.Command(baseCmd, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error retrieving logs for project %s: %v", projectName, err)
	}
	return string(output)
}