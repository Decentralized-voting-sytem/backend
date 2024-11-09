package main

import (
    "github.com/Decentralized-voting-sytem/backend/services/auth/controllers"
    "github.com/Decentralized-voting-sytem/backend/services/auth/db"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    r := gin.Default()

    // CORS middleware configuration
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // Allow requests only from frontend at 3000
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Allowed methods
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"}, // Allow necessary headers
        AllowCredentials: true, // Allow cookies to be included in requests
        MaxAge:           12 * 3600, 
    }))

    db.InitDB()

    // Set up your routes
    r.POST("/login", controllers.Login)
    r.POST("/register-vote", controllers.RegisterVote)
    r.POST("/admin-login", controllers.AdminLogin)

    // Start the server on port 8080
    err := r.Run(":8080")
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
