package systemdUI

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"discover/agents/systemd"
)

// ShowSystemdMenu handles the systemd service menu
func ShowSystemdMenu(serviceName string) {
	// Create a prompt for service actions
	actionPrompt := promptui.Select{
		Label: fmt.Sprintf("🔍 Select an action for service '%s'", serviceName),
		Items: []string{"📜 View Logs", "📊 View Details", "🔄 Restart Service", "⬅️ Back"},
	}
	
	_, actionSelection, err := actionPrompt.Run()
	if err != nil {
		fmt.Printf("Action selection failed: %v\n", err)
		return
	}
	
	// Handle back option
	if actionSelection == "⬅️ Back" {
		return
	}
	
	switch actionSelection {
	case "📜 View Logs":
		logs := systemd.GetSystemdServiceLogs(serviceName)
		fmt.Println(logs)
		
	case "📊 View Details":
		details, err := systemd.GetSystemdServiceStatus(serviceName)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Printf("Service: %s\n", details.Id)
		fmt.Printf("Description: %s\n", details.Description)
		fmt.Printf("Load State: %s\n", details.LoadState)
		fmt.Printf("Active State: %s\n", details.ActiveState)
		fmt.Printf("Sub State: %s\n", details.SubState)
		fmt.Printf("Unit File State: %s\n", details.UnitFileState)
		fmt.Printf("Main PID: %s\n", details.ExecMainPID)
		fmt.Printf("Main Status: %s\n", details.ExecMainStatus)
		fmt.Printf("Type: %s\n", details.Type)
		fmt.Printf("Restart: %s\n", details.Restart)
		
	case "🔄 Restart Service":
		fmt.Printf("Restarting service %s...\n", serviceName)
		if err := systemd.RestartSystemdService(serviceName); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Service %s restarted successfully\n", serviceName)
		}
	}
}