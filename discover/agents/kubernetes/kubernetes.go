package kubernetes

import (
	"fmt"
	"os/exec"
	"strings"
	"encoding/json"

	"discover/models"
)

// GetKubernetesConfigs returns a list of Kubernetes configurations with nested namespace and deployment information
func GetKubernetesConfigs() []models.KubernetesConfig {
	var configs []models.KubernetesConfig

	cmd := exec.Command("kubectl", "config", "current-context")
	currentContext, err := cmd.Output()
	if err != nil {
		fmt.Println("Warning: kubectl command failed, might not be installed")
		return configs
	}

	cmd = exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	contextsOutput, err := cmd.Output()
	if err != nil {
		return configs
	}

	contexts := strings.Split(strings.TrimSpace(string(contextsOutput)), "\n")
	current := strings.TrimSpace(string(currentContext))

	for _, context := range contexts {
		status := "Configured"
		if context == current {
			status = "Active"
		}

		// Get namespaces for this context
		namespaces, err := GetNamespacesForContext(context)
		if err != nil {
			fmt.Printf("Warning: Error getting namespaces for context %s: %v\n", context, err)
			namespaces = []models.KubernetesNamespace{} // Use empty array instead of nil
		}

		configs = append(configs, models.KubernetesConfig{
			Name:       context,
			Status:     status,
			Nodes:      "N/A",
			Namespaces: namespaces,
		})
	}

	return configs
}

// GetNamespacesForContext retrieves all namespaces in a Kubernetes context
func GetNamespacesForContext(contextName string) ([]models.KubernetesNamespace, error) {
	cmd := exec.Command("kubectl", "get", "namespaces", "--context", contextName, "-o", "jsonpath={.items[*].metadata.name}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error retrieving namespaces for context %s: %v", contextName, err)
	}

	namespaceNames := strings.Fields(string(output))
	if len(namespaceNames) == 0 {
		return nil, fmt.Errorf("no namespaces found in context %s", contextName)
	}

	var namespaces []models.KubernetesNamespace
	for _, namespaceName := range namespaceNames {
		// Get deployments for this namespace
		deployments, err := GetDeploymentsForNamespace(contextName, namespaceName)
		if err != nil {
			fmt.Printf("Warning: Error getting deployments for namespace %s in context %s: %v\n", 
				namespaceName, contextName, err)
			deployments = []models.KubernetesDeployment{} // Use empty array instead of nil
		}

		namespaces = append(namespaces, models.KubernetesNamespace{
			Name:        namespaceName,
			Deployments: deployments,
		})
	}

	return namespaces, nil
}

// GetDeploymentsForNamespace retrieves all deployments in a specific namespace
func GetDeploymentsForNamespace(contextName, namespaceName string) ([]models.KubernetesDeployment, error) {
	// Get deployments as JSON to get more details
	cmd := exec.Command("kubectl", "get", "deployments", "-n", namespaceName, "--context", contextName, "-o", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error retrieving deployments for namespace %s in context %s: %v", 
			namespaceName, contextName, err)
	}

	// Parse the JSON response
	var result struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Spec struct {
				Replicas int `json:"replicas"`
			} `json:"spec"`
			Status struct {
				AvailableReplicas int `json:"availableReplicas"`
				ReadyReplicas     int `json:"readyReplicas"`
			} `json:"status"`
		} `json:"items"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("error parsing deployments JSON: %v", err)
	}

	var deployments []models.KubernetesDeployment
	for _, item := range result.Items {
		status := "Healthy"
		if item.Status.ReadyReplicas < item.Spec.Replicas {
			status = fmt.Sprintf("Degraded (%d/%d ready)", item.Status.ReadyReplicas, item.Spec.Replicas)
		}

		deployments = append(deployments, models.KubernetesDeployment{
			Name:     item.Metadata.Name,
			Replicas: item.Spec.Replicas,
			Ready:    item.Status.ReadyReplicas,
			Status:   status,
		})
	}

	return deployments, nil
}

// GetKubernetesDeployments retrieves all deployments in a Kubernetes context
func GetKubernetesDeployments(contextName string) ([]string, error) {
	cmd := exec.Command("kubectl", "get", "deployments", "--all-namespaces", "--context", contextName, "-o", "jsonpath={.items[*].metadata.name}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error retrieving deployments for Kubernetes context %s: %v", contextName, err)
	}
	
	deployments := strings.Fields(string(output))
	if len(deployments) == 0 {
		return nil, fmt.Errorf("no deployments found in context %s", contextName)
	}
	
	return deployments, nil
}

// GetKubernetesLogs retrieves logs for a specific deployment in a Kubernetes context
func GetKubernetesLogs(contextName string, deploymentName string) string {
	// First, find the namespace for this deployment
	cmd := exec.Command("kubectl", "get", "deployment", "--all-namespaces", "--context", contextName, 
		"-o", "jsonpath={range .items[?(@.metadata.name==\""+deploymentName+"\")]}{.metadata.namespace}{end}")
	namespaceOutput, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error finding namespace for deployment %s in context %s: %v", deploymentName, contextName, err)
	}
	
	namespace := string(namespaceOutput)
	if namespace == "" {
		return fmt.Sprintf("Could not find deployment %s in context %s", deploymentName, contextName)
	}
	
	// Get logs using the namespace
	cmd = exec.Command("kubectl", "logs", "deployment/"+deploymentName, "-n", namespace, "--context", contextName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error retrieving logs for deployment %s in namespace %s and context %s: %v", 
			deploymentName, namespace, contextName, err)
	}
	return string(output)
}