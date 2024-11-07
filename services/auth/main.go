package main

import (
    "github.com/Decentralized-voting-sytem/backend/services/auth/controllers"
    "github.com/Decentralized-voting-sytem/backend/services/auth/db"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)


func main() {
    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"}
    config.AllowHeaders = []string{"Content-Type"}
    config.AllowCredentials = true

    r.Use(cors.New(config))
    db.InitDB()
    r.POST("/login", controllers.Login)
    r.POST("/register-vote",controllers.RegisterVote)
    r.POST("/admin-login",controllers.AdminLogin)
    r.Run()
}
