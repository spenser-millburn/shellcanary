package ui

import (
	"fmt"
)

// PauseForUser pauses the program until the user presses Enter
func PauseForUser() {
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
}