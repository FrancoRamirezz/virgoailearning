package main

import (
	"log"
	"net/http"

	"backend/config"
	"backend/database"
	"backend/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	database.Connect(cfg)
	defer database.Close()

	// Run database migrations
	database.Migrate()

	// Setup routes
	router := routes.SetupRoutes(cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Env)
	log.Printf("Database connected: %s", cfg.Database.Name)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
