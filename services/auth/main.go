package main

import (
    "github.com/Decentralized-voting-sytem/backend/services/auth/controllers"
    "github.com/Decentralized-voting-sytem/backend/services/auth/db"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)


func main() {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true, 
    }))

    db.InitDB()
    r.POST("/login", controllers.Login)
    r.POST("/register-vote",controllers.RegisterVote)
    r.POST("/admin-login",controllers.AdminLogin)
    r.Run(":8080")
}
