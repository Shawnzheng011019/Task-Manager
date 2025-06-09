package main

import (
	"log"
	"task-manager/internal/config"
	"task-manager/internal/database"
	"task-manager/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.LoadConfig("configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register routes
	taskHandler := handlers.NewTaskHandler()
	taskHandler.RegisterRoutes(router)
	
	// Register config routes
	configHandler := handlers.NewConfigHandler()
	configHandler.RegisterRoutes(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Task Manager API is running",
		})
	})

	// Start server
	cfg := config.GetConfig()
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
