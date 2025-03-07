package help

import (
	"fmt"
)

// ShowHelpPage displays the application help information
func ShowHelpPage() {
	fmt.Println(`
================================================
DISCOVER - System Resource Management Tool
================================================

This application helps you discover and manage Docker, Kubernetes, and Systemd 
resources on your system.

MAIN FUNCTIONALITY:
------------------
• Resource Discovery: Find and list various system resources
• Log Viewing: Access logs for containers, deployments, and services
• Resource Details: View detailed information about system components
• State Capture: Save the current system state for future reference

RESOURCE TYPES:
--------------
🐳 Docker:
   - View Docker Compose projects and their containers
   - Access logs for specific containers or entire projects

☸️ Kubernetes:
   - Browse Kubernetes contexts, namespaces, and deployments
   - View deployment logs and status information

⚙️ Systemd:
   - List active systemd services
   - View service logs, status details, and perform restarts

NAVIGATION TIPS:
--------------
• Use arrow keys to navigate menus
• Press Enter to select an option
• Select "Back" options to return to previous menus
• Select "Exit Application" from the main menu to quit

COMMAND LINE USAGE:
-----------------
$ discover [OPTION]

Options:
  --help, -h          Display this help information
  --capture-state     Capture the current system state and exit

Running without arguments launches the interactive interface.
The system state is saved to ~/.discover/discover_state.json

================================================
`)
}