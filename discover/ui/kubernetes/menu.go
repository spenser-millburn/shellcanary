package kubernetesUI

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"discover/agents/kubernetes"
)

// ShowKubernetesMenu handles the Kubernetes context menu
func ShowKubernetesMenu(contextName string) {
	// Get all deployments in the selected context
	deployments, err := kubernetes.GetKubernetesDeployments(contextName)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Add back option
	deploymentOptions := append([]string{"⬅️ Back"}, deployments...)
	
	// Create a prompt for selecting a deployment
	deploymentPrompt := promptui.Select{
		Label: fmt.Sprintf("🔍 Select a deployment in context '%s' to view logs", contextName),
		Items: deploymentOptions,
	}
	
	_, deploymentName, err := deploymentPrompt.Run()
	if err != nil {
		fmt.Printf("Deployment selection failed: %v\n", err)
		return
	}
	
	// Handle back option
	if deploymentName == "⬅️ Back" {
		return
	}
	
	// Get logs for the selected deployment
	logs := kubernetes.GetKubernetesLogs(contextName, deploymentName)
	fmt.Println(logs)
}