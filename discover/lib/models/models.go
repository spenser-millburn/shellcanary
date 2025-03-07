package models

import "time"

// SystemdServiceDetail represents detailed information about a systemd service
type SystemdServiceDetail struct {
	Id             string
	Description    string
	LoadState      string
	ActiveState    string
	SubState       string
	UnitFileState  string
	ExecMainPID    string
	ExecMainStatus string
	Type           string
	Restart        string
}

// DockerProject represents a Docker Compose project
type DockerProject struct {
	Name             string
	Path             string
	Containers       int
	Status           string
	ContainerDetails []ContainerInfo
}

// ContainerInfo represents details about a container in a Docker project
type ContainerInfo struct {
	Name   string
	Status string
	Labels map[string]string
}

// KubernetesDeployment represents a deployment in Kubernetes
type KubernetesDeployment struct {
	Name     string
	Replicas int
	Ready    int
	Status   string
}

// KubernetesNamespace represents a namespace in Kubernetes
type KubernetesNamespace struct {
	Name        string
	Deployments []KubernetesDeployment
}

// KubernetesConfig represents a Kubernetes configuration
type KubernetesConfig struct {
	Name       string
	Status     string
	Nodes      string
	Namespaces []KubernetesNamespace
}

// SystemdService represents a systemd service
type SystemdService struct {
	Name        string
	Status      string
	SubStatus   string
	Description string
}

// LogEntry represents a log entry in the state file
type LogEntry struct {
	DataType    string    `json:"data_type"`
	Project     string    `json:"project"`
	Container   string    `json:"container,omitempty"`
	Deployment  string    `json:"deployment,omitempty"`
	LogContent  string    `json:"log_content"`
	Timestamp   time.Time `json:"timestamp"`
}

// SystemState represents the entire system state
type SystemState struct {
	DockerProjects    []DockerProject    `json:"docker_compose_projects"`
	KubernetesConfigs []KubernetesConfig `json:"kubernetes_projects"`
	SystemdServices   []SystemdService   `json:"systemd_services,omitempty"`
	LastUpdated       time.Time          `json:"last_updated"`
}