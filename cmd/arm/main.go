package main

import (
	"fmt"
	"log"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

func main() {
	// Create the main ARM service - it handles its own dependencies
	armService := arm.NewArmService()

	// Example usage
	version := armService.Version()
	fmt.Printf("ARM Version: %+v\n", version)

	// TODO: Add CLI command handling
	log.Println("ARM service initialized successfully")
}
