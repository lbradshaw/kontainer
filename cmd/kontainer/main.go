package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"kontainer/internal/api"
	"kontainer/internal/database"
	"kontainer/internal/service"
)

func main() {
	fmt.Println("Starting Kontainer - Storage Container Inventory Management")

	// Initialize settings service and load settings
	settingsService := service.NewSettingsService()
	settings, err := settingsService.LoadSettings()
	if err != nil {
		log.Fatalf("Failed to load settings: %v", err)
	}

	// Override settings with environment variables if set (for Docker)
	if envPort := os.Getenv("PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil {
			settings.Port = port
			fmt.Printf("Port overridden by environment variable: %d\n", port)
		}
	}

	if envDBPath := os.Getenv("DATABASE_PATH"); envDBPath != "" {
		settings.DatabasePath = envDBPath
		fmt.Printf("Database path overridden by environment variable: %s\n", envDBPath)
	}

	if envTheme := os.Getenv("THEME"); envTheme != "" {
		settings.Theme = envTheme
		fmt.Printf("Theme overridden by environment variable: %s\n", envTheme)
	}

	// Initialize database with configured path
	dbPath := settings.DatabasePath
	if dbPath == "" {
		dbPath = "kontainer.db" // Fallback to default
	}
	fmt.Printf("Using database: %s\n", dbPath)
	
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize service layer
	toteService := service.NewToteService(db)

	// Initialize API routes
	router := api.NewRouter(toteService, settingsService)

	// Start server with configured port
	port := fmt.Sprintf(":%d", settings.Port)
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Printf("Access from other devices on your network at: http://<your-ip>%s\n", port)
	fmt.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
