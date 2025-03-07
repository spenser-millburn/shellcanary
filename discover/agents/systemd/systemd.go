package systemd

import (
	"fmt"
	"os/exec"
	"strings"
	"regexp"

	"discover/models"
)

// GetSystemdServices returns a list of systemd services
func GetSystemdServices() []models.SystemdService {
	var services []models.SystemdService

	// Check if systemctl is available
	cmd := exec.Command("systemctl", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("Warning: systemctl command failed, systemd might not be available")
		return services
	}

	// Get list of all services
	cmd = exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Warning: Failed to list systemd services:", err)
		return services
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// Skip the header line and process each service
	for i, line := range lines {
		if i == 0 { // Skip header
			continue
		}
		
		// Remove extra spaces and split by whitespace
		line = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(line), " ")
		parts := strings.SplitN(line, " ", 5)
		
		// Check if we have enough parts
		if len(parts) < 4 {
			continue
		}
		
		// Extract service name without .service suffix
		serviceName := parts[0]
		if strings.HasSuffix(serviceName, ".service") {
			serviceName = serviceName[:len(serviceName)-8] // Remove ".service"
		}
		
		// Extract status and description
		loadStatus := parts[1]
		activeStatus := parts[2]
		subStatus := parts[3]
		description := ""
		if len(parts) >= 5 {
			description = parts[4]
		}
		
		// Combine load and active status
		status := activeStatus
		if loadStatus != "loaded" {
			status = fmt.Sprintf("%s (%s)", activeStatus, loadStatus)
		}
		
		services = append(services, models.SystemdService{
			Name:        serviceName,
			Status:      status,
			SubStatus:   subStatus,
			Description: description,
		})
	}

	return services
}

// GetSystemdServiceStatus retrieves the detailed status of a specific systemd service
func GetSystemdServiceStatus(serviceName string) (models.SystemdServiceDetail, error) {
	var detail models.SystemdServiceDetail
	
	// Construct the service name with .service suffix if not present
	if !strings.HasSuffix(serviceName, ".service") {
		serviceName = serviceName + ".service"
	}
	
	// Get service properties
	cmd := exec.Command("systemctl", "show", 
		"--property=Id,Description,LoadState,ActiveState,SubState,UnitFileState,ExecMainPID,ExecMainStatus,Type,Restart",
		serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return detail, fmt.Errorf("error retrieving details for service %s: %v", serviceName, err)
	}
	
	// Parse the output
	properties := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			properties[parts[0]] = parts[1]
		}
	}
	
	// Create the service detail
	detail = models.SystemdServiceDetail{
		Id:             properties["Id"],
		Description:    properties["Description"],
		LoadState:      properties["LoadState"],
		ActiveState:    properties["ActiveState"],
		SubState:       properties["SubState"],
		UnitFileState:  properties["UnitFileState"],
		ExecMainPID:    properties["ExecMainPID"],
		ExecMainStatus: properties["ExecMainStatus"],
		Type:           properties["Type"],
		Restart:        properties["Restart"],
	}
	
	return detail, nil
}

// GetSystemdServiceLogs retrieves logs for a specific systemd service
func GetSystemdServiceLogs(serviceName string) string {
	// Ensure service name has .service suffix
	if !strings.HasSuffix(serviceName, ".service") {
		serviceName = serviceName + ".service"
	}
	
	// Use journalctl to get logs for the service
	cmd := exec.Command("journalctl", "-u", serviceName, "--no-pager", "-n", "100")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error retrieving logs for service %s: %v", serviceName, err)
	}
	return string(output)
}

// RestartSystemdService attempts to restart a systemd service
func RestartSystemdService(serviceName string) error {
	// Ensure service name has .service suffix
	if !strings.HasSuffix(serviceName, ".service") {
		serviceName = serviceName + ".service"
	}
	
	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart service %s: %v\nOutput: %s", 
			serviceName, err, string(output))
	}
	
	return nil
}