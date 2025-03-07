package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"discover/models"
)

const stateFileName = "discover_state.json"

// GetStateFilePath returns the path to the state file
func GetStateFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if we can't get home dir
		return stateFileName
	}
	
	// Create a .discover directory in the user's home directory
	discoverDir := filepath.Join(homeDir, ".discover")
	if _, err := os.Stat(discoverDir); os.IsNotExist(err) {
		if err := os.Mkdir(discoverDir, 0755); err != nil {
			fmt.Printf("Warning: Could not create directory %s: %v\n", discoverDir, err)
			return stateFileName
		}
	}
	
	return filepath.Join(discoverDir, stateFileName)
}

// LoadState loads the system state from the state file
func LoadState() (models.SystemState, error) {
	var state models.SystemState
	
	stateFile := GetStateFilePath()
	data, err := ioutil.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty state if the file doesn't exist
			return models.SystemState{
				LastUpdated: time.Now(),
			}, nil
		}
		return state, fmt.Errorf("error reading state file: %v", err)
	}
	
	if err := json.Unmarshal(data, &state); err != nil {
		return state, fmt.Errorf("error parsing state file: %v", err)
	}
	
	return state, nil
}

// SaveState saves the system state to the state file
func SaveState(state models.SystemState) error {
	// Ensure the LastUpdated field is set to current time
	state.LastUpdated = time.Now()
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding state: %v", err)
	}
	
	stateFile := GetStateFilePath()
	if err := ioutil.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("error writing state file: %v", err)
	}
	
	return nil
}


// UpdateSystemState updates the full system state with current data
func UpdateSystemState(
	dockerProjects []models.DockerProject,
	k8sConfigs []models.KubernetesConfig,
	systemdServices []models.SystemdService,
) error {
	// First load existing state to preserve logs
	state, err := LoadState()
	if err != nil {
		return err
	}
	
	// Update the resource data
	state.DockerProjects = dockerProjects
	state.KubernetesConfigs = k8sConfigs
	state.SystemdServices = systemdServices
	
	return SaveState(state)
}
