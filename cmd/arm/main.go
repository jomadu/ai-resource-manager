package main

import (
	"fmt"
	"log"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
)

func main() {
	// Initialize components
	configManager := config.NewFileManager()
	manifestManager := manifest.NewFileManager()
	lockFileManager := lockfile.NewFileManager()
	fileInstaller := installer.NewFileInstaller()

	// Initialize cache and registry components
	keyGen := cache.NewGitKeyGen()
	fileCache := cache.NewFileCache()
	gitRepo := registry.NewGitRepo()
	gitRegistry := registry.NewGitRegistry(fileCache, gitRepo, keyGen)

	// Create the main ARM service
	armService := arm.NewArmService(
		configManager,
		manifestManager,
		lockFileManager,
		fileInstaller,
	)

	// Example usage
	version := armService.Version()
	fmt.Printf("ARM Version: %+v\n", version)

	// TODO: Add CLI command handling
	log.Println("ARM service initialized successfully")
	_ = gitRegistry // Will be used when implementing registry operations
}