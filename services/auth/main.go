package main

import (
	"fmt"
	"log"

	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/services/auth/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	database.Init()

	// Check if database connection was established
	db, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Check database connectivity
	if err := db.DB().Ping(); err != nil {
		fmt.Println("Database connection error:", err)
		return
	} else {
		fmt.Println("Database connection successful")
	}

	// Create a new Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Content-Type"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Set up routes
	r.POST("/login", controllers.Login)

	// Start the server
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
