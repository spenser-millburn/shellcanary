package main

import (
	"fmt"
	"os"

	"discover/ui"
	"discover/ui/help"
)

func main() {
	// Check command line args
	if len(os.Args) > 1 {
		// Process command line flags
		switch os.Args[1] {
		case "--help", "-h":
			// Show help info and exit
			help.ShowHelpPage()
			os.Exit(0)
			
		case "--capture-state":
			// Capture system state and exit
			if err := ui.CaptureSystemState(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
			
		default:
			// Unknown flag, show brief usage and exit
			fmt.Printf("Unknown option: %s\n\n", os.Args[1])
			fmt.Println("Usage: discover [OPTION]")
			fmt.Println("  --help, -h          Display help information")
			fmt.Println("  --capture-state     Capture current system state")
			fmt.Println("\nRun without arguments for interactive mode.")
			os.Exit(1)
		}
	}
	
	// Start the interactive menu system
	ui.StartMainMenu()
}